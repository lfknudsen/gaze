package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

func closePipe(r *io.PipeReader, w *io.PipeWriter) {
	_ = r.Close()
	_ = w.Close()
}

type VideoInfo struct {
	Width, Height, Bitrate int64
	Framerate              float64
}

func probe(path string) *VideoInfo {
	cmd := exec.Command("ffprobe", "-v", "0", "-of", "default=nk=1:nw=1",
		"-select_streams", "v", "-show_entries", "stream=width,height,r_frame_rate",
		"-show_entries", "format=bit_rate", path)
	pipeReader, pipeWriter := io.Pipe()
	cmd.Stdout = pipeWriter
	reader := bufio.NewReader(pipeReader)
	err := cmd.Start()
	if err != nil {
		closePipe(pipeReader, pipeWriter)
		log.Fatal(err)
	}
	widthStr, _ := reader.ReadString('\n')
	width, err := strconv.ParseInt(strings.TrimSpace(widthStr), 10, 64)

	heightStr, _ := reader.ReadString('\n')
	height, err := strconv.ParseInt(strings.TrimSpace(heightStr), 10, 64)

	framerateStr, _ := reader.ReadString('\n')
	aboveStr, belowStr, _ := strings.Cut(framerateStr, "/")
	above, _ := strconv.ParseFloat(strings.TrimSpace(aboveStr), 64)
	below, _ := strconv.ParseFloat(strings.TrimSpace(belowStr), 64)
	framerate := above / below

	bitrateStr, _ := reader.ReadString('\n')
	bitrate, err := strconv.ParseInt(strings.TrimSpace(bitrateStr), 10, 64)
	closePipe(pipeReader, pipeWriter)
	_ = cmd.Wait()
	info := VideoInfo{
		Width:     width,
		Height:    height,
		Bitrate:   bitrate,
		Framerate: framerate,
	}
	return &info
}

func getBitrate(path string) string {
	cmd := exec.Command("ffprobe", "-v", "0", "-of",
		"default=nk=1:nw=1", "-show_entries", "format=bit_rate", path)
	pipeReader, pipeWriter := io.Pipe()
	cmd.Stdout = pipeWriter
	reader := bufio.NewReader(pipeReader)
	err := cmd.Start()
	if err != nil {
		closePipe(pipeReader, pipeWriter)
		log.Fatal(err)
	}
	capturedOutput, err := reader.ReadString('\n')
	capturedOutput = strings.TrimSpace(capturedOutput)
	closePipe(pipeReader, pipeWriter)
	_ = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
	return capturedOutput
}

func main() {
	dir, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	maxLength := 0
	for _, file := range dir {
		length := len(file.Name())
		if maxLength < length {
			maxLength = length
		}
	}
	maxLength += 3
	var p = message.NewPrinter(language.English)
	fmt.Printf("NAME %s| BITRATE | WIDTH | HEIGHT | FPS\n", strings.Repeat(" ", maxLength-len("NAME")))
	for _, entry := range dir {
		if strings.HasSuffix(entry.Name(), ".mp4") {
			info := probe(entry.Name())
			_, _ = p.Printf("%s %s %v   %d    %d    %f\n", entry.Name(),
				strings.Repeat(" ", maxLength-len(entry.Name())),
				number.Decimal(info.Bitrate), info.Width, info.Height, info.Framerate)
		}
	}
}
