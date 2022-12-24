package output

import (
	"bytes"
	"fmt"
	"os/exec"
	"sophiex/internal/logger"
)

func CreateMuxer() StreamWriter {
	return &FFmpegMuxer{}
}

type FFmpegMuxer struct{}

func (muxer *FFmpegMuxer) Launch(streams []*StreamDownloader, outputPath string, done chan bool) {
	args := []string{"-y"}
	for _, stream := range streams {
		args = append(args, "-i", stream.Output.Path)
	}
	args = append(args, "-c", "copy", outputPath)
	//args = append(args, "-f", "matroska", outputPath)

	command := exec.Command("ffmpeg", args...)

	var outb, errb bytes.Buffer
	command.Stdout = &outb
	command.Stderr = &errb

	go func() {
		logger.Log.Debug("%v\n", command.Args)
		err := command.Run()
		if err != nil {
			logger.Log.Debug("out:", outb.String(), "err:", errb.String())
			panic(err)
		}
		close(done)
	}()
}
