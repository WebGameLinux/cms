package models

import (
		"errors"
		"github.com/WebGameLinux/cms/models/conditions"
		"github.com/WebGameLinux/cms/models/types"
		"github.com/WebGameLinux/cms/utils/reflects"
		string2 "github.com/WebGameLinux/cms/utils/string"
		"github.com/astaxie/beego/orm"
		"github.com/sirupsen/logrus"
		"time"
)

type AttachmentsWrapper struct {
		BaseWrapper
}

func NewAttachmentsWrapper() *AttachmentsWrapper {
		wrapper := new(AttachmentsWrapper)
		wrapper.Model = new(Attachments)
		return wrapper
}

func (this *Attachments) TableName() string {
		return `attachments`
}

func GetAttachments(options ...interface{}) *AttachmentsWrapper {
		wrapper := NewAttachmentsWrapper()
		WrapperInitOptions(wrapper, options...)
		return wrapper
}

func (this *AttachmentsWrapper) GetById(id int64) *Attachments {
		model := new(Attachments)
		model.Id = id
		err := this.GetOrm().Read(model)
		if err == orm.ErrNoRows {
				this.Error = err
				return nil
		}
		return model
}

func (this *AttachmentsWrapper) Save(data *Attachments) int64 {
		if !this.check(data) {
				return 0
		}
		return this.create(data)
}

func (this *AttachmentsWrapper) DeleteById(id int64) bool {
		data := new(Attachments)
		data.Id = id
		data.DeletedAt = time.Now()
		if _, err := this.GetOrm().Update(data, "deleted_at"); err != nil {
				return false
		}
		return true
}

func (this *AttachmentsWrapper) DestroyById(id int64) bool {
		data := new(Attachments)
		data.Id = id
		if _, err := this.GetOrm().Delete(data); err != nil {
				return false
		}
		return true
}

func (this *AttachmentsWrapper) Add(attachments []*Attachments) (int, error) {
		var items []*Attachments
		for _, data := range attachments {
				if !this.check(data) {
						continue
				}
				if !this.InitAttachment(data) {
						continue
				}
				items = append(items, data)
		}
		if len(items) == 0 {
				return 0, errors.New("all not allow")
		}
		num, err := this.GetOrm().InsertMulti(len(items), items)
		return int(num), err
}

func (this AttachmentsWrapper) UpdateById(id int64, data map[string]interface{}) bool {
		if len(data) == 0 {
				return false
		}
		attachment := this.GetById(id)
		if attachment == nil {
				return false
		}
		// 无效更新
		if !this.differUpdate(attachment, data) {
				return false
		}
		if rows, err := this.GetOrm().Update(attachment); err == nil && rows > 0 {
				return true
		}
		return false
}

func (this AttachmentsWrapper) differUpdate(data *Attachments, update map[string]interface{}) bool {
		var num = 0
		for key, v := range update {
				switch key {
				case `relate_table`:
						if table, ok := v.(string); ok {
								if data.RelateTable != table {
										data.RelateTable = table
										num++
								}
						}
				case `relate_id`:
						if id, ok := v.(int); ok {
								if data.RelateId != id {
										data.RelateId = id
										num++
								}
						}
						if id, ok := v.(int64); ok {
								if data.RelateId != int(id) {
										data.RelateId = int(id)
										num++
								}
						}
				case `access_url`:
						if url, ok := v.(string); ok {
								if data.AccessUrl != url && url != "" {
										data.AccessUrl = url
										num++
								}
						}
				case `save_path`:
						if save, ok := v.(string); ok {
								if data.SavePath != save && save != "" {
										data.SavePath = save
										num++
								}
						}
				case `uploader_id`:
						if id, ok := v.(int64); ok {
								if data.UploaderId != id && id != 0 {
										data.UploaderId = id
										num++
								}
						}
				case `deleter_id`:
						if id, ok := v.(int64); ok {
								if data.DeleterId != id && id != 0 {
										data.DeleterId = id
										num++
								}
						}
				case `size`:
						if size, ok := v.(string); ok {
								if data.Size != size && size != "" {
										data.Size = size
										num++
								}
						}
				case `filename`:
						if name, ok := v.(string); ok {
								if data.FileName != name && name != "" {
										data.FileName = name
										num++
								}
						}
				case `file_type`:
						if ty, ok := v.(string); ok {
								if data.FileType != ty && ty != "" {
										data.FileType = ty
										num++
								}
						}
				case `extension`:
						if ext, ok := v.(string); ok {
								if data.Extension != ext && ext != "" {
										data.Extension = ext
										num++
								}
						}
				case `deleted_at`:
						if del, ok := v.(time.Time); ok {
								if data.DeletedAt != del {
										data.DeletedAt = del
										num++
								}
						}
				}
		}
		return num > 0
}

func (this *AttachmentsWrapper) create(data *Attachments) int64 {
		if !this.InitAttachment(data) {
				return 0
		}
		if id, err := this.GetOrm().Insert(data); err == nil {
				return id
		}
		return 0
}

func (this *AttachmentsWrapper) InitAttachment(data *Attachments) bool {
		if data.Hash == "" {
				data.Hash = string2.FileHash(data.SavePath)
		}
		if data.SeqId == "" {
				data.SeqId = this.CreateUUid(data.TableName(), SeqKey)
		}
		if data.Hash == "" {
				return false
		}
		return true
}

func (this *AttachmentsWrapper) check(data *Attachments) bool {
		if data == nil || data.Id != 0 {
				return false
		}
		if data.AccessUrl == "" || data.SavePath == "" {
				return false
		}
		return false
}

func (this AttachmentsWrapper) Lists(condition *conditions.PageCondition) (items []orm.ParamsList, meta *types.Meta) {
		meta = types.NewMeta()
		meta.Count = condition.Count
		meta.Page = condition.Page
		resolver := this.SearchQueryResolver(condition.Conditions)
		if resolver == nil {
				return
		}
		if len(resolver.Effects) == 0 {
				return
		}
		keys := resolver.Fields
		if len(keys) == 0 {
				keys = append(keys, "*")
		}
		r2 := this.GetOrm().QueryTable(this.Table())
		err := reflects.Copy(resolver.Query, r2)
		if err != nil {
				logrus.Debug(err)
		} else {
				meta.SetQuery(r2)
		}
		total, _ := resolver.Query.Count()
		meta.Total = int(total)
		if n, err := resolver.Query.Limit(meta.Limit()).Offset(meta.Offset()).ValuesList(&items, keys...); err == nil && n > 0 {
				return
		}
		return
}

func (this *AttachmentsWrapper) SearchQueryResolver(conditions map[string]interface{}) *types.SearchParams {
		builder := this.GetOrm().QueryTable(new(Attachments).TableName())
		var (
				fields  []string
				effects map[string]interface{}
		)
		if arr, ok := conditions[FieldsQueryKey]; ok {
				if v, ok := arr.([]string); ok {
						fields = this.FilterFields(v)
				}
		}
		builder, effects = this.QueryResolver(builder, this.GetFieldKeys(), conditions)
		params := types.NewSearchParams()
		params.Query = builder
		params.Fields = fields
		params.Effects = effects
		return params
}
