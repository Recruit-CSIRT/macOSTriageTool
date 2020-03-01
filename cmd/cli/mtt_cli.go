package main

import (
	"flag"
	"fmt"

	"github.com/Recruit-CSIRT/macOSTriageTool/pkg/log"
	"github.com/Recruit-CSIRT/macOSTriageTool/pkg/triage"
	"github.com/Recruit-CSIRT/macOSTriageTool/pkg/conf"
)

var config conf.Config

func init() {

	flag.StringVar(&config.RootPath, "root", "/", "set the evidence root path.")
	flag.StringVar(&config.OutputRootPath, "output", ".", "set the output path.")

	flag.IntVar(&config.SelectedPresetNum, "preset", 0, "select forensics type: 0(AllList), 1(Malware), 2(Fraud), 3(macripper), 4(only custom). (default: 0)")
	flag.StringVar(&config.UserDefinedFileListName, "file", "", "set user custom file list path.")

	flag.BoolVar(&config.IsEnabledToGetProfiler, "i", false, "get the system information. (default: false)")
	flag.BoolVar(&config.IsEnabledToCalcHash, "c", false, "calc the file hash. (default: false)")
	flag.BoolVar(&config.IsEnabledToGetStatInfo, "s", false, "get the stat info. (default: false)")
	flag.BoolVar(&config.IsEnabledToSaveIntoDmg, "d", false, "save files into a dmg. (default: false)")

}

func main() {

	flag.Parse()

	err := config.Setting()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = log.Init(config.OutputRootPath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = triage.Triage(&config)
	if err != nil {
		log.ToolLogger.Error(err.Error())
		return
	}

}
