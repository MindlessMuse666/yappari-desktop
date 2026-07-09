package main

import (
	"context"
	"encoding/base64"

	"github.com/MindlessMuse666/yappari/db"
	"github.com/MindlessMuse666/yappari/tts"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if err := db.InitDB(); err != nil {
		panic(err)
	}
}

func (a *App) shutdown(ctx context.Context) {
	if err := db.CloseDB(); err != nil {
		panic(err)
	}
}

func (a *App) GetDecks() ([]db.Deck, error) {
	return db.GetDecks()
}

func (a *App) CreateDeck(name string) (*db.Deck, error) {
	return db.CreateDeck(name)
}

func (a *App) UpdateDeck(id int, name string) error {
	return db.UpdateDeck(id, name)
}

func (a *App) DeleteDeck(id int) error {
	return db.DeleteDeck(id)
}

func (a *App) GetCardsByDeck(deckID int) ([]db.Card, error) {
	return db.GetCardsByDeck(deckID)
}

func (a *App) CreateCard(input db.CardInput) (*db.Card, error) {
	return db.CreateCard(input)
}

func (a *App) UpdateCard(id int, input db.CardInput) error {
	return db.UpdateCard(id, input)
}

func (a *App) DeleteCard(id int) error {
	return db.DeleteCard(id)
}

func (a *App) GetTrainingCards(mode string, deckIDs []int) ([]db.TrainingCard, error) {
	return db.GetTrainingCards(mode, deckIDs)
}

func (a *App) SubmitReview(cardID int, grade int) (*db.Card, error) {
	return db.SubmitReview(cardID, grade)
}

func (a *App) ResetCardProgress(cardID int) error {
	return db.ResetCardProgress(cardID)
}

func (a *App) ResetDeckProgress(deckID int) error {
	return db.ResetDeckProgress(deckID)
}

func (a *App) CheckVoicesAvailability() db.VoiceStatus {
	return db.VoiceStatus{
		Ja: true,
		Ru: true,
	}
}

func (a *App) SpeakText(text string, lang string) (string, error) {
	data, err := tts.Speak(text, lang)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (a *App) CheckEdgeTTSAvailability() map[string]any {
	ok, msg := tts.CheckAvailability()
	return map[string]any{
		"available": ok,
		"message":   msg,
	}
}
