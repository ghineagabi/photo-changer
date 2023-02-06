package main

import (
	"fmt"
	"log"
	"os/exec"
	"testing"
)

func TestGeneratebands(t *testing.T) {
	width := 7
	diagLen := 6
	bandLen := 40
	bandStart := 40
	smoothness := 16
	bands := generateBand(width, diagLen, bandLen, bandStart, smoothness)
	for _, val := range *bands {
		fmt.Println(val)
	}
}

func TestFfmpeg(t *testing.T) {
	cmd := exec.Command("cmd", "ffmpeg", "-framerate", "10", "-i", "converted/%03d.jpeg", "-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2", "output.mp4")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
