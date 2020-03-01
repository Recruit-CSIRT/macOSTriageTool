package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	"github.com/Recruit-CSIRT/macOSTriageTool/pkg/log"
)

type FileInfo struct {
	filepath string
	inode 	int
	uid 	int
	gid 	int
	flags	int
	mode 	os.FileMode
	size 	int64
	atime   time.Time
	btime	time.Time
	ctime 	time.Time
	mtime 	time.Time
	hash 	string
}

func NewFileInfo(filepath *string) FileInfo {
	fileInfo := FileInfo{}
	fileInfo.filepath = *filepath
	return fileInfo
}

func (s *FileInfo) SetStat (fileinfo *os.FileInfo) () {
	s.inode = int((*fileinfo).Sys().(*syscall.Stat_t).Ino)
	s.uid 	= int((*fileinfo).Sys().(*syscall.Stat_t).Uid)
	s.gid	= int((*fileinfo).Sys().(*syscall.Stat_t).Gid)
	s.flags = int((*fileinfo).Sys().(*syscall.Stat_t).Flags)
	s.mode 	= (*fileinfo).Mode()
	s.size 	= (*fileinfo).Size()
	s.atime = time.Unix((*fileinfo).Sys().(*syscall.Stat_t).Atimespec.Sec, 9)
	s.btime = time.Unix((*fileinfo).Sys().(*syscall.Stat_t).Birthtimespec.Sec, 9)
	s.ctime = time.Unix((*fileinfo).Sys().(*syscall.Stat_t).Ctimespec.Sec, 9)
	s.mtime = time.Unix((*fileinfo).Sys().(*syscall.Stat_t).Mtimespec.Sec, 9)
}


func (s *FileInfo) SetHash ()  error {
	file, err := os.Open(s.filepath)
	if err != nil {
		log.FileInfoErrorLogger.Error("[-] Failed to open the file. File: " + s.filepath)
		return err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.FileInfoErrorLogger.Error("[-] Failed to calculate md5 hash. File: " + s.filepath)
		return err
	}
	sum := hash.Sum(nil)

	s.hash = fmt.Sprintf("%x", sum)
	return nil
}

func (s *FileInfo) WriteLogger () {

	log.FileInfoLogger.Info(
		fmt.Sprintf(`"%s","%d","%d","%d","%d","%s","%d","%s","%s","%s","%s","%s"`,
			s.filepath,
			s.inode,
			s.uid,
			s.gid,
			s.flags,
			s.mode.String(),
			s.size,
			s.atime,
			s.btime,
			s.ctime,
			s.mtime,
			s.hash,
		))
}

func WriteColumn () {

	log.FileInfoLogger.Info(
		fmt.Sprintf(`"%s","%s","%s","%s","%s","%s","%s","%s","%s","%s","%s","%s"`,
		"filepath",
			"inode",
			"uid",
			"gid",
			"flags",
			"mode",
			"size",
			"atime",
			"btime",
			"ctime",
			"mtime",
			"hash",
		))
}

