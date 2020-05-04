package types

import (
		"encoding/json"
		"fmt"
		"github.com/WebGameLinux/cms/utils/array"
		"github.com/WebGameLinux/cms/utils/mapper"
)

type Paginator interface {
		Items() []interface{}
		Count() int
		Len() int
		fmt.Stringer
		Meta() PaginatorMeta
		Init()
		Current() int
		Map() mapper.Mapper
		Dto() *PaginatorDto
}

type PaginatorWrapper interface {
		Paginator
		Next() (Paginator, bool)
		Prev() (Paginator, bool)
		More() bool
		Total() int
		Store(key string, value interface{})
		Get(key string, def ...interface{}) interface{}
}

const ProviderName = "paginator"
const ProviderMetaName = "paginatorMeta"

type PaginatorMeta interface {
		Url() string
		Get(key string) string
		Var(key string, def ...string) string
		Set(key string, value string)
		fmt.Stringer
		Map() map[string]string
		Json() string
}

type MetaDto struct {
		Uri  string
		Args map[string]string
}

type MetaWrapper struct {
		MetaDto
		MetaProvider Provider
		UrlHandler   func(wrapper *MetaWrapper) string
		VarsHandler  func(key string, wrapper *MetaWrapper) string
}

func NewMeta() PaginatorMeta {
		var meta = new(MetaWrapper)
		meta.MetaProvider, _ = Resolver(ProviderMetaName)
		return meta
}

func (this *MetaWrapper) Url() string {
		if this.UrlHandler == nil && this.MetaProvider != nil {
				if v, ok := this.MetaProvider.Get("url"); ok {
						if handler, ok := v.(func(interface{}) string); ok {
								this.UrlHandler = func(wrapper *MetaWrapper) string {
										return handler(wrapper)
								}
						}
				}
		}
		if this.UrlHandler != nil {
				return this.UrlHandler(this)
		}
		return this.getUrl()
}

func (this *MetaWrapper) getUrl() string {
		var (
				schema = this.Var("schema", "http://")
				domain = this.Var("domain", "localhost")
				port   = this.Var("port", "80")
				path   = this.Var("path", "/")
				page   = this.Var("page", "1")
				count  = this.Var("count", "20")
		)
		return fmt.Sprintf("%s%s:%s/%spage=%s&count=%s", schema, domain, port, path, page, count)
}

func (this *MetaWrapper) Get(key string) string {
		return this.Args[key]
}

func (this *MetaWrapper) Var(key string, def ...string) string {
		if len(def) == 0 {
				def = append(def, "")
		}
		if this.VarsHandler == nil && this.MetaProvider != nil {
				if v, ok := this.MetaProvider.Get("var"); ok {
						if handler, ok := v.(func(string, interface{}) string); ok {
								this.VarsHandler = func(key string, wrapper *MetaWrapper) string {
										return handler(key, wrapper)
								}
						}
				}
		}
		if this.VarsHandler != nil {
				if val := this.VarsHandler(key, this); val != "" {
						return val
				}
		}
		if v, ok := this.Args[key]; ok {
				return v
		}
		return def[0]
}

func (this *MetaWrapper) Set(key string, value string) {
		this.Args[key] = value
}

func (this *MetaWrapper) String() string {
		if this.MetaProvider != nil {
				if v, ok := this.MetaProvider.Get("stringer"); ok {
						if handler, ok := v.(func(interface{}) string); ok {
								return handler(this)
						}
				}
		}
		if v, err := json.Marshal(this.Map()); err == nil {
				return string(v)
		}
		return ""
}

func (this *MetaWrapper) Json() string {
		if v, err := json.Marshal(this.Map()); err == nil {
				return string(v)
		}
		return ""
}

func (this *MetaWrapper) Map() map[string]string {
		this.Args["url"] = this.Url()
		return this.Args
}

type PaginatorDto struct {
		Items []interface{}
		Count int
		Page  int
		Size  int
		Meta  PaginatorMeta
}

type BasePaginator struct {
		PaginatorDto
		NextHandler      func(*BasePaginator) (Paginator, bool) // 下一页
		PrevHandler      func(*BasePaginator) (Paginator, bool) // 上一页
		MoreHandler      func(*BasePaginator) bool              // 是否还有更多
		TotalHandler     func(*BasePaginator) int               // 总共页数
		MarshalHandler   func(interface{}) error                // 反序列化
		UnmarshalHandler func(*BasePaginator) (string, error)   // 序列化
		Boot             func()                                 // 初始化
		Provider         Provider                               // 服务提供器
		StoreMap         map[string]interface{}                 // 数据缓存
}

func (this *BasePaginator) Next() (Paginator, bool) {
		if this.NextHandler != nil {
				return this.NextHandler(this)
		}
		if this.Provider == nil {
				return nil, false
		}
		if v, ok := this.Provider.Get("next"); ok {
				if handler, ok := v.(func(interface{}) (Paginator, bool)); ok {
						this.NextHandler = func(paginator2 *BasePaginator) (paginator Paginator, b bool) {
								return handler(paginator)
						}
				}
		}
		if this.NextHandler != nil {
				return this.NextHandler(this)
		}
		return nil, false
}

func (this *BasePaginator) Prev() (Paginator, bool) {
		if this.PrevHandler != nil {
				return this.PrevHandler(this)
		}
		if this.Provider == nil {
				return nil, false
		}
		if v, ok := this.Provider.Get("prev"); ok {
				if handler, ok := v.(func(interface{}) (Paginator, bool)); ok {
						this.PrevHandler = func(paginator2 *BasePaginator) (paginator Paginator, b bool) {
								return handler(paginator)
						}
				}
		}
		if this.PrevHandler != nil {
				return this.PrevHandler(this)
		}
		return nil, false
}

func (this *BasePaginator) More() bool {
		if this.MoreHandler != nil {
				return this.MoreHandler(this)
		}
		if this.Provider == nil {
				return false
		}
		if v, ok := this.Provider.Get("more"); ok {
				if handler, ok := v.(func(interface{}) bool); ok {
						this.MoreHandler = func(paginator *BasePaginator) bool {
								return handler(paginator)
						}
				}
		}
		if this.MoreHandler != nil {
				return this.MoreHandler(this)
		}
		return false
}

func (this *BasePaginator) Total() int {
		if this.TotalHandler != nil {
				return this.TotalHandler(this)
		}
		if this.Provider == nil {
				return this.Len()
		}
		if v, ok := this.Provider.Get("total"); ok {
				if handler, ok := v.(func(interface{}) int); ok {
						this.TotalHandler = func(paginator *BasePaginator) int {
								return handler(paginator)
						}
				}
		}
		if this.TotalHandler != nil {
				return this.TotalHandler(this)
		}
		return this.Len()
}

func NewBasePaginator() Paginator {
		var paginator = new(BasePaginator)
		paginator.Provider, _ = Resolver(ProviderName)
		return paginator
}

func GetPaginatorProvider() Provider {
		if v, ok := Resolver(ProviderName); ok {
				return v
		}
		return nil
}

func Paginator2BaseProvider(p Paginator) *BasePaginator {
		if base, ok := p.(*BasePaginator); ok {
				return base
		}
		return nil
}

func (this *BasePaginator) Booted() bool {
		if this.Boot != nil {
				if this.NextHandler != nil && this.PrevHandler != nil {
						return true
				}
		}
		if this.Provider != nil {
				if v, ok := this.Provider.Get("boot"); ok {
						if fn, ok := v.(func(interface{})); ok {
								this.Boot = func() {
										fn(this)
								}
						}
				}
		}
		if this.Boot != nil {
				this.Boot()
		}
		return true
}

func (this *BasePaginator) Init() {
		if this.Booted() {
				return
		}
		if this.Boot != nil {
				this.Boot()
		}
}

func (this *BasePaginator) Items() []interface{} {
		return this.PaginatorDto.Items
}

func (this *BasePaginator) Count() int {
		return this.PaginatorDto.Count
}

func (this *BasePaginator) Len() int {
		return len(this.PaginatorDto.Items)
}

func (this *BasePaginator) Current() int {
		if this.Page == 0 {
				return 1
		}
		return this.Page
}

func (this *BasePaginator) String() string {
		if this.UnmarshalHandler == nil && this.Provider != nil {
				if v, ok := this.Provider.Get("unmarshal"); ok {
						if handler, ok := v.(func(interface{}) (string, error)); ok {
								this.UnmarshalHandler = func(paginator *BasePaginator) (s string, err error) {
										return handler(paginator)
								}
						}
				}
		}
		if this.UnmarshalHandler != nil {
				if v, err := this.UnmarshalHandler(this); err == nil {
						return v
				}
		}

		if v, err := json.Marshal(this.Map()); err == nil {
				return string(v)
		}
		return ""
}

func (this *BasePaginator) Meta() PaginatorMeta {
		if this.PaginatorDto.Meta == nil && this.Provider != nil {
				if v, ok := this.Provider.Get("meta"); ok {
						if handler, ok := v.(func(interface{}) PaginatorMeta); ok {
								this.PaginatorDto.Meta = handler(this)
						}
				}
		}
		return this.PaginatorDto.Meta
}

func (this *BasePaginator) Store(key string, value interface{}) {
		switch key {
		case "count":
				if n, ok := value.(int); ok {
						this.PaginatorDto.Count = n
				}
				return
		case "page":
				if n, ok := value.(int); ok {
						this.PaginatorDto.Page = n
				}
				return
		case "items":
				if n, ok := value.([]interface{}); ok {
						this.PaginatorDto.Items = n
				}
				if array.IsArray(value) {
						this.PaginatorDto.Items = array.JoinArrays(this.PaginatorDto.Items, value)
				}
				return
		default:
				if this.Provider != nil {
						if v, ok := this.Provider.Get("store.set"); ok {
								if handler, ok := v.(func(interface{}, interface{}, interface{})); ok {
										handler(this, key, value)
										return
								}
						}
				}
		}
		if this.StoreMap == nil {
				this.StoreMap = make(mapper.Mapper)
		}
		this.StoreMap[key] = value
}

func (this *BasePaginator) Get(key string, def ...interface{}) interface{} {
		if len(def) == 0 {
				def = append(def, nil)
		}
		switch key {
		case "count":
				return this.Count()
		case "page":
				return this.Current()
		case "items":
				return this.Items()
		default:
				if this.Provider != nil {
						if v, ok := this.Provider.Get("store.get"); ok {
								if handler, ok := v.(func(interface{}, interface{}) (interface{}, bool)); ok {
										if d, ok := handler(this, key); ok {
												return d
										}
								}
						}
				}
		}
		if this.StoreMap == nil {
				return def[0]
		}
		if v, ok := this.StoreMap[key]; ok {
				return v
		}
		return def[0]
}

func (this *BasePaginator) Map() mapper.Mapper {
		var m = make(map[string]interface{})
		m["page"] = this.Current()
		m["size"] = this.Len()
		m["items"] = this.Items()
		m["count"] = this.Count()
		m["meta"] = this.Meta()
		return m
}

func (this *BasePaginator) Dto() *PaginatorDto {
		var dto = new(PaginatorDto)
		dto.Items = this.Items()
		dto.Count = this.Count()
		dto.Size = this.Len()
		dto.Meta = this.Meta()
		return dto
}
