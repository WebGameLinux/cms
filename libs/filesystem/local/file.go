package local

import (
		"encoding/json"
		"errors"
		"fmt"
		log "github.com/sirupsen/logrus"
		"io"
		"os"
		"path/filepath"
		"reflect"
		"sync"
		"time"
)

type FileSystem interface {
		Name() string
		Root() string
		GetConfig() string
		GetError() error
		Proxy(name string, args ...interface{}) interface{}
}

const BufferSize = 1024

type FileSystemLocal struct {
		root      string
		name      string
		Config    string
		methods   sync.Map
		configure map[string]interface{}
		Error     error
}

var ErrorExistsMethod = errors.New("method not exists")
var ErrorProxyNotCaller = errors.New("method not func")

func (this *FileSystemLocal) Name() string {
		return this.name
}

func (this *FileSystemLocal) setName(name string) *FileSystemLocal {
		if name == "" {
				return this
		}
		this.name = name
		return this
}

func (this *FileSystemLocal) Root() string {
		return this.root
}

func (this *FileSystemLocal) setRoot(root string) *FileSystemLocal {
		if root != "" {
				if !filepath.IsAbs(root) {
						tmp, err := filepath.Abs(root)
						if err == nil {
								root = tmp
						}
				}
		}
		this.root = root
		return this
}

func (this *FileSystemLocal) GetConfig() string {
		return this.Config
}

func (this *FileSystemLocal) setConfig(config string) *FileSystemLocal {
		if config == "" {
				return this
		}
		this.Config = config
		this.ReloadConfigure()
		return this
}

func (this *FileSystemLocal) ReloadConfigure() bool {
		if this.configure == nil {
				this.configure = make(map[string]interface{})
		}
		if this.Config == "" {
				return false
		}
		if json.Unmarshal([]byte(this.Config), &this.configure) == nil {
				return true
		}
		return false
}

func (this *FileSystemLocal) Proxy(name string, args ...interface{}) interface{} {
		v, ok := this.methods.Load(name)
		if !ok {
				return ErrorExistsMethod
		}
		refType := reflect.TypeOf(v)
		if refType.Kind() != reflect.Func {
				return ErrorProxyNotCaller
		}
		if len(args) == 0 {
				return v
		}
		refValue := reflect.ValueOf(v)
		var values []reflect.Value
		for range args {
				values = append(values, reflect.ValueOf(args))
		}
		return refValue.Call(values)
}

func (this *FileSystemLocal) GetError() error {
		var err = this.Error
		this.Error = nil
		return err
}

func (this *FileSystemLocal) Register(method string, fn interface{}, must ...bool) *FileSystemLocal {
		if len(must) == 0 {
				must = append(must, false)
		}
		if v, ok := this.methods.Load(method); ok {
				if must[0] {
						if v != fn {
								this.methods.Store(must, fn)
						}
				}
				return this
		}
		this.methods.Store(method, fn)
		return this
}

func (this *FileSystemLocal) Open(filename string, args ...interface{}) *os.File {
		if filename == "" {
				return nil
		}
		if this.root != "" && !this.IsAbs(filename) {
				filename = filepath.Join(this.root, filename)
		}
		var (
				flag int
				mode os.FileMode
				argc = len(args)
		)
		if argc > 0 {
				flag, mode = this.ResolverArgs(args...)
				if flag == 0 || flag&os.O_CREATE == 0 {
						if _, err := os.Stat(filename); err != nil {
								this.Error = err
								return nil
						}
				}
				file, err := os.OpenFile(filename, flag, mode)
				if err != nil {
						this.Error = err
						return nil
				}
				return file
		}
		file, err := os.Open(filename)
		if err != nil {
				this.Error = err
				return nil
		}
		return file
}

func (this *FileSystemLocal) Write(filename string, content ...[]byte) int {
		if len(content) == 0 {
				return 0
		}
		file := this.GetAppendWriter(filename)
		if file == nil {
				return 0
		}
		defer this.Close(file)
		return this.write(file, content...)
}

func (this *FileSystemLocal) GetAppendWriter(filename string) *os.File {
		return this.OpenWrite(filename, os.O_CREATE|os.O_APPEND)
}

func (this *FileSystemLocal) GetCoverWriter(filename string) *os.File {
		return this.OpenWrite(filename, os.O_CREATE)
}

func (this *FileSystemLocal) GetFile(filename string, flag ...int) *os.File {
		var (
				f    int
				args []interface{}
				mode os.FileMode
		)
		for _, v := range flag {
				args = append(args, v)
		}
		f, mode = this.ResolverArgs(args...)
		if f != 0 && !this.checkFlag(f) {
				f = os.O_RDONLY
		}
		if f == 0 {
				f = os.O_RDONLY
		}
		return this.Open(filename, f, mode)
}

func (this *FileSystemLocal) checkFlag(flag int) bool {
		return !(flag&os.O_WRONLY == 0 &&
				flag&os.O_RDWR == 0 && flag&os.O_RDONLY == 0 &&
				flag&os.O_EXCL == 0 && flag&os.O_CREATE == 0 &&
				flag&os.O_SYNC == 0 && flag&os.O_TRUNC == 0)
}

func (this *FileSystemLocal) Close(file *os.File) {
		err := file.Close()
		if err != nil {
				this.Error = err
		}
}

func (this *FileSystemLocal) WriteCover(filename string, content ...[]byte) int {
		if len(content) == 0 {
				return 0
		}
		file := this.GetCoverWriter(filename)
		if file == nil {
				return 0
		}
		defer this.Close(file)
		return this.write(file, content...)
}

func (this *FileSystemLocal) Save(filename string, content []byte, flag ...int) int {
		if len(content) == 0 {
				return 0
		}
		var file *os.File
		if len(flag) != 0 {
				if flag[0]&os.O_RDONLY != 0 {
						return 0
				}
				if flag[0]&os.O_TRUNC == 0 {
						flag[0] = flag[0] | os.O_APPEND
				}
				if flag[0]&os.O_CREATE == 0 {
						flag[0] = flag[0] | os.O_CREATE
				}
				file = this.GetFile(filename, flag[0])
		} else {
				file = this.GetAppendWriter(filename)
		}
		if file == nil {
				return 0
		}
		defer this.Close(file)
		return this.write(file, content)
}

func (this *FileSystemLocal) write(file *os.File, content ...[]byte) int {
		if file == nil {
				return 0
		}
		var size = 0
		for _, buf := range content {
				n, err := file.Write(buf)
				if err != nil {
						this.Error = err
						return size
				}
				size = size + n
		}
		return size
}

func (this *FileSystemLocal) OpenWrite(filename string, args ...interface{}) *os.File {
		var (
				flag int
				mode os.FileMode
		)
		flag, mode = this.ResolverArgs(args...)
		if flag == 0 || (flag&os.O_WRONLY == 0 && flag&os.O_RDWR == 0) {
				flag = flag | os.O_WRONLY
		}
		return this.Open(filename, flag, mode)
}

func (this *FileSystemLocal) ResolverArgs(args ...interface{}) (int, os.FileMode) {
		var (
				flag int
				mode os.FileMode
		)
		for _, v := range args {
				if v == nil {
						continue
				}
				if n, ok := v.(int); ok && n != 0 && flag == 0 {
						flag = n
				}
				if m, ok := v.(os.FileMode); ok && m != 0 && mode == 0 {
						mode = m
				}
				if m, ok := v.(uint32); ok && m != 0 && mode == 0 {
						mode = os.FileMode(m)
				}
		}
		if mode == 0 {
				mode = os.ModePerm
		}
		return flag, mode
}

func (this *FileSystemLocal) IsAbs(filename string) bool {
		return filepath.IsAbs(filename)
}

func (this *FileSystemLocal) Abs(filename string, withRoot ...bool) string {
		if len(withRoot) == 0 {
				withRoot = append(withRoot, true)
		}
		if this.root != "" && withRoot[0] {
				if this.IsAbs(filename) {
						return filename
				}
				filename = filepath.Join(this.root, filename)
		}
		abs, err := filepath.Abs(filename)
		if err != nil {
				this.Error = err
				return ""
		}
		return abs
}

func (this *FileSystemLocal) GetDir(dirname string, flag ...int) interface{} {
		return GetDir(dirname, flag...)
}

func (this *FileSystemLocal) NewDir(dirname string, mode os.FileMode, recursion ...bool) error {
		if mode <= 0 {
				mode = os.ModePerm
		}
		if len(recursion) != 0 && recursion[0] {
				return os.MkdirAll(dirname, mode)
		}
		return os.Mkdir(dirname, mode)
}

func (this *FileSystemLocal) IsDir(dirname string) bool {
		return GetFileInfo(dirname).IsDir
}

func (this *FileSystemLocal) Copy(src, dst string) bool {
		var ok bool
		ok, this.Error = FileCopy(src, dst)
		return ok
}

func (this *FileSystemLocal) ReName(old, name string) bool {
		this.Error = os.Rename(old, name)
		return this.Error == nil
}

func (this *FileSystemLocal) Delete(name string, all ...bool) bool {
		if len(all) != 0 && all[0] {
				this.Error = os.RemoveAll(name)
		} else {
				this.Error = os.Remove(name)
		}
		return this.Error == nil
}

func (this *FileSystemLocal) SoftDelete(name string) bool {
		if !FileExists(name) {
				return true
		}
		dirname := this.CreateDir(filepath.Join(this.GetUserHome(), "trash"))
		if dirname == "" {
				dirname = "./"
		}
		newName := filepath.Join(dirname, fmt.Sprintf(".%d_%s", time.Now().Unix(), name))
		return this.ReName(name, newName)
}

func (this *FileSystemLocal) CreateDir(dirname string) string {
		state, err := os.Stat(dirname)
		if err != nil {
				return ""
		}
		if !state.IsDir() {
				return ""
		}
		err = os.MkdirAll(dirname, os.ModePerm)
		if err != nil {
				return ""
		}
		return dirname
}

func (this *FileSystemLocal) Read(filename string) []byte {
		var content []byte

		return content
}

func (this *FileSystemLocal) ReadDir(dirname string, filter ...func(string) bool) []string {
		return NewDir(dirname).Read(-1, filter...)
}

func (this *FileSystemLocal) GetSystemUserConfigDir() string {
		var dirname string
		dirname, this.Error = os.UserConfigDir()
		return dirname
}

func (this *FileSystemLocal) GetUserHome() string {
		var dirname string
		dirname, this.Error = os.UserHomeDir()
		return dirname
}

func (this *FileSystemLocal) GetUserCacheDir() string {
		var dirname string
		dirname, this.Error = os.UserCacheDir()
		return dirname
}

func (this *FileSystemLocal) Init() *FileSystemLocal {
		this.init()
		return this
}

func (this *FileSystemLocal) init() {
		this.Register("Abs", this.Abs)
		this.Register("AppendWrite", this.Write)
		this.Register("WriteOrCreate", this.Write)
		this.Register("CoverWrite", this.write)
		this.Register("Open", this.Open)
		this.Register("IsDir", this.IsDir)
		this.Register("GetFile", this.GetFile)
		this.Register("GetDir", this.GetDir)
		this.Register("Save", this.Save)
		this.Register("SetRoot", this.setRoot)
		this.Register("MkDir", this.NewDir)
		this.Register("MKdir", this.NewDir)
		this.Register("mkdir", this.NewDir)
		this.Register("Copy", this.Copy)
		this.Register("rename", this.ReName)
		this.Register("mv", this.ReName)
		this.Register("Move", this.ReName)
		this.Register("AutoDir", this.CreateDir)
		this.Register("CreateDir", this.CreateDir)
}

func NewFileSystemLocal(root string) *FileSystemLocal {
		local := new(FileSystemLocal)
		local.root = root
		local.Init()
		return local
}

var localSystem = NewFileSystemLocal("")

func GetFileLocalSystem() *FileSystemLocal {
		return localSystem
}

func GetDisk(name string) FileSystem {
		GetFileSystemManager().Get(name)
		return nil
}

func FileCopy(src, dst string) (bool, error) {
		sourceFileStat, err := os.Stat(src)
		if err != nil {
				return false, err
		}

		if !sourceFileStat.Mode().IsRegular() {
				return false, fmt.Errorf("%s is not a regular file", src)
		}
		source, err := os.Open(src)
		if err != nil {
				return false, err
		}
		defer closeFile(source)
		destination, err := os.Create(dst)
		if err != nil {
				return false, err
		}
		defer closeFile(destination)
		buf := make([]byte, BufferSize)
		for {
				n, err := source.Read(buf)
				if err != nil && err != io.EOF {
						return false, err
				}
				if n == 0 {
						break
				}

				if _, err := destination.Write(buf[:n]); err != nil {
						return false, err
				}
		}
		return true, nil
}

func closeFile(file *os.File) {
		err := file.Close()
		if err == nil {
				return
		}
		if err == os.ErrClosed {
				return
		}
		log.Printf(err.Error())
}
