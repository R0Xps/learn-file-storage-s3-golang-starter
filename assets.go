package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

type videoFileStreams struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

func getVideoAspectRatio(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	buf := bytes.NewBuffer(nil)
	cmd.Stdout = buf

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	decoder := json.NewDecoder(buf)
	streams := videoFileStreams{}
	err = decoder.Decode(&streams)
	if err != nil {
		return "", err
	}

	width := streams.Streams[0].Width
	height := streams.Streams[0].Height

	if width/16 == height/9 {
		return "16:9", nil
	}

	if width/9 == height/16 {
		return "9:16", nil
	}

	return "other", nil
}
