# 音频切割合并工具 (audio-splitter-tool)

[English](./README.md) | [简体中文]

基于 SRT 字幕切割音频并自动合并的专业命令行工具。专为语音识别（16k 采样率、单声道、S16LE 编码）任务优化。

---

## ✨ 核心功能

* **环境感知**: 严格区分终端运行（完成后直接退出）与文件夹运行（完成后停住窗口，防止闪退）。
* **智能匹配**: 自动识别同级目录下同名的 `.srt` 文件，实现“一键切割”，无需手动指定字幕路径。
* **规范输出**:
    * 独立切片存至 `{输入音频文件名}_output/audio/`，命名格式为 `{关键词}/{关键词}_{片段序号}.wav`。
    * 关键词合并文件存至 `{输入音频文件名}_output/merged/`，命名格式为 `{关键词}_total_{片段数}.wav`。
    * 总合并文件存至 `{输入音频文件名}_output/merged/`，命名格式为 `total_{片段数}.wav`。
* **全平台原生支持**: 提供 Windows (amd64/arm64), macOS (Intel/M1/M2), Linux (amd64/arm64) 的原生二进制包。

---

## 🚀 安装指南

### 1. 前提条件
本工具依赖 **FFmpeg**。请确保其已安装并已添加到系统环境变量。

* **Windows**: 使用 `winget install ffmpeg` (系统自带) 或 `scoop install ffmpeg`。
* **macOS**: `brew install ffmpeg`
* **Linux**: `sudo apt install ffmpeg`

### 2. 下载工具
请前往 [Releases](../../releases) 页面下载适合您系统架构的最新版本。

---

## 🛠 使用方法

### 模式 A：直接拖拽（推荐）
只需将音频文件（WAV/MP3 等）**拖拽**到 `audio-splitter-tool` 的可执行文件图标上即可。
* 如果同级目录下有同名 SRT，程序将自动开始处理。
* 处理完成后窗口会保持打开状态，方便查看输出结果。

### 模式 B：命令行
```bash
# 基础用法
./audio-splitter-tool -a input.wav -t input.srt

# 使用智能匹配（如果 input.srt 在同一文件夹下）
./audio-splitter-tool -a input.wav

# 查看版本
./audio-splitter-tool -v
