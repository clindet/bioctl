package io

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	clog "github.com/openbiox/bioctl/log"
)

var log = clog.Logger
var logBash = clog.LoggerBash

// Open connect the file using os.OpenFile (auto create if not exist the file)
func Open(name string) (*os.File, error) {
	var fn *os.File
	var err error
	if hasFile, _ := PathExists(name); !hasFile {
		CreateFile(name)
	}
	fn, err = os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0664)
	return fn, err
}

// NewOutStream set stdout or outfn
func NewOutStream(outfn string, url string) *os.File {
	var of *os.File
	if outfn == "" {
		of = os.Stdout
	} else {
		var err error
		of, err = os.OpenFile(outfn, os.O_CREATE|os.O_WRONLY, 0664)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		wd, _ := os.Getwd()
		if url != "" {
			log.Infof("Trying save %s => %s", url, path.Join(wd, outfn))
		}
	}
	return of
}

// CreateFile create file
func CreateFile(name string) error {
	fo, err := os.Create(name)
	if err != nil {
		return err
	}
	defer func() {
		fo.Close()
	}()
	return nil
}

// CreateFileParDir create parant dir of file
func CreateFileParDir(name string) error {
	var err error
	if isExists, _ := PathExists(path.Dir(name)); !isExists {
		err = os.MkdirAll(path.Dir(name), os.ModePerm)
	}
	return err
}

// CreateDir create dir with os.MkdirAll
func CreateDir(name string) error {
	var err error
	if isExists, _ := PathExists(name); !isExists {
		err = os.MkdirAll(name, os.ModePerm)
	}
	return err
}

// ReadLines import a file as []string
func ReadLines(fn string) []string {
	var final []string
	if hasFile, _ := PathExists(fn); !hasFile {
		log.Fatalf("%s not existed.", fn)
	}
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	buf := bufio.NewReader(f)
	for {
		b, errR := buf.ReadBytes('\n')
		if errR != nil {
			if errR == io.EOF {
				break
			}
			log.Fatal(errR.Error())
		}
		final = append(final, strings.ReplaceAll(string(b), "\n", ""))
	}
	return final
}

// CopyFile copy file from srcFilePath to dstFilePath
func CopyFile(dstFilePath string, srcFilePath string) (written int64, err error) {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()
	reader := bufio.NewReader(srcFile)

	dstFile, err := os.OpenFile(dstFilePath, os.O_WRONLY|os.O_CREATE, 0655)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(dstFile)
	defer dstFile.Close()
	return io.Copy(writer, reader)
}

// FileSize get file size
func FileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

// LinkFileToDir create os.Symlink
func LinkFileToDir(oldname, newname string) error {
	hasTargetFile, _ := PathExists(newname)
	if !hasTargetFile {
		if err := CreateFileParDir(newname); err != nil {
			return err
		}
		err := os.Symlink(oldname, newname)
		if err != nil {
			return err
		}
	}
	return nil
}

// PathExists check wheather file is existed
func PathExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//CheckFilesOverwrite check wheather to overwrite existing files
func CheckFilesOverwrite(files []string, overwrite bool) int {
	for _, v := range files {
		hasFile, _ := PathExists(v)
		if hasFile && overwrite {
			if err := os.Remove(v); err != nil {
				log.Fatalf("Can not overwrite the output file %s.", v)
				return 1
			}
		} else if hasFile {
			log.Infof("%s existed.", v)
			return 0
		}
	}
	return -1
}

//CheckRenameOverwrite check wheather to rename existing files
func CheckRenameOverwrite(old string, new string, overwrite bool) int {
	hasNew, _ := PathExists(new)
	hasOld, _ := PathExists(old)
	if hasNew && !overwrite {
		return 0
	} else if hasOld && !overwrite {
		log.Infoln(fmt.Sprintf("%s existed.", old))
		err := os.Rename(old, new)
		logBash.Infof("mv %s %s", old, new)
		if err != nil {
			log.Fatalln(err)
			return 1
		}
	}
	return -1
}

// CheckFilesFinal batch to check wheathre final files is existed
func CheckFilesFinal(files []string) bool {
	for _, v := range files {
		if hasFile, _ := PathExists(v); hasFile {
			return true
		}
	}
	return false
}

// CheckRenameFinal check final file wheather renamed
func CheckRenameFinal(old string, new string) bool {
	hasNew, _ := PathExists(new)
	hasOld, _ := PathExists(old)
	if !hasNew && hasOld {
		log.Infoln(fmt.Sprintf("%s existed.", old))
		err := os.Rename(old, new)
		logBash.Infof("mv %s %s", old, new)
		if err != nil {
			log.Fatalln(err)
			return false
		}
	}
	return true
}

// CheckInputFiles check input files
func CheckInputFiles(files []string) error {
	for _, v := range files {
		if v != "" {
			if strings.Contains(v, " ") {
				v = strings.Split(v, " ")[0]
			}
			var err error
			hasFile, err := PathExists(v)
			if err != nil {
				return err
			}
			pt, _ := exec.LookPath(v)
			if !hasFile && pt == "" {
				err = fmt.Errorf("file %s not existed", v)
				return err
			}
		}
	}
	return nil
}
