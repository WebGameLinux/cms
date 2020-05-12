package local

import (
		"os"
		"time"
)

type FileInfo struct {
		Size           int64
		Exists         bool
		Mode           os.FileMode
		IsDir          bool
		ModTime        time.Time
		Name           string
		PermissionAble bool
}

func NewFileInfo() *FileInfo {
		return new(FileInfo)
}

func (this *FileInfo) Load(info os.FileInfo, perm ...bool) *FileInfo {
		this.Size = info.Size()
		this.Name = info.Name()
		this.Mode = info.Mode()
		this.ModTime = info.ModTime()
		this.Exists = true
		this.IsDir = info.IsDir()
		if len(perm) > 0 {
				this.PermissionAble = perm[0]
		}
		return this
}

func (this *FileInfo) GetFileInfo(filename string) *FileInfo {
		var (
				info = NewFileInfo()
				perm = false
		)
		state, err := os.Stat(filename)
		if err != nil {
				if err == os.ErrNotExist {
						return info
				}
				if err == os.ErrPermission {
						perm = false
				}
		}
		if state != nil {
				return info.Load(state, perm)
		}
		return info
}

func (this *FileInfo) State(filename string) (os.FileInfo, bool) {
		if filename == "" {
				return nil, false
		}
		state, err := os.Stat(filename)
		if err != nil {
				if err == os.ErrNotExist {
						return nil, false
				}
				if err == os.ErrPermission {
						return state, false
				}
				if err == os.ErrInvalid {
						return nil, false
				}
		}
		return state, true
}

var nilInfo = NewFileInfo()

func GetFileInfo(filename string) *FileInfo {
		return nilInfo.GetFileInfo(filename)
}

func GetFileSize(filename string) int64 {
		state, err := os.Stat(filename)
		if err != nil {
				return 0
		}
		return state.Size()
}

func GetFileSizeFormat(filename string, format ...string) string {
		size := GetFileSize(filename)
		if size == 0 {
				return EmptySize.String()
		}
		if len(format) != 0 {
				return EmptySize.ParseInt(size).Format(format[0])
		}
		return EmptySize.ParseInt(size).String()
}

func FileExists(filename string) bool {
		_, err := os.Stat(filename)
		if err != nil {
				if err == os.ErrNotExist {
						return false
				}
				if err != os.ErrPermission {
						return false
				}
		}
		return true
}

func PermissionAble(filename string) bool {
		_, err := os.Stat(filename)
		if err != nil {
				return false
		}
		return true
}

func IsDir(filename string) bool {
		state, err := os.Stat(filename)
		if err != nil {
				return false
		}
		return state.IsDir()
}

func IsMode(filename string, mode os.FileMode) bool {
		state, err := os.Stat(filename)
		if err != nil {
				return false
		}
		m := state.Mode()
		return m&m == 0 || m == mode
}

func IsSymlink(filename string) bool {
		return IsMode(filename, os.ModeSymlink)
}

func IsPipe(filename string) bool {
		return IsMode(filename, os.ModeNamedPipe)
}

func IsDevice(filename string) bool {
		return IsMode(filename, os.ModeDevice)
}

func IsCharDevice(filename string) bool {
		return IsMode(filename, os.ModeCharDevice)
}

func IsIrregular(filename string) bool {
		return IsMode(filename, os.ModeIrregular)
}
