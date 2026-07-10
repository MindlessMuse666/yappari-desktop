// Пакет tts реализует синтез речи через Microsoft Edge TTS (edge-tts)
// и встроенный Windows TTS (SAPI через PowerShell).
//
// Порядок выбора движка:
//  1. edge-tts (Python) — наилучшее качество, поддержка японского и русского
//  2. Windows TTS (System.Speech) — fallback, если Python не установлен
//
// Для edge-tts требуется установка:
//	pip install edge-tts
//
// Windows TTS доступен на Windows Vista+ без дополнительных зависимостей,
// но требует установленного языкового пакета для нужного языка.
package tts

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// voiceFor возвращает имя голоса edge-tts для указанного кода языка.
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

// tryInvoke пытается выполнить edge-tts с указанными аргументами через три
// возможных способа: прямой вызов edge-tts, python -m edge_tts, python3 -m edge_tts.
// Возвращает первый успешный результат.
func tryInvoke(args []string) error {
	attempts := []struct {
		name string
		args []string
	}{
		{"edge-tts", args},
		{"python", append([]string{"-m", "edge_tts"}, args...)},
		{"python3", append([]string{"-m", "edge_tts"}, args...)},
	}

	var lastErr error
	for _, a := range attempts {
		cmd := exec.Command(a.name, a.args...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err == nil {
			return nil
		}
		lastErr = fmt.Errorf("%s: %s", a.name, stderr.String())
	}

	return fmt.Errorf("edge-tts не найден. Установите: pip install edge-tts\n  %s", lastErr)
}

// tryInvokeWithOutput пытается выполнить edge-tts и вернуть stdout.
// Возвращает первый успешный результат с выводом.
func tryInvokeWithOutput(args []string) (bool, string, error) {
	attempts := []struct {
		name string
		args []string
	}{
		{"edge-tts", args},
		{"python", append([]string{"-m", "edge_tts"}, args...)},
		{"python3", append([]string{"-m", "edge_tts"}, args...)},
	}

	var lastErr error
	for _, a := range attempts {
		cmd := exec.Command(a.name, a.args...)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err == nil {
			return true, stdout.String(), nil
		}
		lastErr = fmt.Errorf("%s: %s", a.name, stderr.String())
	}

	return false, lastErr.Error(), lastErr
}

// speakEdgeTTS вызывает edge-tts для синтеза указанного текста.
// Возвращает MP3-аудио в виде среза байт.
func speakEdgeTTS(text, lang string) ([]byte, error) {
	voice := voiceFor(lang)

	tmpDir, err := os.MkdirTemp("", "yappari-tts")
	if err != nil {
		return nil, fmt.Errorf("не удалось создать временную папку: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	outPath := filepath.Join(tmpDir, "speak.mp3")

	args := []string{
		"--voice", voice,
		"--text", text,
		"--write-media", outPath,
	}

	if err := tryInvoke(args); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать аудиофайл: %w", err)
	}

	return data, nil
}

// speakWindowsTTS вызывает встроенный Windows TTS (System.Speech) через PowerShell.
// Возвращает WAV-аудио в виде среза байт.
//
// Использует PowerShell-скрипт, который:
//  1. Загружает сборку System.Speech
//  2. Ищет установленный голос по коду языка (ja-JP / ru-RU)
//  3. Синтезирует речь в WAV-файл
//
// Текст передаётся через временный файл, чтобы избежать проблем с экранированием.
func speakWindowsTTS(text, lang string) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "yappari-tts")
	if err != nil {
		return nil, fmt.Errorf("Windows TTS: не удалось создать временную папку: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	outPath := filepath.Join(tmpDir, "speak.wav")
	textPath := filepath.Join(tmpDir, "text.txt")
	scriptPath := filepath.Join(tmpDir, "speak.ps1")

	// Пишем текст во временный файл (UTF-8 без BOM)
	if err := os.WriteFile(textPath, []byte(text), 0644); err != nil {
		return nil, fmt.Errorf("Windows TTS: не удалось записать текст: %w", err)
	}

	// Приводим язык к формату Windows: ja-JP, ru-RU
	culture := strings.ReplaceAll(lang, "_", "-")

	// PowerShell-скрипт: ищет голос по культуре, синтезирует речь в WAV
	psScript := fmt.Sprintf(`
$text = Get-Content '%s' -Raw -Encoding UTF8
Add-Type -AssemblyName System.Speech
$s = New-Object System.Speech.Synthesis.SpeechSynthesizer
$voice = $s.GetInstalledVoices() | Where-Object { $_.VoiceInfo.Culture.Name -eq '%s' } | Select-Object -First 1
if ($voice -ne $null) { $s.SelectVoice($voice.VoiceInfo.Name) }
$s.SetOutputToWaveFile('%s')
$s.Speak([string]$text)
$s.Dispose()
`, textPath, culture, outPath)

	if err := os.WriteFile(scriptPath, []byte(psScript), 0644); err != nil {
		return nil, fmt.Errorf("Windows TTS: не удалось записать скрипт: %w", err)
	}

	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("Windows TTS: %s", strings.TrimSpace(stderr.String()))
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		return nil, fmt.Errorf("Windows TTS: не удалось прочитать аудиофайл: %w", err)
	}

	return data, nil
}

// Speak вызывает синтез речи, последовательно перебирая доступные движки:
//
//  1. edge-tts (Python) — MP3
//  2. Windows TTS (System.Speech) — WAV
//
// Возвращает аудиоданные, MIME-тип и ошибку.
//
// Параметры:
//   - text — текст для озвучивания
//   - lang — код языка ("ja" для японского, "ru" для русского)
func Speak(text, lang string) ([]byte, string, error) {
	// Попытка 1: edge-tts (наилучшее качество)
	data, err := speakEdgeTTS(text, lang)
	if err == nil {
		return data, "audio/mpeg", nil
	}

	// Попытка 2: Windows TTS (fallback)
	data, err = speakWindowsTTS(text, lang)
	if err == nil {
		return data, "audio/wav", nil
	}

	return nil, "",
		fmt.Errorf("TTS недоступен. Установите edge-tts (pip install edge-tts)")
}

// checkWindowsTTSAvailability проверяет, доступен ли Windows TTS для указанного языка.
func checkWindowsTTSAvailability(lang string) (bool, string) {
	culture := strings.ReplaceAll(lang, "_", "-")

	psScript := fmt.Sprintf(`
Add-Type -AssemblyName System.Speech
$s = New-Object System.Speech.Synthesis.SpeechSynthesizer
$voice = $s.GetInstalledVoices() | Where-Object { $_.VoiceInfo.Culture.Name -eq '%s' } | Select-Object -First 1
if ($voice -ne $null) { $true } else { $false }
`, culture)

	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psScript)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return false, strings.TrimSpace(stderr.String())
	}

	result := strings.TrimSpace(stdout.String())
	if result == "True" {
		return true, fmt.Sprintf("Windows TTS доступен для %s", culture)
	}
	return false, fmt.Sprintf("Windows TTS: голос для %s не найден", culture)
}

// CheckAvailability проверяет, доступен ли хотя бы один движок синтеза речи.
// Возвращает признак доступности и информационное сообщение.
func CheckAvailability() (bool, string) {
	// Проверка 1: edge-tts
	if ok, _, _ := tryInvokeWithOutput([]string{"--list-voices"}); ok {
		return true, "edge-tts доступен"
	}

	// Проверка 2: Windows TTS (хотя бы для одного из целевых языков)
	for _, lang := range []string{"ja-JP", "ru-RU"} {
		if ok, msg := checkWindowsTTSAvailability(lang); ok {
			return true, msg
		}
	}

	return false,
		"TTS недоступен. Установите edge-tts (pip install edge-tts) " +
			"или языковой пакет для японского/русского в Windows"
}
