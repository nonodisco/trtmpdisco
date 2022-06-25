package processor

import (
	"net/url"
	"os"
	"os/exec"
	"path"
)

var name string

func getName() string {
	if name != "" {
		return name
	}
	name = os.Getenv("FFMPEG_PATH")
	if name == "" {
		name = "ffmpeg"
	}
	return name
}

var output string

func getOutput() (string, error) {
	if output != "" {
		return output, nil
	}
	key := os.Getenv("RTMP_KEY")
	url, err := url.Parse(os.Getenv("RTMP_URL"))
	if err != nil {
		return "", err
	}
	url.Path = path.Join(url.Path, key)
	output = url.String()
	return output, nil
}

var (
	args = []string{
                "-re",
		"-stream_loop", "-1",
		"-f", "image2",
		"-i", "radio-background.jpg",
		"-i", "",
		"-c:a", "libfdk_aac",
		"-c:v", "libx264",
		"-vf", "scale=400:-2",
		"-preset", "ultrafast",
		"-b:v", "2500k",
		"-tune", "stillimage",
		"-ac", "2",
		"-ar", "44100",
		"-pix_fmt", "yuv420p",
		"-g", "60",
		"-f", "flv",
	}
	inputIndex = 7
)

func getArgs(input string) ([]string, error) {
	toReturn := args
	toReturn[inputIndex] = input
	output, err := getOutput()
	if err != nil {
		return toReturn, err
	}
	toReturn = append(toReturn, output)
	return toReturn, nil
}

var process *exec.Cmd

func Stop() (bool, error) {
	if Processing() {
		err := process.Process.Kill()
		process = nil
		return true, err
	}
	return false, nil
}

func Process(input string, errc chan error) error {
	if Processing() {
		if _, err := Stop(); err != nil {
			return err
		}
	}
	args, err := getArgs(input)
	if err != nil {
		return err
	}
	process = exec.Command(getName(), args...)
	process.Stderr = os.Stderr
	err = process.Start()
	if err == nil {
		go func() {
			err = process.Wait()
			if err != nil {
				errc <- err
				close(errc)
			}
			Stop()
		}()
	}
	return err
}

func Processing() bool {
	return process != nil
}
