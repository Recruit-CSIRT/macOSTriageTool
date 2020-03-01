package main

import (
	"github.com/Recruit-CSIRT/macOSTriageTool/pkg/conf"
	"github.com/Recruit-CSIRT/macOSTriageTool/pkg/gui"
)

func main() {
	var config conf.Config
	config.SetCurrentDir()
	gui.Run(&config)
}
