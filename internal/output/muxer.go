package output

import (
	"bytes"
	"fmt"
	"os/exec"
)

func CreateMuxer() StreamWriter {
	return &FFmpegMuxer{}
}

type FFmpegMuxer struct{}

func (muxer *FFmpegMuxer) Launch(streams []*StreamDownloader, outputPath string) {
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
		fmt.Printf("%v\n", command.Args)
		err := command.Run()
		if err != nil {
			fmt.Println("out:", outb.String(), "err:", errb.String())
			panic(err)
		}
	}()
}
