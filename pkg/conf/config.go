package conf

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/Recruit-CSIRT/macOSTriageTool/pkg/log"
)

type Config struct {

	// Option
	IsEnabledToGetProfiler bool
	IsEnabledToCalcHash    bool
	IsEnabledToGetStatInfo bool
	IsEnabledToSaveIntoDmg bool

	// Path
	RootPath       string
	OutputRootPath string
	OutputDirPath  string
	CurrentDirPath string

	UserDefinedFileListName string

	SelectedPresetNum int

	DmgSize int64
}

func (c *Config) SetCurrentDir(){
	dir, err := os.Getwd()
	if err != nil {
		log.ToolLogger.Error("[-] Failed to get a current directory.")
		log.ToolLogger.Error("[-] Change the output directory to /tmp.")
		dir = TmpDir
	}

	if dir == "/" {
		dir = TmpDir
	}

	c.CurrentDirPath = dir
}

func (c *Config) Setting() error {

	c.SetCurrentDir()

	c.OutputRootPath = filepath.Join(c.OutputRootPath, OutputDirName)
	err := os.MkdirAll(c.OutputRootPath, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}


func (c *Config) CreateAndMountDmg() error{

	dmgSize := float64(c.DmgSize / 1024 ) * 1.2

	// hdiutil create -fs APFS -size xxxk -volname Disk /Volumes/Name/Disk.dmg
	cmd := exec.Command(CmdHdiutil, "create",
		"-fs", DmgFileSystem,
		"-size", strconv.FormatInt(int64(dmgSize), 10) + "k",
		"-volname", VolumeName,
		"-nospotlight",
		"-attach",
		filepath.Join(c.OutputRootPath, VolumeName+".dmg"))
	_, err := cmd.Output()
	if err != nil {
		log.ToolLogger.Error(err.Error())
		log.ToolLogger.Error("Failed to create dmg file.")
		return err
	}

	c.OutputDirPath = filepath.Join(VolumePath, VolumeName)

	return nil
}

func (c *Config) UnmountDmg () error {
	// hdiutil unmount /Volumes/Disk
	cmd := exec.Command(CmdHdiutil, "unmount", c.OutputDirPath)

	_, err := cmd.Output()
	if err != nil {
		log.ToolLogger.Error("Failed to unmount dmg.")
		return err
	}

	return nil
}

func (c *Config) MakeOutputDir() error {
	c.OutputDirPath = filepath.Join(c.OutputRootPath, VolumeName)
	err := os.MkdirAll(c.OutputDirPath, os.ModePerm)

	if err != nil {
		log.ToolLogger.Error("Failed to make output dir.")
		return err
	}

	return nil
}
