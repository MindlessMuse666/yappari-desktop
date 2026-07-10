"""
Скрипт первоначальной настройки TTS-окружения.

Создаёт виртуальное окружение Python, устанавливает зависимости
и скачивает модели Silero и Kokoro для офлайн-работы.

Вызывается автоматически при первом запуске приложения (из tts.go/ensureVenv)
или вручную через `make install-tts`.

В интерактивном терминале — человекочитаемый вывод.
При захваченном stdout (Go подпроцесс) — JSON-строки для парсинга.
"""

import argparse
import json
import os
import subprocess
import sys


def log(msg: str, phase: str = "progress"):
    """Выводит сообщение: в терминале — человеческим языком, при захвате — JSON."""
    if sys.stdout.isatty():
        icons = {"venv": "📁", "deps": "📦", "models": "🤖", "error": "❌", "done": "✅"}
        print(f"{icons.get(phase, '•')} {msg}", flush=True)
    else:
        data = {"type": phase, "phase": phase, "message": msg, "percent": 0}
        sys.stdout.write(json.dumps(data, ensure_ascii=False) + "\n")
        sys.stdout.flush()


def log_done(success: bool, message: str = ""):
    """Выводит сообщение о завершении."""
    if sys.stdout.isatty():
        status = "✅" if success else "❌"
        print(f"{status} {message}", flush=True)
    else:
        data = {"type": "complete", "success": success, "message": message}
        sys.stdout.write(json.dumps(data, ensure_ascii=False) + "\n")
        sys.stdout.flush()


def run_pip(pip_exe: str, args: list, desc: str) -> bool:
    """Запускает pip и логирует результат."""
    # --no-cache-dir чтобы избежать битого кэша на Windows
    cmd = [pip_exe] + args + ["--no-cache-dir"]
    try:
        log(f"{desc}...", "deps")
        result = subprocess.run(cmd, capture_output=True, text=True)
        if result.returncode != 0:
            stderr = result.stderr.strip()
            stdout = result.stdout.strip()
            log(f"STDERR: {stderr[-500:]}", "error")
            if stdout:
                log(f"STDOUT: {stdout[-500:]}", "error")
            return False
        # Показываем последние строки pip (прогресс установки)
        out = result.stdout.strip()
        if out:
            for line in out.split("\n")[-3:]:
                line = line.strip()
                if line:
                    log(f"  {line}", "deps")
        return True
    except FileNotFoundError:
        log(f"{desc}: команда не найдена", "error")
        return False


def setup_venv(venv_dir: str) -> bool:
    """Создаёт виртуальное окружение, если его нет."""
    python_exe = os.path.join(venv_dir, "Scripts", "python.exe")
    if os.path.exists(python_exe):
        log("Виртуальное окружение уже существует", "venv")
        return True

    log("Создание виртуального окружения...", "venv")
    result = subprocess.run(
        [sys.executable, "-m", "venv", venv_dir],
        capture_output=True,
        text=True,
    )
    if result.returncode != 0:
        log(f"Ошибка создания venv: {result.stderr[:200]}", "error")
        return False

    log("Виртуальное окружение создано", "venv")
    return True


def install_dependencies(venv_dir: str, requirements_path: str) -> bool:
    """Устанавливает зависимости через pip."""
    pip_exe = os.path.join(venv_dir, "Scripts", "pip.exe")

    # — Пропускаем pip install --upgrade pip, на Windows он падает —
    #    свежий pip уже есть в venv "из коробки".

    # PyTorch CPU-only
    log("Установка PyTorch (CPU-only, ~800MB)...", "deps")
    result = subprocess.run(
        [pip_exe, "install", "torch", "torchaudio",
         "--index-url", "https://download.pytorch.org/whl/cpu",
         "--no-cache-dir"],
        capture_output=True, text=True,
    )
    if result.returncode != 0:
        log(f"PyTorch STDERR: {result.stderr.strip()[-500:]}", "error")
        return False
    out = result.stdout.strip()
    if out:
        for line in out.split("\n")[-3:]:
            line = line.strip()
            if line:
                log(f"  {line}", "deps")

    # Остальные зависимости: soundfile, numpy, misaki
    log("Установка soundfile, numpy, misaki...", "deps")
    if not run_pip(pip_exe, ["install", "-r", requirements_path], "Зависимости"):
        return False

    # Kokoro ставим без зависимостей — misaki уже есть,
    # а misaki[en] (тянет spacy с C++ сборкой) нам не нужен для японского
    log("Установка kokoro (без зависимостей)...", "deps")
    if not run_pip(pip_exe, ["install", "kokoro>=0.1.0", "--no-deps"], "Kokoro"):
        return False

    return True


def download_silero(venv_dir: str, models_dir: str) -> bool:
    """Скачивает Silero TTS модель."""
    python_exe = os.path.join(venv_dir, "Scripts", "python.exe")
    silero_dir = os.path.join(models_dir, "silero")
    model_path = os.path.join(silero_dir, "v4_ru.pt")

    os.makedirs(silero_dir, exist_ok=True)

    if os.path.exists(model_path):
        log("Silero: модель уже скачана", "models")
        return True

    log("Скачивание Silero TTS (~30MB)...", "models")

    # Экранируем обратные слеши Windows для вставки в Python-скрипт
    sd = silero_dir.replace("\\", "/")
    mp = model_path.replace("\\", "/")

    script = f"""
import torch
import os

os.makedirs('{sd}', exist_ok=True)

model, example_text = torch.hub.load(
    repo_or_dir='snakers4/silero-models',
    model='silero_tts',
    language='ru',
    speaker='v4_ru',
)
torch.save(model.state_dict(), '{mp}')
print('Silero: модель сохранена')
"""

    result = subprocess.run(
        [python_exe, "-c", script],
        capture_output=True,
        text=True,
    )
    if result.returncode != 0:
        log(f"Silero: {result.stderr[:200]}", "error")
        return False

    log("Silero: модель скачана", "models")
    return True


def download_kokoro(venv_dir: str, models_dir: str) -> bool:
    """Скачивает Kokoro TTS модель (кэшируется в HF cache)."""
    python_exe = os.path.join(venv_dir, "Scripts", "python.exe")

    log("Скачивание Kokoro TTS (~300-500MB, кэшируется HuggingFace)...", "models")

    script = """
from kokoro import KPipeline
pipeline = KPipeline(lang_code='j')
print('Kokoro: модель загружена')
generator = pipeline('テスト')
for _ in generator:
    break
print('Kokoro: тестовый синтез выполнен')
"""

    result = subprocess.run(
        [python_exe, "-c", script],
        capture_output=True,
        text=True,
    )
    if result.returncode != 0:
        log(f"Kokoro: {result.stderr[:200]}", "error")
        return False

    log("Kokoro: модель скачана", "models")
    return True


def get_default_dir(name: str) -> str:
    """Возвращает путь по умолчанию в папке пользователя."""
    appdata = os.environ.get("APPDATA") or os.path.expanduser("~/.config")
    return os.path.join(appdata, "Yappari", name)


def main():
    parser = argparse.ArgumentParser(description="Установка TTS-окружения Yappari")
    parser.add_argument(
        "--venv-dir", default=None,
        help="Путь к папке venv (по умолчанию: %%APPDATA%%/Yappari/tts_env)",
    )
    parser.add_argument(
        "--models-dir", default=None,
        help="Путь к папке с моделями (по умолчанию: %%APPDATA%%/Yappari/tts_models)",
    )
    args = parser.parse_args()

    if args.venv_dir is None:
        args.venv_dir = get_default_dir("tts_env")
    if args.models_dir is None:
        args.models_dir = get_default_dir("tts_models")

    log(f"📂 Venv: {args.venv_dir}", "venv")
    log(f"📂 Models: {args.models_dir}", "models")

    try:
        # Шаг 1: venv
        if not setup_venv(args.venv_dir):
            log_done(False, "Ошибка создания виртуального окружения")
            sys.exit(1)

        # Шаг 2: Зависимости
        script_dir = os.path.dirname(os.path.abspath(__file__))
        req_path = os.path.join(script_dir, "requirements-tts.txt")

        if not install_dependencies(args.venv_dir, req_path):
            log_done(False, "Ошибка установки зависимостей")
            sys.exit(1)

        # Шаг 3: Silero
        if not download_silero(args.venv_dir, args.models_dir):
            log_done(False, "Ошибка скачивания Silero")
            sys.exit(1)

        # Шаг 4: Kokoro
        if not download_kokoro(args.venv_dir, args.models_dir):
            log_done(False, "Ошибка скачивания Kokoro")
            sys.exit(1)

        log_done(True, "TTS-окружение готово! Модели установлены.")
        sys.exit(0)

    except Exception as e:
        log_done(False, str(e))
        sys.exit(1)


if __name__ == "__main__":
    main()
