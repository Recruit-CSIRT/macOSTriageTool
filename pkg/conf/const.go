package conf

var PresetString = []string{"AllList", "Malware", "Fraud", "macripper", "Only use custom list"}

const (
	CmdHdiutil = "/usr/bin/hdiutil"
	CmdSystemProfiler = "/usr/sbin/system_profiler"
	CmdDitto = "/usr/bin/ditto"

	TmpDir = "/tmp"
	VolumeName = "ROOT"
	OutputDirName = "evidence"
	VolumePath = "/Volumes"
	ScheduledTriageFileListName = "scheduled_triage_file_list.txt"
	SystemProfilerFileName = "system_profiler.txt"

	DmgFileSystem = "APFS"
)