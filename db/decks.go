package db

import (
	"time"
)

func GetDecks() ([]Deck, error) {
	rows, err := DB.Query("SELECT id, name, created_at FROM decks ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var decks []Deck
	for rows.Next() {
		var d Deck
		if err := rows.Scan(&d.ID, &d.Name, &d.CreatedAt); err != nil {
			return nil, err
		}
		decks = append(decks, d)
	}
	return decks, nil
}

func CreateDeck(name string) (*Deck, error) {
	result, err := DB.Exec("INSERT INTO decks (name) VALUES (?)", name)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	d := &Deck{
		ID:        int(id),
		Name:      name,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	return d, nil
}

func UpdateDeck(id int, name string) error {
	_, err := DB.Exec("UPDATE decks SET name = ? WHERE id = ?", name, id)
	return err
}

func DeleteDeck(id int) error {
	_, err := DB.Exec("DELETE FROM decks WHERE id = ?", id)
	return err
}
