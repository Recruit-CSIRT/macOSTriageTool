package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/Recruit-CSIRT/macOSTriageTool/pkg/conf"
	"github.com/Recruit-CSIRT/macOSTriageTool/pkg/log"
)

var presetFileList = [][]string{
	{"/private/var/audit/", "true", "true", "true", "false", "User Account, Audit"},
	{"/.fseventsd/", "true", "true", "true", "false", "File System Event"},
	{"/.Spotlight-V100/Store-V2/*/.store.db", "true", "true", "true", "true", "Spotlight"},
	{"/.MobileBackups", "true", "false", "true", "false", "TimeMachine"},
	{"/.Spotlight-V100/Store-V2/*/store.db", "true", "true", "true", "true", "Spotlight"},
	{"/private/var/db/uuidtext/", "true", "true", "true", "true", "Unified Logs"},
	{"/private/var/db/diagnostics/", "true", "true", "true", "true", "Unified Logs"},
	{"/private/var/log/system.log*", "false", "true", "true", "false", "System.log"},
	{"/private/var/log/asl/", "false", "true", "true", "false", "ASL (Apple System Log)"},
	{"/Applications/", "true", "false", "false", "false", "Execution History"},
	{"/Library/Extensions/", "true", "false", "false", "false", "Kext"},
	{"/private/var/log/daily.out", "false", "true", "true", "false", "Disk Usasge History, Maintenance log"},
	{"/private/var/log/monthly.out", "false", "true", "true", "false", "Maintenance log"},
	{"/private/var/log/weekly.out", "false", "true", "true", "false", "Maintenance log"},
	{"/private/var/vm/sleepimage", "true", "false", "false", "false", "RAM"},
	{"/Library/Preferences/.GlobalPreferences.plist ", "false", "false", "false", "false", "Timezone"},
	{"/Library/Preferences/com.apple.accounts.plist", "false", "false", "false", "false", "User Account"},
	{"/Library/Preferences/com.apple.alf.plist", "false", "false", "false", "false", "Firewall"},
	{"/Library/Preferences/com.apple.Bluetooth.plist", "false", "false", "false", "false", "Bluetooth"},
	{"/Library/Preferences/com.apple.loginwindow.plist ", "false", "false", "false", "false", "User Account, User Behavior"},
	{"/Library/Preferences/com.apple.preferences.accounts.plist", "false", "false", "false", "false", "User Account, User Behavior"},
	{"/Library/Preferences/com.apple.SoftwareUpdate.plist", "false", "false", "false", "false", "App update"},
	{"/Library/Preferences/com.apple.TimeMachine.plist", "false", "false", "false", "false", "TimeMachine"},
	{"/Library/Preferences/OpenDirectory/Configurations/Active Directory/", "false", "false", "false", "false", "Active Directory"},
	{"/Library/Preferences/org.cups.printers.plist", "false", "false", "false", "false", "Print"},
	{"/Library/Preferences/SystemConfiguration/com.apple.airport.preferences.plist", "false", "false", "false", "false", "Network Analysis (Wifi), User Behavior"},
	{"/Library/Preferences/SystemConfiguration/com.apple.smb.server.plist", "false", "false", "false", "false", "User Settings"},
	{"/Library/Preferences/SystemConfiguration/NetworkInterfaces.plist", "false", "false", "false", "false", "Network Analysis"},
	{"/Library/Preferences/SystemConfiguration/preferences.plist", "false", "false", "false", "false", "Network Analysis"},
	{"/private/var/vm/swapfile", "true", "false", "false", "false", "Swapfile"},
	{"/.DocumentRevisions-V100/", "true", "true", "true", "false", "DocumentRevisions"},
	{"/.DocumentRevisions-V100/db-V1/db.sqlite", "true", "true", "true", "false", "DocumentRevisions"},
	{"/.Spotlight-V100/Store-V2/*/Cache/", "true", "true", "true", "false", "Spotlight"},
	{"/.Spotlight-V100/VolumeConfiguration.plist", "true", "true", "true", "false", "Spotlight, TimeMachine"},
	{"/private/etc/emond.d/", "true", "false", "false", "false", "Event Monitor Daemon"},
	{"/Library/LaunchAgents/*.plist", "true", "true", "false", "true", "Persistance"},
	{"/private/etc/kcpassword", "true", "false", "false", "false", "User Account"},
	{"/Library/LaunchDaemons/*.plist", "true", "true", "false", "true", "Persistance"},
	{"/Library/Logs/DiagnosticReports/", "true", "true", "true", "false", "Execution History"},
	{"/Library/Preferences/", "true", "true", "true", "false", "Preferences"},
	{"/Library/Receipts/InstallHistory.plist", "true", "true", "true", "false", "App install"},
	{"/Library/ScriptingAdditions", "true", "true", "false", "false", "Persistance"},
	{"/Library/StartupItems/", "true", "true", "false", "false", "Persistance"},
	{"/private/etc/crontab", "true", "true", "false", "false", "Cron"},
	{"/private/etc/defaults/periodic.conf", "true", "true", "false", "false", "Cron"},
	{"/private/etc/hosts", "true", "true", "false", "false", "Network Analysis"},
	{"/private/etc/localtime", "true", "true", "true", "true", "Timezone"},
	{"/private/etc/periodic.conf", "true", "true", "false", "false", "Cron"},
	{"/private/var/db/CoreDuet/Knowledge/", "true", "false", "true", "false", "IME"},
	{"/private/var/db/dhcpclient/DUID_IA.plist", "true", "false", "false", "false", "Network Analysis"},
	{"/private/etc/periodic.conf.local", "true", "true", "false", "false", "Cron"},
	{"/private/etc/periodic/", "true", "true", "false", "false", "Cron"},
	{"/private/var/at/tabs/", "true", "true", "false", "false", "Cron"},
	{"/private/var/db/.AppleInstallType.plist", "true", "true", "true", "false", "Mac system installation"},
	{"/private/var/db/.AppleSetupDone", "true", "true", "true", "false", "Mac system installation"},
	{"/private/var/db/RemoteManagement/caches", "true", "false", "false", "false", "Apple Remote Desktop"},
	{"/private/var/db/RemoteManagement/ClientCaches/", "true", "false", "false", "false", "Apple Remote Desktop"},
	{"/private/var/db/RemoteManagement/RMDB/", "true", "false", "false", "false", "Apple Remote Desktop"},
	{"/private/var/db/.LastGKReject", "true", "true", "true", "false", "Gatekeeper and XProtect"},
	{"/private/var/folders/*/*/*/com.apple.notificationcenter/", "true", "false", "false", "false", "Notifications"},
	{"/private/var/folders/*/*/*/com.apple.ScreenTimeAgent/Store/", "true", "false", "true", "false", "Screentime"},
	{"/private/var/folders/*/*/C/com.apple.QuickLook.thumbnailcache/*", "true", "false", "true", "false", "Quicklook"},
	{"/private/var/db/analyticsd/aggregates/", "true", "true", "true", "false", "Execution History"},
	{"/private/var/db/com.apple.xpc.launchd/disabled.*.plist", "true", "true", "false", "false", "Autorun"},
	{"/private/var/log/", "true", "false", "false", "false", "Logs"},
	{"/private/var/db/dhcpclient/leases/", "true", "true", "true", "false", "Network Analysis"},
	{"/private/var/log/cups/", "false", "false", "true", "false", "Print"},
	{"/private/var/db/dslocal/nodes/Default/groups/admin.plist", "true", "true", "true", "false", "Group Account"},
	{"/private/var/db/dslocal/nodes/Default/users/*.plist", "true", "true", "true", "false", "User Account"},
	{"/private/var/db/launchd.db/com.apple.launchd/overrides.plist", "true", "true", "true", "false", "User Settings"},
	{"/private/var/db/receipts/", "true", "true", "true", "false", "App install"},
	{"/private/var/folders/zz/zyxvpxvq6csfxvn_n00000sm00006d/C/cache_encryptedA.db", "true", "true", "true", "false", "Network Analysis (Wifi)"},
	{"/private/var/folders/zz/zyxvpxvq6csfxvn_n00000sm00006d/C/lockCache_encryptedA.db", "true", "true", "true", "false", "Network Analysis (Wifi)"},
	{"/private/var/log/install.log", "false", "true", "true", "false", "Mac system installation, App install, App update"},
	{"/private/var/log/wifi.log", "false", "true", "true", "false", "Wifi"},
	{"/private/var/log/wifi.log.*.bz2", "false", "true", "true", "false", "Wifi"},
	{"/private/var/networkd/netusage.sqlite*", "true", "true", "true", "false", "Network Analysis"},
	{"/private/var/spool/cups/", "true", "false", "true", "false", "Print"},
	{"/private/var/run/resolv.conf", "true", "true", "false", "false", "Network Analysis"},
	{"/private/var/run/utmpx", "true", "true", "true", "false", "utmpx"},
	{"/System/Library/CoreServices/SystemVersion.plist", "true", "true", "true", "true", "Systeminfo (Mac OS Ver.)"},
	{"/System/Library/Extensions/", "true", "false", "false", "false", "Kext"},
	{"/System/Library/LaunchAgents/*.plist", "true", "true", "false", "false", "Persistance"},
	{"/System/Library/LaunchDaemons/*.plist", "true", "true", "true", "true", "Persistance, Maintenance Log"},
	{"/System/Library/LaunchDaemons/com.apple.periodic-daily.plist", "false", "false", "false", "false", "Maintenance log"},
	{"/System/Library/LaunchDaemons/com.apple.periodic-monthly.plist", "false", "false", "false", "false", "Maintenance log"},
	{"/System/Library/LaunchDaemons/com.apple.periodic-weekly.plist", "false", "false", "false", "false", "Maintenance log"},
	{"/System/Library/ScriptingAdditions", "true", "false", "false", "false", "Persistance"},
	{"/System/Library/StartupItems", "true", "true", "true", "false", "Persistance"},
	{"/System/Volumes/Data/.fseventsd/", "true", "true", "true", "false", "File System Event"},
	{"/System/Volumes/Data/.Spotlight-V100/Store-V2/*/.store.db", "true", "true", "true", "true", "Spotlight"},
	{"/System/Volumes/Data/.Spotlight-V100/Store-V2/*/store.db", "true", "true", "true", "true", "Spotlight"},
	{"/System/Volumes/Data/private/var/db/.AppleSetupDone", "true", "true", "true", "false", "Mac system installation"},
	{"/System/Volumes/Data/private/var/db/Spotlight-V100/BootVolume/Store-V2/*/.store.db", "true", "true", "true", "true", "Spotlight"},
	{"/System/Volumes/Data/private/var/db/Spotlight-V100/BootVolume/Store-V2/*/store.db", "true", "true", "true", "true", "Spotlight"},
	{"/System/Volumes/Data/private/var/db/Spotlight-V100/BootVolume/VolumeConfiguration.plist", "true", "true", "true", "false", "Spotlight, TimeMachine"},
	{"/tmp/", "true", "true", "false", "false", "Tmp Dir"},
	{"/Users/*/.*sh_history", "true", "true", "true", "true", "bash_histry"},
	{"/Users/*/.bash_sessions/", "true", "true", "true", "true", "bash_histry"},
	{"/Users/*/.ssh/", "true", "true", "true", "false", "SSH"},
	{"/Users/*/.Trash", "true", "true", "true", "true", "User Behavior"},
	{"/Users/*/Library/Metadata/CoreSpotlight/index.spotlightV3/.store.db", "true", "true", "true", "false", "Spotlight"},
	{"/Users/*/Applications/", "true", "true", "false", "true", "User Behavior"},
	{"/Users/*/Desktop/", "true", "false", "true", "false", "User Behavior"},
	{"/Users/*/Documents/", "true", "false", "false", "false", "User Behavior"},
	{"/Users/*/Downloads/", "true", "true", "true", "false", "Download Files"},
	{"/Users/*/Library/Accounts/", "true", "true", "true", "false", "User Settings"},
	{"/Users/*/Library/Accounts/VerifiedBackup/", "false", "false", "false", "false", "User Settings"},
	{"/Users/*/Library/Application Support/AddressBook/", "true", "false", "true", "false", "User Behavior"},
	{"/Users/*/Library/Application Support/App Store/", "true", "false", "false", "false", "Applications"},
	{"/Users/*/Library/Application Support/CloudDocs/session/db/", "true", "false", "true", "false", "User Settings"},
	{"/Users/*/Library/Application Support/com.apple.backgroundtaskmanagementagent/backgrounditems.btm", "true", "true", "false", "false", "Persistance"},
	{"/Users/*/Library/Application Support/com.apple.sharedfilelist/", "true", "true", "true", "true", "Most Recently Used"},
	{"/Users/*/Library/Application Support/com.apple.sharedfilelist/com.apple.LSSharedFileList.ApplicationRecentDocuments/", "false", "false", "false", "false", "Most Recently Used"},
	{"/Users/*/Library/Application Support/com.apple.sharedfilelist/com.apple.LSSharedFileList.FavoriteVolumes.sfl2", "false", "false", "false", "false", "Most Recently Used"},
	{"/Users/*/Library/Application Support/com.apple.sharedfilelist/com.apple.LSSharedFileList.Recent*", "false", "false", "false", "false", "Most Recently Used"},
	{"/Users/*/Library/Application Support/com.apple.spotlight.Shortcuts", "true", "true", "true", "true", "Spotlight"},
	{"/Users/*/Library/Application Support/com.apple.spotlight/appList.dat", "true", "true", "true", "false", "Applications"},
	{"/Users/*/Library/Application Support/Firefox/profiles.ini", "true", "true", "true", "false", "Firefox"},
	{"/Users/*/Library/Application Support/Firefox/Profiles/", "true", "true", "true", "false", "Firefox"},
	{"/Users/*/Library/Application Support/Google/Chrome/", "true", "false", "false", "false", "Google Chrome"},
	{"/Users/*/Library/Application Support/Google/Chrome/Default/", "false", "true", "true", "false", "Google Chrome"},
	{"/Users/*/Library/Application Support/Google/Chrome/Profile*/", "false", "true", "true", "false", "Google Chrome"},
	{"/Users/*/Library/Application Support/icdd/", "true", "true", "true", "false", "User Behavior"},
	{"/Users/*/Library/Application Support/icdd/deviceInfoCache.plist", "false", "false", "false", "false", "User Behavior"},
	{"/Users/*/Library/Application Support/iCloud/", "true", "false", "true", "false", "iCloud"},
	{"/Users/*/Library/Application Support/Knowledge/", "true", "false", "true", "false", "IME"},
	{"/Users/*/Library/Application Support/Microsoft/Skype for Desktop", "true", "false", "true", "false", "Skype"},
	{"/Users/*/Library/Application Support/MobileSync/Backup/", "true", "false", "true", "false", "User Settings"},
	{"/Users/*/Library/Application Support/Skype/", "true", "false", "true", "false", "Skype"},
	{"/Users/*/Library/Caches/com.apple.Safari/Cache.db", "true", "true", "true", "false", "Safari"},
	{"/Users/*/Library/Caches/com.apple.Safari/TabSnapshots/", "true", "false", "false", "false", "Safari"},
	{"/Users/*/Library/Caches/com.apple.Safari/TabSnapshots/Metadata.db", "false", "false", "true", "false", "Safari"},
	{"/Users/*/Library/Caches/com.apple.Safari/Webpage Previews/", "true", "true", "true", "false", "Safari"},
	{"/Users/*/Library/Caches/Firefox/Profiles/", "true", "false", "false", "false", "Firefox"},
	{"/Users/*/Library/Caches/Firefox/Profiles/*/cache2/entries/", "false", "true", "true", "false", "Firefox"},
	{"/Users/*/Library/Caches/Google/Chrome/", "true", "false", "true", "false", "Google Chrome"},
	{"/Users/*/Library/Calendars/", "true", "false", "true", "false", "Calendars"},
	{"/Users/*/Library/Calendars/*.caldav/*/Events/*.ics", "false", "false", "false", "false", "Calendar"},
	{"/Users/*/Library/Calendars/*.calendar/Events/*.ics", "false", "false", "false", "false", "Calendar"},
	{"/Users/*/Library/Calendars/Calendar Cache*", "false", "false", "false", "false", "Calendar"},
	{"/Users/*/Library/Containers/com.apple.mail/Data/Library/Mail Downloads/", "true", "false", "false", "false", "Mail"},
	{"/Users/*/Library/Containers/com.apple.mail/Data/Library/Preferences/com.apple.mail.plist", "true", "true", "true", "false", "Mail"},
	{"/Users/*/Library/Containers/com.apple.Maps/", "true", "false", "true", "false", "Map"},
	{"/Users/*/Library/Containers/com.apple.Maps/Data/Library/Maps/GeoBookmarks.plist", "false", "false", "false", "false", "Map"},
	{"/Users/*/Library/Containers/com.apple.Maps/Data/Library/Maps/GeoHistory.mapsdata", "false", "false", "false", "false", "Map"},
	{"/Users/*/Library/Containers/com.apple.Notes/Data/Library/Application Support/Notes", "true", "false", "true", "false", "Notes"},
	{"/Users/*/Library/Containers/com.apple.Notes/Data/Library/Notes/", "true", "false", "true", "false", "Notes"},
	{"/Users/*/Library/Containers/com.microsoft.*/Data/Library/Preferences/com.microsoft.*.securebookmarks.plist", "true", "true", "true", "true", "Most Recently Used(MS Office)"},
	{"/Users/*/Library/Cookies/Cookies.binarycookies", "true", "true", "true", "false", "Safari"},
	{"/Users/*/Library/Dictionaries/JapaneseInputMethod/", "true", "true", "true", "true", "IME"},
	{"/Users/*/Library/Group Containers/group.com.apple.notes/", "true", "false", "true", "false", "Notes"},
	{"/Users/*/Library/Keychains/", "true", "true", "true", "false", "Keychain"},
	{"/Users/*/Library/Keychains/login.keychain-db", "false", "false", "false", "false", "User Account"},
	{"/Users/*/Library/LaunchAgents/", "true", "true", "false", "true", "Persistance"},
	{"/Users/*/Library/Logs/fsck_hfs.log", "true", "true", "true", "false", "dmg open"},
	{"/Users/*/Library/Mail Downloads/", "true", "false", "false", "false", "Mail"},
	{"/Users/*/Library/Mail/", "true", "true", "true", "false", "Mail"},
	{"/Users/*/Library/Messages/Archive/", "true", "false", "true", "false", "Chat"},
	{"/Users/*/Library/Messages/chat.db*", "true", "false", "true", "false", "Chat"},
	{"/Users/*/Library/Metadata/CoreSpotlight/index.spotlightV3/store.db", "true", "true", "true", "false", "Spotlight"},
	{"/Users/*/Library/Preferences/.GlobalPreferences.plist", "true", "true", "true", "false", "Timezone"},
	{"/Users/*/Library/Preferences/*LSSharedFileList.plist", "true", "true", "true", "true", "Most Recently Used"},
	{"/Users/*/Library/Preferences/ByHost/com.apple.loginwindow.*.plist", "true", "true", "true", "false", "User Account"},
	{"/Users/*/Library/Preferences/com.apple.AddressBook.plist", "true", "true", "true", "false", "User Behavior"},
	{"/Users/*/Library/Preferences/com.apple.dock.plist", "true", "true", "true", "false", "Dock"},
	{"/Users/*/Library/Preferences/com.apple.finder.plist", "true", "true", "true", "true", "Most Recently Used"},
	{"/Users/*/Library/Preferences/com.apple.iChat.*", "true", "false", "true", "false", "Chat"},
	{"/Users/*/Library/Preferences/com.apple.iChat.AIM.plist", "false", "false", "false", "false", "Chat"},
	{"/Users/*/Library/Preferences/com.apple.iChat.Jabber.plist", "false", "false", "false", "false", "Chat"},
	{"/Users/*/Library/Preferences/com.apple.iChat.plist", "false", "false", "false", "false", "Chat"},
	{"/Users/*/Library/Preferences/com.apple.iChat.Yahoo.plist", "false", "false", "false", "false", "Chat"},
	{"/Users/*/Library/Preferences/com.apple.imservice.*", "true", "false", "true", "false", "Chat"},
	{"/Users/*/Library/Preferences/com.apple.imservice.ids.FaceTime.plist", "false", "false", "false", "false", "Chat"},
	{"/Users/*/Library/Preferences/com.apple.imservice.ids.iMessage.plist", "false", "false", "false", "false", "Chat"},
	{"/Users/*/Library/Preferences/com.apple.iPod.plist", "true", "false", "true", "false", "User Behavior"},
	{"/Users/*/Library/Preferences/com.apple.LaunchServices.QuarantineEvents*", "false", "false", "false", "false", "Gatekeeper and XProtect"},
	{"/Users/*/Library/Preferences/com.apple.LaunchServices.QuarantineEventsV2", "true", "true", "true", "true", "Gatekeeper and XProtect"},
	{"/Users/*/Library/Preferences/com.apple.loginitems.plist", "true", "true", "true", "false", "Persistance"},
	{"/Users/*/Library/Preferences/com.apple.loginwindow.plist", "true", "true", "true", "false", "User Account"},
	{"/Users/*/Library/Preferences/com.apple.recentitems.plist", "true", "true", "true", "true", "Most Recently Used"},
	{"/Users/*/Library/Preferences/com.apple.sidebarlists.plist", "true", "true", "true", "true", "Spotlight"},
	{"/Users/*/Library/Preferences/com.apple.spotlight.plist", "true", "true", "true", "false", "Spotlight"},
	{"/Users/*/Library/Preferences/com.microsoft.*.plist", "true", "true", "true", "true", "Most Recently Used(MS Office)"},
	{"/Users/*/Library/Preferences/MobileMeAccounts.plist", "true", "true", "true", "false", "User Settings"},
	{"/Users/*/Library/Safari/Bookmarks.plist", "true", "false", "true", "false", "Safari"},
	{"/Users/*/Library/Safari/Downloads.plist", "true", "true", "true", "false", "Safari"},
	{"/Users/*/Library/Safari/History.db", "true", "true", "true", "false", "Safari"},
	{"/Users/*/Library/Safari/LastSession.plist", "true", "true", "true", "false", "Safari"},
	{"/Users/*/Library/Safari/TopSites.plist", "true", "true", "true", "false", "Safari"},
	{"/Users/*/Library/Saved Application State/", "true", "false", "false", "false", "Saved State"},
	{"/Users/*/Library/SyncedPreferences/", "true", "false", "true", "false", "User Settings"},
	{"/Users/*/Parallels/", "true", "false", "false", "false", "VM"},
	{"/Users/*/Pictures/", "true", "false", "true", "false", "Pictures"},
	{"/Users/*/Virtual Machines.localized/", "true", "false", "false", "false", "VM"},
	{"/Users/*/VirtualBox VMs/", "true", "false", "false", "false", "VM"},
	{"/usr/local/Cellar/", "true", "false", "true", "false", "App install"},
	{"/usr/local/var/homebrew/linked/", "true", "false", "true", "false", "App install"},
	{"/var/root/.*sh_history", "true", "true", "true", "true", "bash_histry"},
	{"/var/root/.Trash", "true", "true", "true", "false", "User Behavior"},
	{"/var/root/Library/Preferences/com.apple.LaunchServices.QuarantineEventsV2", "true", "true", "true", "true", "Gatekeeper and XProtect"},
	{"/var/tmp", "true", "true", "false", "false", "Tmp Dir"},
}

type TriageFileList struct {
	FileList []string
	UserFileList []string
	NormalizedFileList []string
	presetNum int
}

func (f *TriageFileList) New (config *conf.Config) error {

	f.presetNum = config.SelectedPresetNum

	if config.SelectedPresetNum >= 0 && config.SelectedPresetNum < 4 {
		f.loadPresetList()
	}

	if len(config.UserDefinedFileListName) > 1 {
		err := f.LoadCustomList(config)
		if err != nil {
			log.ToolLogger.Warn(err.Error())
		}
	}

	f.WriteSettingsToLog(config)

	f.NormalizeFileList(config)

	return nil
}

func (f *TriageFileList) loadPresetList () {
	for _, file := range presetFileList {

		if file[f.presetNum+1] == "true" {
			f.FileList = append(f.FileList, file[0])
		}
	}
}

func (f *TriageFileList) LoadCustomList (config *conf.Config) error{

	fp, err := os.Open(config.UserDefinedFileListName)
	if err != nil {
		log.ToolLogger.Error("[-] Failed to open a file list")
		return err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		f.UserFileList = append(f.UserFileList, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.ToolLogger.Error(err.Error())
		return err
	}
	return nil
}

func (f *TriageFileList) NormalizeFileList(config *conf.Config)  {

	for _, line := range append(f.FileList, f.UserFileList...) {
		targetFilePath := filepath.Join(config.RootPath, line)

		if strings.Contains(targetFilePath, "*") {

			matches, err := filepath.Glob(targetFilePath)
			if err != nil {
				log.ToolLogger.Error("[-] Failed to extract * in the path. Path: " + targetFilePath)
				log.ToolLogger.Error(err.Error())
			}

			f.NormalizedFileList = append(f.NormalizedFileList, matches...)

		} else {
			f.NormalizedFileList = append(f.NormalizedFileList, targetFilePath)
		}

	}

}


func (f *TriageFileList) WriteSettingsToLog (config *conf.Config){

	logMessage := ""

	if f.presetNum == 4 {
		logMessage = conf.PresetString[f.presetNum]

	} else if len((f).UserFileList) > 0 {
		logMessage = conf.PresetString[f.presetNum]
		logMessage += " and user defined file " + config.UserDefinedFileListName
	} else {
		logMessage = conf.PresetString[f.presetNum]
	}

	log.ToolLogger.Info("[+] Profile: " + logMessage)

}
