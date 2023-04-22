package yolov3

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Params struct {
	InputSize [2]int
	ConfigPath string
}

func (p Params) GetConfigPath() string {
	if p.ConfigPath == "" {
		return "cfg/yolov3.cfg"
	} else {
		return p.ConfigPath
	}
}

func CreateParams(fname string, p Params, training bool) {
	// prepare configuration with this width/height
	configPath := p.GetConfigPath()
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join("lib/darknet/", configPath)
	}
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	for _, line := range strings.Split(string(bytes), "\n") {
		line = strings.TrimSpa