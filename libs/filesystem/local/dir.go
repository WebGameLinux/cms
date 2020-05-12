package local

import (
		"os"
		"path/filepath"
		"sort"
)

type DirectionFile struct {
		dirname  string
		file     *os.File
		isOpen   bool
		err      error
		flag     int
		readable bool
		mode     os.FileMode
		Info     *FileInfo
}

func NewDir(dirname string) *DirectionFile {
		var dir = new(DirectionFile)
		name, _ := filepath.Abs(dirname)
		dir.dirname = name
		return dir
}

func (this *DirectionFile) Open() bool {
		if this.isOpen && nil != this.file {
				return true
		}
		// 自动创建检查,权限,类型检查
		if !this.Check() {
				this.isOpen = false
				return false
		}
		var (
				err  error
				file *os.File
		)
		flag := this.getFlag()
		if flag == 0 {
				file, err = os.Open(this.dirname)
		} else {
				file, err = os.OpenFile(this.dirname, flag, this.getMode())
		}
		if err != nil {
				this.err = err
				this.isOpen = false
				return false
		}
		this.file = file
		this.isOpen = true
		return true
}

func (this *DirectionFile) Check() bool {
		var info = this.GetInfo()
		if !info.Exists {
				if !this.Auto() {
						return false
				}
				if err := os.MkdirAll(this.dirname, this.getMode()); err != nil {
						this.err = err
						return false
				}
				this.Info = nil
				info = this.GetInfo()

		}
		if !info.IsDir || !info.PermissionAble {
				return false
		}
		return info.Exists
}

func (this *DirectionFile) SetFlag(flag int) *DirectionFile {
		this.flag = flag
		return this
}

func (this *DirectionFile) Exists() bool {
		return this.GetInfo().Exists
}

func (this *DirectionFile) GetInfo() *FileInfo {
		if this.Info != nil {
				return this.Info
		}
		this.Info = GetFileInfo(this.dirname)
		return this.Info
}

func (this *DirectionFile) Close() {
		this.isOpen = false
		this.err = this.file.Close()
}

func (this *DirectionFile) SetMode(mode os.FileMode) *DirectionFile {
		this.mode = mode
		return this
}

func (this *DirectionFile) GetMode() os.FileMode {
		return this.mode
}

func (this *DirectionFile) getMode() os.FileMode {
		if this.mode == 0 {
				return os.ModePerm
		}
		return this.mode
}

func (this *DirectionFile) GetFlag() int {
		return this.flag
}

func (this *DirectionFile) getFlag() int {
		if this.flag == 0 {
				return os.O_CREATE | os.O_RDWR
		}
		return this.flag
}

func (this *DirectionFile) Auto() bool {
		if this.getFlag()&os.O_CREATE != 0 {
				return true
		}
		return false
}

func (this *DirectionFile) Error() error {
		var e = this.err
		this.err = nil
		return e
}

func (this *DirectionFile) Read(level int, filter ...func(string) bool) []string {
		if level == 0 {
				return []string{}
		}
		if level < 0 {
				level = -1
		}
		files := this.read(level)
		if len(files) == 0 || len(filter) == 0 {
				return files
		}
		var arr []string
		for _, name := range files {
				ok := true
				for _, fn := range filter {
						if !fn(name) {
								ok = false
								break
						}
				}
				if ok {
						arr = append(arr, name)
				}
		}
		return arr
}

func (this *DirectionFile) ReadDir(level int, filter ...func(os.FileInfo) bool) []os.FileInfo {
		if level == 0 {
				return nil
		}
		if level < 0 {
				level = -1
		}
		files := this.readDir(level)
		if len(files) == 0 || len(filter) == 0 {
				return files
		}
		var arr []os.FileInfo
		for _, info := range files {
				ok := true
				for _, fn := range filter {
						if !fn(info) {
								ok = false
								break
						}
				}
				if ok {
						arr = append(arr, info)
				}
		}
		return arr
}

func (this *DirectionFile) readDir(n int) []os.FileInfo {
		var (
				err   error
				files []os.FileInfo
		)
		if !this.Open() {
				return files
		}
		files, err = this.file.Readdir(n)
		if err != nil {
				this.err = err
		}
		sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })
		return files
}

func (this *DirectionFile) read(n int) []string {
		var (
				err   error
				files []string
		)
		if !this.Open() {
				return files
		}
		files, err = this.file.Readdirnames(n)
		if err != nil {
				this.err = err
		}
		sort.Slice(files, func(i, j int) bool { return files[i] < files[j] })
		return files
}

func (this *DirectionFile) Prev() *DirectionFile {
		dirname := filepath.Dir(this.dirname)
		if dirname == this.dirname {
				return this
		}
		return NewDir(dirname)
}

var dir = NewDir(os.Args[0])

func GetCurrentDir() *DirectionFile {
		return dir
}

func GetDir(dirname string, flag ...int) *DirectionFile {
		var dirFs = NewDir(dirname)
		if len(flag) > 0 {
				dirFs.SetFlag(flag[0])
		}
		return dirFs
}

func Dirname(file string) string {
		return filepath.Dir(file)
}
