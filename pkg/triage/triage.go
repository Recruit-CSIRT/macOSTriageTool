	package triage

	import (
		"bytes"
		"fmt"
		"os"
		"os/exec"
		"path/filepath"
		"strings"

		"github.com/Recruit-CSIRT/macOSTriageTool/pkg/conf"
		"github.com/Recruit-CSIRT/macOSTriageTool/pkg/log"
		"github.com/Recruit-CSIRT/macOSTriageTool/pkg/utils"
	)

func Triage(config *conf.Config) error {

	log.ToolLogger.Info("[+] Tool Start")

	triageList := utils.TriageFileList{}
	err := triageList.New(config)
	if err != nil {
		return err
	}

	fsnfl, err := os.OpenFile(filepath.Join(config.OutputRootPath, conf.ScheduledTriageFileListName), os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.ToolLogger.Error("[-] Failed to open the file: " + err.Error())
		return err
	}
	defer fsnfl.Close()

	_, err = fsnfl.WriteString(strings.Join(triageList.NormalizedFileList, "\n"))
	if err != nil {
		log.ToolLogger.Error("[-] Failed to write the list: " + err.Error())
		return err
	}


	if config.IsEnabledToCalcHash || config.IsEnabledToGetStatInfo || config.IsEnabledToSaveIntoDmg {
		log.ToolLogger.Info("[+] Start processing the options")

		utils.WriteColumn()

		for _, f := range triageList.NormalizedFileList {
			err := processFileOptions(&f, config)
			if err != nil {
				log.FileInfoErrorLogger.Error(err.Error())
			}
		}

		log.ToolLogger.Info("[+] Finish processing the options")
	}

	// create dir or dmg for file
	log.ToolLogger.Info("[+] Creating the output dir or dmg")
	if config.IsEnabledToSaveIntoDmg {
		err = config.CreateAndMountDmg()
		defer config.UnmountDmg()
	} else {
		err = config.MakeOutputDir()
	}
	if err != nil {
		return nil
	}


	log.ToolLogger.Info("[+] Start copying files")
	for _, f := range triageList.NormalizedFileList {
		copyFile(&f, config)
	}
	log.ToolLogger.Info("[+] Finish copying files")


	if config.IsEnabledToGetProfiler {
		log.ToolLogger.Info("[+] Start getting system profiler")
		err := getSystemProfiler(config)
		if err != nil {
			log.ToolLogger.Info("[-] Failed to get system profiler")
		}
		log.ToolLogger.Info("[+] Finish getting system profiler")
	}

	log.ToolLogger.Info("[+] Tool Finish")

	return nil
}

func copyFile(fileFrom *string, config *conf.Config) {

	// ditto --rsrc --extattr --qtn --acl fileFrom fileTo
	cmd := exec.Command(conf.CmdDitto, "--rsrc", "--extattr", "--qtn", "--acl", *fileFrom, filepath.Join(config.OutputDirPath, *fileFrom))

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.FileCopyLogger.Error("[-] Failed to copy: " + *fileFrom)
		log.FileCopyLogger.Error("[-] Message: " + strings.TrimRight(stderr.String(), "\n"))
	} else {
		log.FileCopyLogger.Info("[+] Copied: " + *fileFrom)
	}

}

func processFileOptions(filePath *string, config *conf.Config) error {

	err := filepath.Walk(*filePath, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			log.FileInfoErrorLogger.Error("[-] Prevent panic by handling failure accessing a path. Path: " + path + ", Error: " + err.Error())
			return err
		}

		// stat
		customFileInfo := utils.NewFileInfo(&path)
		if (config.IsEnabledToGetStatInfo || config.IsEnabledToSaveIntoDmg) && info != nil {
			customFileInfo.SetStat(&info)
		}

		// hash
		if config.IsEnabledToCalcHash && !info.IsDir()  {
			_ = customFileInfo.SetHash()
		}

		if config.IsEnabledToCalcHash || config.IsEnabledToGetStatInfo {
			customFileInfo.WriteLogger()
		}

		// calc size of all files
		if config.IsEnabledToSaveIntoDmg {
			config.DmgSize += int64(info.Size())
		}

		return nil
	})
	if err != nil {
		log.FileInfoErrorLogger.Error("[-] Error walking the path" + err.Error())
		return err
	}

	return nil
}

func getSystemProfiler(config *conf.Config) error {

	cmd := exec.Command(conf.CmdSystemProfiler)

	output, err := cmd.Output()
		if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(config.OutputRootPath, conf.SystemProfilerFileName))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, string(output))
	if err != nil {
		return err
	}

	return nil
}
