package configure

import (
		"github.com/WebGameLinux/cms/utils/mapper"
		"github.com/tietang/props/consul"
		"github.com/tietang/props/kvs"
		"github.com/tietang/props/nacos"
		"github.com/tietang/props/zk"
		"io/ioutil"
		"net"
		"os"
		"path/filepath"
		"strings"
		"sync"
		"time"
		"unicode"
)

// 配置
type Configuration struct {
		uri      string
		mode     string
		params   map[string]interface{}
		composer *kvs.CompositeConfigSource
		sync.RWMutex
}

const (
		NetDriver       = "net"
		FileDriver      = "file"
		ZookeeperSchema = "zookeeper://"
		ConsulSchema    = "consul://"
		NacosSchema     = "nacos://"
		FileSchema      = "file://"
		DefConfPath     = "conf/"
		DefEnvMode      = "development"
		DefRootPath     = "/"
		SchemaFlag      = "://"
		ExtIni          = ".ini"
		ExtConf         = ".conf"
		ExtYaml         = ".yaml"
		ExtProps        = ".properties"
		ScanModeFile    = "file"
		ScanModeDir     = "dir"
)

var (
		instance          *Configuration
		localhost         string
		supportExtensions = []string{ExtIni, ExtYaml, ExtProps, ExtConf}
		schemas           = []string{ZookeeperSchema, ConsulSchema, NacosSchema}
)

// 配置构造
func NewConfiguration(uri string, mode string) *Configuration {
		config := new(Configuration)
		config.uri = uri
		config.mode = mode
		config.params = make(map[string]interface{})
		config.init()
		return config
}

// 获取配置对象
func GetConfigure(args ...string) *Configuration {
		if len(args) == 0 {
				args = append(args, DefConfPath, DefEnvMode)
		}
		if instance == nil {
				if len(args) > 1 {
						instance = NewConfiguration(args[0], args[1])
				} else {
						instance = NewConfiguration(args[0], "")
				}
		}
		return instance
}

func (this *Configuration) init() {
		this.Lock()
		defer this.Unlock()
		if this.uri == "" {
				panic("configure url param miss")
		}
		var driver = this.driver()

		switch driver {
		case NetDriver:
				this.params["driver"] = NetDriver
				this.initParams()
				this.loadConfigureByNetWork()
		case FileDriver:
				fallthrough
		default:
				this.params["driver"] = FileDriver
				this.initParams()
				this.loadConfigureByDir()
		}
}

func (this *Configuration) loadConfigureByDir() {
		var (
				params          = this.parse(FileDriver, FileSchema)
				name            = params.Get("name", DefConfPath+this.mode+"/")
				appendSystemEnv = params.GetBool("appendSystemEnv", true)
				isDir           = params.GetBool("_dir")
				isExists        = params.GetBool("_exists")
				onlyOneFile     = params.GetBool("onlyOneFile")
				extArr          = params.Strings("extension", supportExtensions)
		)
		this.composer = kvs.NewCompositeConfigSource(name, appendSystemEnv)
		if isExists {
				panic("config file not exists")
		}
		if isDir {
				if onlyOneFile {
						file := filepath.Join(params.Get("_path"), this.mode)
						this.resolverFile(file, extArr)
						return
				}
				files := this.ScanDir(name, extArr)
				if len(files) == 0 {
						return
				}
				for _, fs := range files {
						n := this.resolverName(fs)
						node := kvs.NewPropertiesConfigSourceByFile(n, fs)
						if node != nil {
								this.composer.Add(node)
						}
				}
		}
}

func (this *Configuration) resolverFile(file string, extArr []string) {
		var (
				err   error
				state os.FileInfo
		)
		for _, ext := range extArr {
				fs := filepath.Join(file, ext)
				if state, err = os.Stat(fs); err != nil {
						continue
				}
				if state.IsDir() {
						continue
				}
				name := this.resolverName(fs)
				node := kvs.NewPropertiesConfigSourceByFile(name, fs)
				this.composer.Add(node)
				break
		}
}

func (this *Configuration) ScanDir(dir string, extArr []string, modes ...string) []string {
		var (
				mode  string
				files []string
				tmp   string
				err   error
				state os.FileInfo
				split = string(filepath.Separator)
		)
		if len(modes) == 0 {
				modes = append(modes, "")
		}
		if mode == "" {
				mode = modes[0]
		}
		modeType := this.GetMode()
		// 无环境模型
		if modeType == "" {
				mode = ""
		} else {
				if !strings.Contains(dir, modeType) && (mode == "" || mode == ScanModeDir) {
						tmp = filepath.Join(dir, modeType)
						if state, err = os.Stat(tmp); err == nil {
								if state.IsDir() {
										dir = tmp
										mode = ScanModeDir
								} else {
										mode = ScanModeFile
								}
						}
				}
		}
		dir, _ = filepath.Abs(dir)
		rd, err := ioutil.ReadDir(dir)
		if err != nil {
				return files
		}
		for _, fi := range rd {
				if fi.IsDir() {
						arr := this.ScanDir(dir+fi.Name()+split, extArr, mode)
						if len(arr) != 0 {
								files = append(files, arr...)
						}
						continue
				}
				exclude := true
				name := fi.Name()
				if mode == ScanModeFile && !strings.Contains(name, this.mode) {
						continue
				}
				length := len(name)
				for _, ext := range extArr {
						size := len(ext)
						if length <= size {
								continue
						}
						n := name[length-size:]
						if n == ext {
								exclude = false
								break
						}
				}
				if exclude {
						continue
				}
				if absName, err := filepath.Abs(filepath.Join(dir, name)); err == nil {
						files = append(files, absName)
				}
		}
		return files
}

func (this *Configuration) loadConfigureByNetWork() {
		schema := this.getNetWorkDriverSchema()
		switch schema {
		case NacosSchema:
				this.createNacosKvs()
		case ZookeeperSchema:
				this.createZookeeperKvs()
		case ConsulSchema:
				this.createConsulKvs()
		}
}

func (this *Configuration) getNetWorkDriverSchema() string {
		for _, sch := range schemas {
				if strings.Index(this.uri, sch) == 0 {
						return sch
				}
		}
		return ConsulSchema
}

func (this *Configuration) parse(driver string, schema string) mapper.Mapper {
		var v map[string]interface{}
		switch driver {
		case FileDriver:
				return this.parseFileParams(schema)
		case NetDriver:
				return this.parseNetWorkParams(schema)
		}
		v = make(mapper.Mapper)
		return v
}

func (this *Configuration) parseNetWorkParams(schema string) mapper.Mapper {
		var m = mapper.Mapper(this.params)
		switch schema {
		case NacosSchema:
				// address
				fallthrough
		case ConsulSchema:
				// address
				address := m.Get("_url")
				if address != "" && !m.Exists("address") {
						m["address"] = address
				}
		case ZookeeperSchema:
				connStr := m.Strings("_url")
				if len(connStr) != 0 && !m.Exists("connStr") {
						m["connStr"] = connStr
				}
		}
		return m
}

func (this *Configuration) parseFileParams(schema string) mapper.Mapper {
		var (
				ok bool
				v  interface{}
				m  = mapper.Mapper(this.params)
		)
		if schema != FileSchema {
				return m
		}
		if m.Get("driver") == "" {
				return m
		}
		if v, ok = this.params["name"]; ok {
				m["name"] = v
		} else {
				tmp := m.Get("_path", DefConfPath)
				m["name"] = this.resolverName(tmp)
		}
		m["_dir"] = false
		m["_perm"] = false
		m["_exists"] = false
		state, err := os.Stat(m.Get("name"))
		if err == nil {
				m["_dir"] = state.IsDir()
		} else {
				m["_exists"] = os.IsExist(err)
				m["_perm"] = os.IsPermission(err)
		}
		if this.mode == "" && m.Exists("mode") {
				this.mode = m.Get("mode")
		}
		return m
}

func (this *Configuration) resolverName(str string) string {
		if strings.Contains(str, string(filepath.Separator)) {
				strArr := strings.SplitN(str, string(filepath.Separator), -1)
				str = strArr[len(strArr)-1]
				if !this.filterSpace(str) {
						for i := len(strArr) - 1; i >= 0; i-- {
								if !this.filterSpace(strArr[i]) {
										continue
								}
								str = strArr[i]
						}
				}
				str = strings.Trim(str, " ")
		}
		if strings.Contains(str, ".") {
				strArr := strings.SplitN(str, ".", 2)
				if strArr[0] == "" {
						return strArr[1]
				}
				return strArr[0]
		}
		return str
}

func (this *Configuration) createZookeeperKvs() {
		var (
				params   = this.parse(NetDriver, ZookeeperSchema)
				connStr  = params.Strings("connStr", []string{localhost})
				timeout  = params.GetDuration("timeout", 5*time.Second)
				contexts = params.Strings("contexts", this.resolverArray(this.GetMode(DefRootPath)))
		)
		this.composer = zk.NewZookeeperCompositeConfigSource(contexts, connStr, timeout)
}

func (this *Configuration) GetMode(def ...string) string {
		if len(def) == 0 {
				def = append(def, "")
		}
		if this.mode == "" {
				m := mapper.Mapper(this.params)
				tmp := m.Get("mode", def[0])
				if tmp != "" {
						return tmp
				}
		}
		return this.mode
}

func (this *Configuration) createConsulKvs() {
		var (
				params   = this.parse(NetDriver, ConsulSchema)
				address  = params.Get("address", localhost)
				contexts = params.Strings("contexts", this.resolverArray(this.GetMode(DefRootPath)))
		)
		this.composer = consul.NewConsulKeyValueCompositeConfigSource(contexts, address)
}

func (this *Configuration) resolverArray(str string) []string {
		var arr []string
		if str != "" {
				return strings.SplitN(str, ",", -1)
		}
		return arr
}

func (this *Configuration) createNacosKvs() {
		var (
				params  = this.parse(NetDriver, NacosSchema)
				group   = params.Get("group", DefEnvMode)
				tenant  = params.Get("tenant", DefEnvMode)
				address = params.Get("address", localhost)
				dataIds = params.Strings("dataIds", []string{DefRootPath})
		)
		this.composer = nacos.NewNacosPropsCompositeConfigSource(address, group, tenant, dataIds)
}

func (this *Configuration) driver() string {
		if this.net() {
				return NetDriver
		}
		if this.file() {
				return FileDriver
		}
		return FileDriver
}

func (this *Configuration) net() bool {
		for _, sch := range schemas {
				if strings.Contains(this.uri, sch) {
						return true
				}
		}
		if addr := net.ParseIP(this.uri); addr == nil {
				return false
		}
		return true
}

func (this *Configuration) file() bool {
		if strings.Contains(this.uri, FileSchema) {
				return true
		}
		_, err := os.Stat(this.uri)
		if err != nil {
				return false
		}
		return true
}

func (this *Configuration) Kvs() kvs.ConfigSource {
		this.Lock()
		defer this.Unlock()
		return this.composer.Properties
}

func (this *Configuration) Lists() []kvs.ConfigSource {
		this.Lock()
		defer this.Unlock()
		return this.composer.ConfigSources
}

func (this *Configuration) Scope(name string) kvs.ConfigSource {
		this.Lock()
		defer this.Unlock()
		if this.composer.Name() == name {
				return this.composer
		}
		for _, value := range this.composer.ConfigSources {
				if value.Name() == name {
						return value
				}
		}
		return nil
}

func (this *Configuration) Name() string {
		this.Lock()
		defer this.Unlock()
		return this.composer.Name()
}

func (this *Configuration) Size() int {
		this.Lock()
		defer this.Unlock()
		return this.composer.Size()
}

// 初始化协议参数
func (this *Configuration) initParams() {
		var (
				arr    []string
				params string
				tmp    = this.uri
		)
		if strings.Contains(tmp, SchemaFlag) {
				arr = strings.SplitN(tmp, SchemaFlag, 2)
		}
		if len(arr) >= 2 {
				this.params["driver"] = strings.ToLower(strings.Trim(arr[0], " "))
				tmp = strings.Replace(tmp, arr[0]+SchemaFlag, "", 1)
		}
		if strings.Contains(tmp, "?") {
				arr = strings.SplitN(tmp, "?", 2)
				if len(arr) >= 2 {
						this.params["body"] = strings.Trim(arr[0], " ")
						params = arr[1]
				}
				arr = arr[:0]
		}
		if params == "" {
				this.params["body"] = tmp
		} else {
				if strings.Contains(params, "&") {
						arr = strings.SplitN(params, "&", -1)
				}
				if len(arr) == 0 && params != "" {
						arr = append(arr, params)
				}
				if len(arr) > 0 {
						for _, v := range arr {
								v = strings.Trim(v, " ")
								if !this.filterSpace(v) {
										continue
								}
								if strings.Contains(v, "=") {
										kv := strings.SplitN(v, "=", 2)
										if len(kv) != 2 {
												continue
										}
										this.params[strings.Trim(kv[0], " ")] = strings.Trim(kv[1], " ")
								}
						}
				}
		}
		if this.params["driver"] == FileDriver {
				this.params["_path"] = this.params["body"]
		}
		if this.params["driver"] == NetDriver {
				this.params["_url"] = this.params["body"]
		}
}

// 重新载入
func (this *Configuration) ReLoad() {
		this.Lock()
		this.composer = nil
		this.Unlock()
		this.init()
}

// 过滤空白字符串 空白为false,非空白true
func (this *Configuration) filterSpace(str string) bool {
		if str == "" {
				return false
		}
		for _, ch := range []rune(str) {
				if !unicode.IsSpace(ch) {
						return true
				}
		}
		return false
}

func (this *Configuration) KeyValue(key string) *kvs.KeyValue {
		this.Lock()
		defer this.Unlock()
		return this.composer.KeyValue(key)
}

func (this *Configuration) Strings(key string) []string {
		this.Lock()
		defer this.Unlock()
		return this.composer.Strings(key)
}

func (this *Configuration) Ints(key string) []int {
		this.Lock()
		defer this.Unlock()
		return this.composer.Ints(key)
}

func (this *Configuration) Float64s(key string) []float64 {
		this.Lock()
		defer this.Unlock()
		return this.composer.Float64s(key)
}

func (this *Configuration) Durations(key string) []time.Duration {
		this.Lock()
		defer this.Unlock()
		return this.composer.Durations(key)
}

func (this *Configuration) Get(key string) (string, error) {
		this.Lock()
		defer this.Unlock()
		return this.composer.Get(key)
}

func (this *Configuration) GetDefault(key, defaultValue string) string {
		this.Lock()
		defer this.Unlock()
		return this.composer.GetDefault(key, defaultValue)
}

func (this *Configuration) GetInt(key string) (int, error) {
		this.Lock()
		defer this.Unlock()
		return this.composer.GetInt(key)
}

func (this *Configuration) GetIntDefault(key string, defaultValue int) int {
		this.Lock()
		defer this.Unlock()
		return this.composer.GetIntDefault(key, defaultValue)
}

func (this *Configuration) GetDuration(key string) (time.Duration, error) {
		this.Lock()
		defer this.Unlock()
		return this.composer.GetDuration(key)
}

func (this *Configuration) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
		this.Lock()
		defer this.Unlock()
		return this.composer.GetDurationDefault(key, defaultValue)
}

func (this *Configuration) GetBool(key string) (bool, error) {
		this.Lock()
		defer this.Unlock()
		return this.composer.GetBool(key)
}

func (this *Configuration) GetBoolDefault(key string, defaultValue bool) bool {
		this.Lock()
		defer this.Unlock()
		return this.composer.GetBoolDefault(key, defaultValue)
}

func (this *Configuration) GetFloat64(key string) (float64, error) {
		this.Lock()
		defer this.Unlock()
		return this.composer.GetFloat64(key)
}

func (this *Configuration) GetFloat64Default(key string, defaultValue float64) float64 {
		this.Lock()
		defer this.Unlock()
		return this.composer.GetFloat64Default(key, defaultValue)
}

func (this *Configuration) Set(key, val string) {
		this.Lock()
		defer this.Unlock()
		this.composer.Set(key, val)
}

func (this *Configuration) SetAll(values map[string]string) {
		this.Lock()
		defer this.Unlock()
		this.composer.SetAll(values)
}

func (this *Configuration) Keys() []string {
		this.Lock()
		defer this.Unlock()
		return this.composer.Keys()
}

func (this *Configuration) Unmarshal(t interface{}) error {
		this.Lock()
		defer this.Unlock()
		return this.composer.Unmarshal(t)
}
