package tts

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// voiceFor returns the edge-tts voice name for the given language code.
func voiceFor(lang string) string {
	switch lang {
	case "ja", "ja-JP":
		return "ja-JP-NanamiNeural"
	case "ru", "ru-RU":
		return "ru-RU-SvetlanaNeural"
	default:
		return "ja-JP-NanamiNeural"
	}
}

// runEdgeTTS runs edge-tts with the given args, trying multiple invocation methods.
// It tries: edge-tts → python -m edge_tts → python3 -m edge_tts.
// The audio is written to --write-media path; this function only returns errors.
func runEdgeTTS(args []string) error {
	// Try direct edge-tts command first
	cmd := exec.Command("edge-tts", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err == nil {
		return nil
	}
	edgeTTSErr := stderr.String()

	// Fallback: python -m edge_tts
	pyArgs := append([]string{"-m", "edge_tts"}, args...)
	cmd = exec.Command("python", pyArgs...)
	stderr.Reset()
	cmd.Stderr = &stderr
	if err := cmd.Run(); err == nil {
		return nil
	}
	pyErr := stderr.String()

	// Fallback: python3 -m edge_tts
	cmd = exec.Command("python3", pyArgs...)
	stderr.Reset()
	cmd.Stderr = &stderr
	if err := cmd.Run(); err == nil {
		return nil
	}
	py3Err := stderr.String()

	return fmt.Errorf("edge-tts не найден. Установите: pip install edge-tts\n"+
		"  edge-tts: %s\n  python -m edge_tts: %s\n  python3 -m edge_tts: %s",
		edgeTTSErr, pyErr, py3Err)
}

// Speak calls edge-tts to synthesize the given text in the given language.
// Returns the MP3 audio bytes.
func Speak(text, lang string) ([]byte, error) {
	voice := voiceFor(lang)

	tmpDir, err := os.MkdirTemp("", "yappari-tts")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	outPath := filepath.Join(tmpDir, "speak.mp3")

	args := []string{
		"--voice", voice,
		"--text", text,
		"--write-media", outPath,
	}

	if err := runEdgeTTS(args); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio file: %w", err)
	}

	return data, nil
}

// CheckAvailability checks if edge-tts is available.
// Returns true if available, along with a human-readable message.
func CheckAvailability() (bool, string) {
	// Try edge-tts --list-voices
	cmd := exec.Command("edge-tts", "--list-voices")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err == nil {
		return true, "edge-tts доступен"
	}

	// Try python -m edge_tts --list-voices
	cmd = exec.Command("python", "-m", "edge_tts", "--list-voices")
	var pyOut bytes.Buffer
	cmd.Stdout = &pyOut

	if err := cmd.Run(); err == nil {
		return true, "edge-tts доступен (python -m)"
	}

	// Try python3 -m edge_tts --list-voices
	cmd = exec.Command("python3", "-m", "edge_tts", "--list-voices")
	var py3Out bytes.Buffer
	cmd.Stdout = &py3Out

	if err := cmd.Run(); err == nil {
		return true, "edge-tts доступен (python3 -m)"
	}

	return false, "edge-tts не найден. Установите: pip install edge-tts"
}
