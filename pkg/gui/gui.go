package gui

import (
	"os"

	"github.com/Recruit-CSIRT/macOSTriageTool/pkg/triage"
	"github.com/Recruit-CSIRT/macOSTriageTool/pkg/conf"
	"github.com/Recruit-CSIRT/macOSTriageTool/pkg/log"
	"github.com/therecipe/qt/widgets"
)

func Run(config *conf.Config){
	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := createWindow(config)

	window.Show()

	app.Exec()
}

func createWindow(config *conf.Config) *widgets.QMainWindow {

	// gui group of the evidence root path
	var (
		rootPathGroup    = widgets.NewQGroupBox2("Root path setting", nil)
		rootPathLineEdit = widgets.NewQLineEdit2("/", nil)
		rootPathButton   = widgets.NewQPushButton2("...", nil)
	)
	rootPathButton.ConnectClicked(func(bool) {
		openFileDialog(rootPathLineEdit.Text(), rootPathLineEdit)
	})

	// gui group of the output path
	var (
		outputPathGroup    = widgets.NewQGroupBox2("Output path setting", nil)
		outputPathLineEdit = widgets.NewQLineEdit2(config.CurrentDirPath, nil)
		outputPathButton   = widgets.NewQPushButton2("...", nil)
	)
	outputPathButton.ConnectClicked(func(bool) {
		openFileDialog(outputPathLineEdit.Text(), outputPathLineEdit)
	})

	// gui group of preset and custom file list
	var (
		presetPathLabel = widgets.NewQLabel2("Preset: ", nil, 0)
		presetComboBox  = widgets.NewQComboBox(nil)

		filelistPathGroup    = widgets.NewQGroupBox2("Preset setting", nil)
		filelistPathLabel    = widgets.NewQLabel2("Select your file list: ", nil, 0)
		filelistPathButton   = widgets.NewQPushButton2("...", nil)
		filelistPathLineEdit = widgets.NewQLineEdit2("", nil)

	)

	filelistPathButton.ConnectClicked(func(bool) {
		openFileDialog(filelistPathLineEdit.Text(), filelistPathLineEdit)
	})
	presetComboBox.AddItems(conf.PresetString)

	// gui group of option
	var (
		optionsGroup            = widgets.NewQGroupBox2("Options", nil)
		optionsCheckBoxProfiler = widgets.NewQCheckBox2("Get system information", nil)
		optionsCheckBoxFileHash = widgets.NewQCheckBox2("Get file hash", nil)
		optionsCheckBoxStat     = widgets.NewQCheckBox2("Get stat info of file", nil)
		optionsCheckBoxDmg     	= widgets.NewQCheckBox2("Save file into dmg", nil)
	)

	// run button and action
	runButton := widgets.NewQPushButton2("Run", nil)
	runButton.ConnectClicked(func(bool) {

		// get fields value
		config.RootPath = rootPathLineEdit.Text()
		config.OutputRootPath = outputPathLineEdit.Text()

		config.SelectedPresetNum = presetComboBox.CurrentIndex()
		config.UserDefinedFileListName = filelistPathLineEdit.Text()

		// get options value
		config.IsEnabledToGetProfiler = optionsCheckBoxProfiler.IsChecked()
		config.IsEnabledToCalcHash = optionsCheckBoxFileHash.IsChecked()
		config.IsEnabledToGetStatInfo = optionsCheckBoxStat.IsChecked()
		config.IsEnabledToSaveIntoDmg = optionsCheckBoxDmg.IsChecked()


		err := config.Setting()
		if err != nil {
			widgets.QMessageBox_Information(nil, "Status", "Failed to setup", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			return
		}

		err = log.Init(config.OutputRootPath)
		if err != nil {
			log.ToolLogger.Error(err.Error())
			return
		}

		err = triage.Triage(config)
		if err != nil {
			message := "Failed to triage: " + err.Error()
			widgets.QMessageBox_Information(nil, "Status", message, widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		}

		widgets.QMessageBox_Information(nil, "Status", "Finished!", widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})

	// layout setting
	// layout of the evidence root path
	var rootPathLayout = widgets.NewQGridLayout2()
	rootPathLayout.AddWidget2(rootPathLineEdit, 0, 0, 0)
	rootPathLayout.AddWidget2(rootPathButton, 0, 2, 0)
	rootPathGroup.SetLayout(rootPathLayout)

	// layout of output path
	var outputPathLayout = widgets.NewQGridLayout2()
	outputPathLayout.AddWidget2(outputPathLineEdit, 0, 0, 0)
	outputPathLayout.AddWidget2(outputPathButton, 0, 2, 0)
	outputPathGroup.SetLayout(outputPathLayout)

	// layout of preset/filelist path
	var filelistPathLayout = widgets.NewQGridLayout2()
	filelistPathLayout.AddWidget2(presetPathLabel, 0, 0, 0)
	filelistPathLayout.AddWidget2(presetComboBox, 1, 0, 0)
	filelistPathLayout.AddWidget2(filelistPathLabel, 2, 0, 0)
	filelistPathLayout.AddWidget2(filelistPathLineEdit, 3, 0, 0)
	filelistPathLayout.AddWidget2(filelistPathButton, 3, 2, 0)
	filelistPathGroup.SetLayout(filelistPathLayout)

	// layout of options
	var optionsLayout = widgets.NewQGridLayout2()
	optionsLayout.AddWidget2(optionsCheckBoxProfiler, 0, 0, 0)
	optionsLayout.AddWidget2(optionsCheckBoxFileHash, 1, 0, 0)
	optionsLayout.AddWidget2(optionsCheckBoxStat, 2, 0, 0)
	optionsLayout.AddWidget2(optionsCheckBoxDmg, 3, 0, 0)
	optionsGroup.SetLayout(optionsLayout)

	// add each group to widget
	var layout = widgets.NewQGridLayout2()
	layout.AddWidget2(rootPathGroup, 0, 0, 0)
	layout.AddWidget2(outputPathGroup, 1, 0, 0)
	layout.AddWidget2(filelistPathGroup, 2, 0, 0)
	layout.AddWidget2(optionsGroup, 3, 0, 0)
	layout.AddWidget2(runButton, 4, 0, 3)

	// make window
	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(400, 400)
	window.SetWindowTitle("macOS Triage Tool")

	// set layout
	var widget = widgets.NewQWidget(window, 0)
	widget.SetLayout(layout)
	window.SetCentralWidget(widget)

	return window
}

func openFileDialog(path string, lineEdit *widgets.QLineEdit) {

	fileDialog := widgets.NewQFileDialog2(nil, "Path", path, "")
	if fileDialog.Exec() != int(widgets.QDialog__Accepted) {
		return
	}

	selectedFilePath := fileDialog.SelectedFiles()[0]

	lineEdit.SetText(selectedFilePath)
}
