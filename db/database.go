package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

type Deck struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type Card struct {
	ID            int     `json:"id"`
	DeckID        int     `json:"deck_id"`
	KanjiText     string  `json:"kanji_text"`
	FuriganaText  *string `json:"furigana_text"`
	Translation   string  `json:"translation"`
	EaseFactor    float64 `json:"ease_factor"`
	Interval      int     `json:"interval"`
	Repetitions   int     `json:"repetitions"`
	NextReview    string  `json:"next_review"`
	LastReview    *string `json:"last_review"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type CardInput struct {
	DeckID       int     `json:"deck_id"`
	KanjiText    string  `json:"kanji_text"`
	FuriganaText *string `json:"furigana_text"`
	Translation  string  `json:"translation"`
}

type TrainingCard struct {
	ID           int     `json:"id"`
	KanjiText    string  `json:"kanji_text"`
	FuriganaText *string `json:"furigana_text"`
	Translation  string  `json:"translation"`
}

type VoiceStatus struct {
	Ja bool `json:"ja"`
	Ru bool `json:"ru"`
}

func InitDB() error {
	appData, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config dir: %w", err)
	}

	appDir := filepath.Join(appData, "Yappari")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("failed to create app dir: %w", err)
	}

	dbPath := filepath.Join(appDir, "database.db")
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func runMigrations() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS decks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT (datetime('now'))
		);`,
		`CREATE TABLE IF NOT EXISTS cards (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			deck_id INTEGER NOT NULL,
			kanji_text TEXT NOT NULL,
			furigana_text TEXT,
			translation TEXT NOT NULL,
			ease_factor REAL NOT NULL DEFAULT 2.5,
			interval INTEGER NOT NULL DEFAULT 0,
			repetitions INTEGER NOT NULL DEFAULT 0,
			next_review TEXT NOT NULL,
			last_review TEXT,
			created_at TEXT NOT NULL DEFAULT (datetime('now')),
			updated_at TEXT NOT NULL DEFAULT (datetime('now')),
			FOREIGN KEY (deck_id) REFERENCES decks(id) ON DELETE CASCADE
		);`,
	}

	for _, m := range migrations {
		if _, err := DB.Exec(m); err != nil {
			return err
		}
	}

	return nil
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
