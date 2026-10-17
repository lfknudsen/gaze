package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

var concurrency = true

func main() {
	if len(os.Args) == 2 {
		concurrency = false
	}

	dir, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	names := make([]string, 0, len(dir))
	for _, entry := range dir {
		if strings.HasSuffix(entry.Name(), ".mp4") {
			names = append(names, entry.Name())
		}
	}
	infoList := make([]*VideoInfo, len(names))
	var wg *sync.WaitGroup
	if concurrency {
		wg = new(sync.WaitGroup)
	}
	for i := range names {
		if concurrency {
			wg.Go(func() {
				infoList[i] = probe(names[i])
			})
		} else {
			infoList[i] = probe(names[i])
		}
	}
	if concurrency {
		wg.Wait()
	}
	tabber := tabwriter.NewWriter(os.Stdout, 2, 3, 2, ' ', tabwriter.AlignRight)
	_, err = fmt.Fprint(tabber, "BITRATE\tWIDTH\tHEIGHT\tFPS\t\tNAME\n")
	if err != nil {
		log.Fatal(err)
	}
	var p = message.NewPrinter(language.English)
	for i, info := range infoList {
		str := p.Sprintf("%d\t%d\t%d\t%d\t\t%s\n",
			info.Bitrate, info.Width, info.Height, int64(math.Round(info.Framerate)), names[i])
		_, err = fmt.Fprint(tabber, str)
		if err != nil {
			log.Fatal(err)
		}
	}
	_ = tabber.Flush()
}
