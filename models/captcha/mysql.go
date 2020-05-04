package captcha

import (
		"fmt"
		"github.com/WebGameLinux/cms/models"
		"github.com/WebGameLinux/cms/utils/captcha/drivers"
		"github.com/WebGameLinux/cms/utils/captcha/stores"
		"log"
		"time"
)

type MysqlCaptchaProvider struct {
		models.BaseWrapper
}

type MysqlCaptchaWrapper interface {
		stores.MysqlStoreProvider
		NewKey(string, string) string
}

type MysqlCaptchaModel struct {
		stores.CaptchaCodeDto
		models.SoftDeleteDate
}

func NewMysqlCaptchaWrapper() *MysqlCaptchaProvider {
		var wrapper = new(MysqlCaptchaProvider)
		wrapper.Model = new(MysqlCaptchaModel)
		return wrapper
}

func (this *MysqlCaptchaModel) TableName() string {
		return "captcha_code_log"
}

func (this *MysqlCaptchaProvider) Get(model *stores.CaptchaCodeDto) bool {
		var res interface{}
		if model == nil {
				return false
		}
		query, err := this.GetQuery()
		if err != nil {
				this.Error = err
				return false
		}
		if model.CaptchaId == "" {
				return false
		}
		query = query.Select("*").From(this.Table()).Where("captcha_id=?")
		if model.Key != "" {
				query = query.And("key=?")
		}
		m := new(MysqlCaptchaModel)
		query = query.And("expired_at > ?").SetModel(m).And(" deleted_at IS NULL ")
		if model.Key == "" {
				res = query.First(model.CaptchaId, time.Now())
		} else {
				res = query.First(model.CaptchaId, model.Key, time.Now())
		}
		this.Error = query.GetError()
		if res != nil && this.Error == nil {
				return true
		}
		return false
}

func (this *MysqlCaptchaProvider) Set(model *stores.CaptchaCodeDto) bool {
		if model == nil {
				return false
		}
		var (
				n int64
				m = new(MysqlCaptchaModel)
		)
		m.InitByDto(model)
		now := time.Now()
		if m.ExpireAt.Equal(time.Time{}) || now.Equal(m.ExpireAt) || now.After(m.ExpireAt) {
				m.ExpireAt = time.Now().Add(stores.DefaultExpireDuration)
		}
		if n, this.Error = this.GetOrm().Insert(m); this.Error == nil && n > 0 {
				return true
		}
		return false
}

func (this *MysqlCaptchaProvider) Delete(model *stores.CaptchaCodeDto) {
		var (
				n int64
				m = new(MysqlCaptchaModel)
		)
		if model.Id == 0 {
				if !this.Get(model) {
						if this.Error != nil {
								log.Fatal(this.GetError())
						}
						return
				}
		}
		m = m.InitByDto(model)
		m.DeletedAt = time.Now()
		if n, this.Error = this.GetOrm().Update(m, "deleted_at"); this.Error == nil && n > 0 {
				return
		}
		log.Fatal(this.Error)
}

func (this *MysqlCaptchaProvider) NewKey(mobile, code string) string {
		var key = ""
		query, err := this.GetQuery()
		if err != nil {
				this.Error = err
				return ""
		}
		sql := query.Select("*").From(this.Table()).Where("captcha_id=?").
				And("code=?").
				And("`key`=?").And("deleted_at is null").Limit(1).String()
		i := 3

		for i > 0 {
				u := new(MysqlCaptchaModel)
				key = fmt.Sprintf("%s%d", drivers.RandText(10), time.Now().Unix())
				if err = this.GetOrm().Raw(sql, key, code, mobile).QueryRow(u); err != nil {
						if u.Id == 0 {
								return key
						}
				}
				i--
		}

		return key
}

func (this *MysqlCaptchaModel) InitByDto(dto *stores.CaptchaCodeDto) *MysqlCaptchaModel {
		this.Id = dto.Id
		this.Key = dto.Key
		this.CaptchaId = dto.CaptchaId
		if dto.Code != "" {
				this.Code = dto.Code
		}
		return this
}
