package output

import (
	"io"
	"os"
	"os/exec"
	"sophiex/internal/logger"
)

func CreateMuxer() *FFmpegMuxer {
	return &FFmpegMuxer{}
}

type FFmpegMuxer struct{}

func (muxer *FFmpegMuxer) WriteTo(streams []*StreamDownloader, output io.Writer, done chan bool) {
	logger.Log.Debug("FfmpegMuxer: Writing to stdout")

	args := []string{"-y"}
	for _, stream := range streams {
		args = append(args, "-i", stream.Output.Path)
	}
	args = append(args, "-c:a", "aac", "-c:v", "copy", "-f", "matroska", "-")

	command := exec.Command("ffmpeg", args...)
	command.Stdout = output

	go func() {
		logger.Log.Debug("%v\n", command.Args)

		err := command.Run()
		if err != nil {
			panic(err)
		}

		close(done)
	}()
}

func (muxer *FFmpegMuxer) WriteToFile(streams []*StreamDownloader, outputPath string, done chan bool) {
	logger.Log.Debug("FfmpegMuxer: Writing to file")

	args := []string{"-y"}
	for _, stream := range streams {
		args = append(args, "-i", stream.Output.Path)
	}
	args = append(args, "-c", "copy", outputPath)

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
