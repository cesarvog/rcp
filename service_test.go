package main

import (
	"testing"
	"time"
	"os"
)

func TestMustDeleteFile(t *testing.T) {
	fi, err := os.Stat(TmpDir + FilePrefix + "xxx")
	if err != nil {
		t.Errorf("test file not found")
		return
	}

	fit := NewFileInfoTest(fi, time.Now())
	if MustDeleteFile(fit) {
		t.Errorf("Must not delete file created now")	
	}
	fit.modTime = time.Now().Add(time.Hour * -2)

	if !MustDeleteFile(fit) {
		t.Errorf("Must not delete file created now")	
	}
}


//Should delte just files with prefix f_
func TestMustNotDeleteDifPrefixFile(t *testing.T) {
	fi, err := os.Stat(TmpDir + "banana")
	if err != nil {
		t.Errorf("test file not found")
		return
	}

	fit := NewFileInfoTest(fi, time.Now().Add(time.Hour * -2))
	if MustDeleteFile(fit) {
		t.Errorf("Must not delete file created now")	
	}
}


type FileInfoTest struct {
    name string       
    size int64        
    mode os.FileMode     
    modTime time.Time 
    isDir bool        
    sys interface{}   
}

func (fit FileInfoTest) Name() string {
	return fit.name
}
func (fit FileInfoTest) Size() int64 {
	return fit.size 
}
func (fit FileInfoTest) Mode() os.FileMode {
	return fit.mode 
}
func (fit FileInfoTest) ModTime() time.Time {
	return fit.modTime 
}
func (fit FileInfoTest) IsDir() bool {
	return fit.isDir
}
func (fit FileInfoTest) Sys() interface{} {
	return fit.sys 
}

func NewFileInfoTest(fi os.FileInfo, t time.Time) FileInfoTest {
	info := FileInfoTest{}
	info.name = fi.Name()
	info.size = fi.Size()
	info.mode = fi.Mode()
	info.modTime = t
	info.isDir = fi.IsDir()
	info.sys = fi.Sys()

	return info
}


