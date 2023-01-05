package output

import (
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type MpvPlayer struct {
	bin string
}

func CreateStreamPlayer() *MpvPlayer {
	return &MpvPlayer{
		bin: "mpv",
	}
}

func (player *MpvPlayer) PlayFrom(input io.Reader) {
	playerCommand := exec.Command(player.bin, "-", "--player-operation-mode=pseudo-gui")
	playerCommand.Stdin = input

	go func() {
		err := playerCommand.Run()
		if err != nil {
			panic(err)
		}
	}()
}

func (player *MpvPlayer) Launch(streams []*StreamDownloader, manager *sync.WaitGroup) {
	var playerCommand *exec.Cmd
	if len(streams) == 1 {
		// playerCommand = exec.Command(player.bin, streams[0].Output.Path)
		playerCommand = exec.Command(player.bin, "FIXME")

		go func() {
			err := playerCommand.Run()
			if err != nil {
				fmt.Print(err)
			}
		}()

		manager.Done()
		return
	}

	playerCommand = exec.Command(player.bin, "-")
	playerStdin, _ := playerCommand.StdinPipe()

	muxer := CreateMuxer()
	// FIXME
	muxer.WriteTo(nil, playerStdin, manager)

	go func() {
		err := playerCommand.Run()
		if err != nil {
			fmt.Print(err)
		}
	}()
}
