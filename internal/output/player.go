package output

import (
	"fmt"
	"os/exec"
)

type VlcPlayer struct {
}

func CreateStreamPlayer() StreamPlayer {
	return &VlcPlayer{}
}

func (player *VlcPlayer) Launch(streams []*StreamDownloader) {
	var command *exec.Cmd
	if len(streams) == 1 {
		command = exec.Command("vlc", streams[0].Output.Path)
		go func() {
			err := command.Run()
			if err != nil {
				fmt.Print(err)
			}
		}()

	} else {
		muxerOutput := CreateNamedPipe()
		muxerOutput.Open()
		defer muxerOutput.Close()

		//command = exec.Command("vlc", muxerOutput.Path)

		//go func() {
		//	err := command.Run()
		//	if err != nil {
		//		fmt.Print(err)
		//	}
		//}()

		muxer := CreateMuxer()
		muxer.Launch(streams, muxerOutput.Path)
	}
}
