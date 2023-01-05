package output

import (
	"io"
	"os"
	"os/exec"
	"sophiex/internal/logger"
	"sync"
)

type OsPath interface {
	Path() string
}

func CreateMuxer() *FFmpegMuxer {
	return &FFmpegMuxer{}
}

type FFmpegMuxer struct{}

func (muxer *FFmpegMuxer) WriteTo(inputs []OsPath, output io.WriteCloser, manager *sync.WaitGroup) {
	logger.Log.Debug("FfmpegMuxer: Writing to stdout")
	manager.Add(1)

	args := []string{"-y"}
	for _, input := range inputs {
		args = append(args, "-i", input.Path())
	}
	args = append(args, "-c:a", "aac", "-c:v", "copy", "-f", "matroska", "-")

	command := exec.Command("ffmpeg", args...)
	command.Stdout = output
	command.Stderr = os.Stderr

	go func() {
		logger.Log.Debug("%v\n", command.Args)

		err := command.Run()
		if err != nil {
			panic(err)
		}

		output.Close()
		manager.Done()
	}()
}

func (muxer *FFmpegMuxer) WriteToFile(inputs []OsPath, output OsPath, done chan bool) {
	logger.Log.Debug("FfmpegMuxer: Writing to file")

	args := []string{"-y"}
	for _, input := range inputs {
		args = append(args, "-i", input.Path())
	}
	args = append(args, "-c:a", "copy", "-c:v", "copy", output.Path())

	command := exec.Command("ffmpeg", args...)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	go func() {
		logger.Log.Debug("%v\n", command.Args)

		err := command.Run()
		if err != nil {
			panic(err)
		}

		close(done)
	}()
}
