package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// The version variable is defined at the package level.
// It will be overwritten by -ldflags "-X 'main.version=$VERSION'" during build.
var version = "dev-build"

func main() {
	// 1. Environment detection
	independent := isIndependentWindow()

	// 2. Flags definition and Usage override
	flag.Usage = func() {
		fullExeName := filepath.Base(os.Args[0])
		exeName := strings.TrimSuffix(fullExeName, filepath.Ext(fullExeName))
		fmt.Printf("Audio Splitter Tool %s\n", version)
		fmt.Printf("Usage: %s [-a <audio>] [-t <srt>] [-v] [-h]\n", exeName)
		if independent {
			fmt.Println("\nPress Enter to exit...")
			pause()
		}
		os.Exit(0)
	}

	audioFlag := flag.String("a", "", "Input audio file (.wav)")
	srtFlag := flag.String("t", "", "Input SRT file")
	showVer := flag.Bool("v", false, "Show version")
	flag.Parse()

	if *showVer {
		fmt.Printf("Version: %s\n", version)
		if independent {
			pause()
		}
		return
	}

	var audioPath, srtPath string
	audioPath = *audioFlag
	srtPath = *srtFlag

	// 3. Drag-and-drop support for independent window mode
	// Fix: Check flag.Args() instead of os.Args directly because flag.Parse() consumes flags
	if audioPath == "" && flag.NArg() > 0 {
		firstArg := strings.Trim(flag.Arg(0), ` "`)
		if isValidFile(firstArg) {
			audioPath = firstArg
		}
	}

	// 4. Input validation loop for Audio
	for audioPath == "" || !isValidFile(audioPath) {
		if audioPath != "" {
			fmt.Printf("Invalid file: '%s'\n", audioPath)
		}
		fmt.Print("Please input Audio file path: ")
		audioPath = readLine()
	}

	// 5. Input validation loop for SRT (with auto-match)
	if srtPath == "" || !isValidFile(srtPath) {
		base := strings.TrimSuffix(audioPath, filepath.Ext(audioPath))
		potentialSrt := base + ".srt"
		if isValidFile(potentialSrt) {
			srtPath = potentialSrt
			fmt.Printf("Auto-matched SRT: %s\n", filepath.Base(srtPath))
		} else {
			for srtPath == "" || !isValidFile(srtPath) {
				fmt.Print("Please input SRT file path: ")
				srtPath = readLine()
			}
		}
	}

	// 6. FFmpeg dependency check
	ffmpeg, err := findFFmpeg()
	if err != nil {
		fmt.Printf("\nFATAL ERROR: %v\n", err)
		if independent {
			pause()
		}
		os.Exit(1)
	}

	// 7. Execute processing logic
	doWork(ffmpeg, audioPath, srtPath, independent)
}

type Segment struct {
	Start, End, Content string
}

func doWork(ffmpeg, audioPath, srtPath string, independent bool) {
	absPath, _ := filepath.Abs(audioPath)
	baseDir := filepath.Dir(absPath)
	fileName := strings.TrimSuffix(filepath.Base(absPath), filepath.Ext(absPath))
	outputRoot := filepath.Join(baseDir, fileName + "_output")
	audioDir := filepath.Join(outputRoot, "audio")
	totalDir := filepath.Join(outputRoot, "merged")

	_ = os.MkdirAll(audioDir, 0755)
	_ = os.MkdirAll(totalDir, 0755)

	pkgName := calcPkgName()
	appTempDir := filepath.Join(os.TempDir(), pkgName)
	_ = os.MkdirAll(appTempDir, 0755)

	silentWav := filepath.Join(appTempDir, "silent_5s.wav")
	_ = exec.Command(ffmpeg, "-y", "-f", "lavfi", "-i", "anullsrc=r=16000:cl=mono:d=5", 
		"-c:a", "pcm_s16le", silentWav).Run()
	defer os.Remove(silentWav)

	segments := parseSRT(srtPath)
	keywordPaths := make(map[string][]string)
	var allFiles []string
	keywordCounter := make(map[string]int)

	fmt.Println("\n[1/3] Processing segments by time sequence...")
	for _, seg := range segments {
		kw := seg.Content

		kwSubDir := filepath.Join(audioDir, kw)
		_ = os.MkdirAll(kwSubDir, 0755)

		keywordCounter[kw]++
		currentIndex := keywordCounter[kw]

		sliceFileName := fmt.Sprintf("%s_%06d.wav", kw, currentIndex)
		targetPath := filepath.Join(kwSubDir, sliceFileName)

		// Slice and resample: 16k, 16bit, mono
		cmd := exec.Command(ffmpeg, "-y",
			"-ss", seg.Start, "-to", seg.End,
			"-i", audioPath,
			"-ar", "16000", "-ac", "1", "-sample_fmt", "s16",
			targetPath)
		
		if err := cmd.Run(); err == nil {
			keywordPaths[kw] = append(keywordPaths[kw], targetPath)
			allFiles = append(allFiles, targetPath)
			fmt.Printf("  -> Created: %s\n", sliceFileName)
		} else {
			fmt.Printf("  -> Failed to create: %s\n", sliceFileName)
		}
	}

	// Keyword-based merge
	fmt.Println("\n[2/3] Merging per keyword groups...")
	for kw, files := range keywordPaths {
		outName := fmt.Sprintf("%s_total_%d.wav", kw, len(files))
		mergeFiles(ffmpeg, files, silentWav, filepath.Join(totalDir, outName))
	}

	// Global sequential merge
	fmt.Println("\n[3/3] Final total merge...")
	finalName := fmt.Sprintf("total_%d.wav", len(allFiles))
	mergeFiles(ffmpeg, allFiles, silentWav, filepath.Join(totalDir, finalName))

	fmt.Printf("\nProcessing Complete. Output: %s\n", outputRoot)
	if independent {
		fmt.Println("Press Enter to exit.")
		pause()
	}
}

func mergeFiles(ffmpeg string, files []string, silentWav string, outPath string) {
	if len(files) == 0 { return }
	listPath := outPath + ".txt"
	f, _ := os.Create(listPath)
	for i, file := range files {
		abs, _ := filepath.Abs(file)
		fmt.Fprintf(f, "file '%s'\n", strings.ReplaceAll(abs, "'", "'\\''"))

		if i < len(files)-1 {
			absSilent, _ := filepath.Abs(silentWav)
			fmt.Fprintf(f, "file '%s'\n", strings.ReplaceAll(absSilent, "'", "'\\''"))
		}
	}
	f.Close()
	defer os.Remove(listPath)

	cmd := exec.Command(ffmpeg, "-y", "-f", "concat", "-safe", "0", "-i", listPath, "-c", "copy", outPath)
	_ = cmd.Run()
}

func parseSRT(path string) []Segment {
	var segs []Segment
	file, _ := os.Open(path)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	step := 0
	curr := Segment{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" { step = 0; continue }
		switch step {
		case 0: step = 1
		case 1:
			times := strings.Split(line, " --> ")
			if len(times) == 2 {
				curr.Start = strings.Replace(times[0], ",", ".", 1)
				curr.End = strings.Replace(times[1], ",", ".", 1)
				step = 2
			}
		case 2:
			curr.Content = line
			segs = append(segs, curr)
			step = 0
		}
	}
	return segs
}

func calcPkgName() string {
	fullExeName := filepath.Base(os.Args[0])
	exeName := strings.TrimSuffix(fullExeName, filepath.Ext(fullExeName))
	suffixes := []string{"_amd64", "_arm64"}
	pkgName := exeName
	for _, s := range suffixes {
		if strings.HasSuffix(pkgName, s) {
			pkgName = pkgName[:len(pkgName)-len(s)]
			break
		}
	}
	return pkgName
}

func isValidFile(p string) bool { i, e := os.Stat(p); return e == nil && !i.IsDir() }
func readLine() string { s := bufio.NewScanner(os.Stdin); if s.Scan() { return strings.Trim(s.Text(), ` "`) }; return "" }
func pause() { bufio.NewReader(os.Stdin).ReadBytes('\n') }
func findFFmpeg() (string, error) {
	exe, _ := os.Executable()
	local := filepath.Join(filepath.Dir(exe), "ffmpeg.exe")
	if _, err := os.Stat(local); err == nil { return local, nil }
	p, err := exec.LookPath("ffmpeg")
	if err == nil { return p, nil }
	return "", fmt.Errorf("ffmpeg not found")
}
