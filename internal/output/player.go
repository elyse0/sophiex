package output

import (
	"io"
	"os/exec"
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
