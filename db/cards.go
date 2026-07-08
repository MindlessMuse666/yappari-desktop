package db

import (
	"math/rand"
	"strings"
	"time"

	"github.com/MindlessMuse666/yappari/sm2"
)

func GetCardsByDeck(deckID int) ([]Card, error) {
	rows, err := DB.Query("SELECT id, deck_id, kanji_text, furigana_text, translation, ease_factor, interval, repetitions, next_review, last_review, created_at, updated_at FROM cards WHERE deck_id = ? ORDER BY created_at DESC", deckID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []Card
	for rows.Next() {
		var c Card
		err := rows.Scan(
			&c.ID,
			&c.DeckID,
			&c.KanjiText,
			&c.FuriganaText,
			&c.Translation,
			&c.EaseFactor,
			&c.Interval,
			&c.Repetitions,
			&c.NextReview,
			&c.LastReview,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, nil
}

func CreateCard(input CardInput) (*Card, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	result, err := DB.Exec(`
		INSERT INTO cards (deck_id, kanji_text, furigana_text, translation, next_review)
		VALUES (?, ?, ?, ?, ?)
	`, input.DeckID, input.KanjiText, input.FuriganaText, input.Translation, now)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	c := &Card{
		ID:           int(id),
		DeckID:       input.DeckID,
		KanjiText:    input.KanjiText,
		FuriganaText: input.FuriganaText,
		Translation:  input.Translation,
		EaseFactor:   2.5,
		Interval:     0,
		Repetitions:  0,
		NextReview:   now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return c, nil
}

func UpdateCard(id int, input CardInput) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := DB.Exec(`
		UPDATE cards
		SET deck_id = ?, kanji_text = ?, furigana_text = ?, translation = ?, updated_at = ?
		WHERE id = ?
	`, input.DeckID, input.KanjiText, input.FuriganaText, input.Translation, now, id)
	return err
}

func DeleteCard(id int) error {
	_, err := DB.Exec("DELETE FROM cards WHERE id = ?", id)
	return err
}

func GetTrainingCards(mode string, deckIDs []int) ([]TrainingCard, error) {
	if len(deckIDs) == 0 {
		return []TrainingCard{}, nil
	}

	placeholders := make([]any, len(deckIDs))
	for i, id := range deckIDs {
		placeholders[i] = id
	}

	query := `
		SELECT id, kanji_text, furigana_text, translation
		FROM cards
		WHERE deck_id IN (?` + strings.Repeat(",?", len(deckIDs)-1) + `)
	`

	if mode == "interval" {
		now := time.Now().UTC().Format(time.RFC3339)
		query += ` AND next_review <= ?`
		placeholders = append(placeholders, now)
	}

	rows, err := DB.Query(query, placeholders...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []TrainingCard
	for rows.Next() {
		var c TrainingCard
		if err := rows.Scan(&c.ID, &c.KanjiText, &c.FuriganaText, &c.Translation); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}

	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	return cards, nil
}

func GetCardByID(id int) (*Card, error) {
	var c Card
	err := DB.QueryRow(`
		SELECT id, deck_id, kanji_text, furigana_text, translation, ease_factor, interval, repetitions, next_review, last_review, created_at, updated_at
		FROM cards WHERE id = ?
	`, id).Scan(
		&c.ID,
		&c.DeckID,
		&c.KanjiText,
		&c.FuriganaText,
		&c.Translation,
		&c.EaseFactor,
		&c.Interval,
		&c.Repetitions,
		&c.NextReview,
		&c.LastReview,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func SubmitReview(cardID int, grade int) (*Card, error) {
	if err := sm2.ValidateGrade(grade); err != nil {
		return nil, err
	}

	card, err := GetCardByID(cardID)
	if err != nil {
		return nil, err
	}

	result := sm2.Calculate(card.EaseFactor, card.Interval, card.Repetitions, grade)
	now := time.Now().UTC().Format(time.RFC3339)

	_, err = DB.Exec(`
		UPDATE cards
		SET ease_factor = ?, interval = ?, repetitions = ?, next_review = ?, last_review = ?, updated_at = ?
		WHERE id = ?
	`, result.EaseFactor, result.Interval, result.Repetitions, result.NextReview, now, now, cardID)
	if err != nil {
		return nil, err
	}

	updatedCard, err := GetCardByID(cardID)
	return updatedCard, err
}

func ResetCardProgress(cardID int) error {
	result := sm2.Reset()
	now := time.Now().UTC().Format(time.RFC3339)

	_, err := DB.Exec(`
		UPDATE cards
		SET ease_factor = ?, interval = ?, repetitions = ?, next_review = ?, last_review = NULL, updated_at = ?
		WHERE id = ?
	`, result.EaseFactor, result.Interval, result.Repetitions, result.NextReview, now, cardID)
	return err
}

func ResetDeckProgress(deckID int) error {
	result := sm2.Reset()
	now := time.Now().UTC().Format(time.RFC3339)

	_, err := DB.Exec(`
		UPDATE cards
		SET ease_factor = ?, interval = ?, repetitions = ?, next_review = ?, last_review = NULL, updated_at = ?
		WHERE deck_id = ?
	`, result.EaseFactor, result.Interval, result.Repetitions, result.NextReview, now, deckID)
	return err
}
