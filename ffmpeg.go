package runal

import (
	"fmt"
	"os/exec"
	"strconv"
)

const (
	ffmpegBinary = "ffmpeg"
)

func checkFFMPEG() bool {
	_, err := exec.LookPath(ffmpegBinary)
	return err == nil
}

func framesToMP4Videos(fps int, input, output string) error {
	cmd := exec.Command(ffmpegBinary, buildArgs(fps, input, output)...)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %s", err.Error())
	}
	return nil
}

func buildArgs(fps int, input, output string) []string {
	return []string{
		"-framerate",
		strconv.Itoa(fps),
		"-y",
		"-i",
		input,
		"-c:v",
		"libx264",
		"-pix_fmt",
		"yuv420p",
		output,
	}
}
