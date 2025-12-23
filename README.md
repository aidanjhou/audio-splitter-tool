# Audio Splitter Tool

[English] | [简体中文](./README_zh.md)

A professional CLI tool designed to split audio files based on SRT subtitles and auto-merge them into high-quality WAV files. Optimized for speech recognition workflows (16k, Mono, S16LE).

---

## ✨ Features

* **Smart Environment Awareness**: Strictly distinguishes between Terminal launch (CLI mode) and File Explorer launch (Drag & Drop mode).
* **Silent Matching**: Automatically detects and uses the `.srt` file if it shares the same name as the audio in the same directory.
* **Professional Output**:
    * Individual clip file: `{input_audio_filename}_output/audio/{some_keyword}/{some_keyword}_{clip_index}.wav`
    * Merged keyword file: `{input_audio_filename}_output/merged/{some_keyword}_total_{clip_count}.wav`
    * Merged file: `{input_audio_filename}_output/merged/total_{clip_count}.wav`
* **Cross-Platform**: Native binaries for Windows (amd64/arm64), macOS (Intel/Apple Silicon), and Linux (amd64/arm64).

---

## 🚀 Installation

### 1. Prerequisites
This tool requires **FFmpeg** to be installed and available in your system's PATH.

* **Windows**: `winget install ffmpeg` or `scoop install ffmpeg`
* **macOS**: `brew install ffmpeg`
* **Linux**: `sudo apt install ffmpeg`

### 2. Download
Download the latest version for your system from the [Releases](../../releases) page.

---

## 🛠 Usage

### Mode A: Drag & Drop (Desktop)
Simply **drag** your audio file (WAV/MP3/etc.) and **drop** it onto the `audio-splitter-tool` executable icon.
* If a matching `.srt` is found, it processes automatically.
* The window will stay open after completion for you to review results.

### Mode B: Command Line (CLI)
```bash
# Basic usage
./audio-splitter-tool -a input.wav -t input.srt

# Use smart matching (if input.srt exists in the same folder)
./audio-splitter-tool -a input.wav

# Show version
./audio-splitter-tool -v
