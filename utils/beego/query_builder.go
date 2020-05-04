package beego

import (
		"errors"
		"github.com/WebGameLinux/cms/utils/reflects"
		"github.com/astaxie/beego/orm"
		"reflect"
)

type SqlQueryBuilder interface {
		Select(fields ...string) SqlQueryBuilder
		ForUpdate() SqlQueryBuilder
		From(tables ...string) SqlQueryBuilder
		InnerJoin(table string) SqlQueryBuilder
		LeftJoin(table string) SqlQueryBuilder
		RightJoin(table string) SqlQueryBuilder
		On(cond string) SqlQueryBuilder
		Where(cond string) SqlQueryBuilder
		And(cond string) SqlQueryBuilder
		Or(cond string) SqlQueryBuilder
		In(arr ...string) SqlQueryBuilder
		OrderBy(fields ...string) SqlQueryBuilder
		Asc() SqlQueryBuilder
		Desc() SqlQueryBuilder
		Limit(limit int) SqlQueryBuilder
		Offset(offset int) SqlQueryBuilder
		GroupBy(fields ...string) SqlQueryBuilder
		Having(cond string) SqlQueryBuilder
		Update(tables ...string) SqlQueryBuilder
		Set(kv ...string) SqlQueryBuilder
		Delete(tables ...string) SqlQueryBuilder
		InsertInto(table string, fields ...string) SqlQueryBuilder
		Values(values ...string) SqlQueryBuilder
		Subquery(sub string, alias string) string
		String() string
		First(v ...interface{}) interface{}
		Last(v ...interface{}) interface{}
		One(v ...interface{}) interface{}
		GetError(clean ...bool) error
		Get(v ...interface{}) interface{}
		SetModel(m interface{}) SqlQueryBuilder
		Paginator(page int, count int, columns []string, args ...interface{}) (interface{}, bool)
}

type SqlQueryBuilderDto struct {
		Query orm.QueryBuilder
		Exec  orm.Ormer
		Model interface{}
		Error error
}

type SqlQueryBuilderWrapper struct {
		SqlQueryBuilderDto
}

func NewQueryBuilderWrapper(query orm.QueryBuilder, exec orm.Ormer, model ...interface{}) SqlQueryBuilder {
		var wrapper = new(SqlQueryBuilderWrapper)
		wrapper.Exec = exec
		wrapper.Query = query
		if len(model) > 0 {
				wrapper.Model = model[0]
		}
		return wrapper
}

func (this *SqlQueryBuilderWrapper) Select(fields ...string) SqlQueryBuilder {
		this.Query = this.Query.Select(fields...)
		return this
}

func (this *SqlQueryBuilderWrapper) ForUpdate() SqlQueryBuilder {
		this.Query = this.Query.ForUpdate()
		return this
}

func (this *SqlQueryBuilderWrapper) From(tables ...string) SqlQueryBuilder {
		this.Query = this.Query.From(tables...)
		return this
}

func (this *SqlQueryBuilderWrapper) InnerJoin(table string) SqlQueryBuilder {
		this.Query = this.Query.InnerJoin(table)
		return this
}

func (this *SqlQueryBuilderWrapper) LeftJoin(table string) SqlQueryBuilder {
		this.Query = this.Query.LeftJoin(table)
		return this
}

func (this *SqlQueryBuilderWrapper) RightJoin(table string) SqlQueryBuilder {
		this.Query = this.Query.RightJoin(table)
		return this
}

func (this *SqlQueryBuilderWrapper) On(cond string) SqlQueryBuilder {
		this.Query = this.Query.On(cond)
		return this
}

func (this *SqlQueryBuilderWrapper) Where(cond string) SqlQueryBuilder {
		this.Query = this.Query.Where(cond)
		return this
}

func (this *SqlQueryBuilderWrapper) And(cond string) SqlQueryBuilder {
		this.Query = this.Query.And(cond)
		return this
}

func (this *SqlQueryBuilderWrapper) Or(cond string) SqlQueryBuilder {
		this.Query = this.Query.Or(cond)
		return this
}

func (this *SqlQueryBuilderWrapper) In(arr ...string) SqlQueryBuilder {
		this.Query = this.Query.In(arr...)
		return this
}

func (this *SqlQueryBuilderWrapper) OrderBy(fields ...string) SqlQueryBuilder {
		this.Query = this.Query.OrderBy(fields...)
		return this
}

func (this *SqlQueryBuilderWrapper) Asc() SqlQueryBuilder {
		this.Query = this.Query.Asc()
		return this
}

func (this *SqlQueryBuilderWrapper) Desc() SqlQueryBuilder {
		this.Query = this.Query.Desc()
		return this
}

func (this *SqlQueryBuilderWrapper) Limit(limit int) SqlQueryBuilder {
		this.Query = this.Query.Limit(limit)
		return this
}

func (this *SqlQueryBuilderWrapper) Offset(offset int) SqlQueryBuilder {
		this.Query = this.Query.Offset(offset)
		return this
}

func (this *SqlQueryBuilderWrapper) GroupBy(fields ...string) SqlQueryBuilder {
		this.Query = this.Query.GroupBy(fields...)
		return this
}

func (this *SqlQueryBuilderWrapper) Having(cond string) SqlQueryBuilder {
		this.Query = this.Query.Having(cond)
		return this
}

func (this *SqlQueryBuilderWrapper) Update(tables ...string) SqlQueryBuilder {
		this.Query = this.Query.Update(tables...)
		return this
}

func (this *SqlQueryBuilderWrapper) Set(kv ...string) SqlQueryBuilder {
		this.Query = this.Query.Set(kv...)
		return this
}

func (this *SqlQueryBuilderWrapper) Delete(tables ...string) SqlQueryBuilder {
		this.Query = this.Query.Delete(tables...)
		return this
}

func (this *SqlQueryBuilderWrapper) InsertInto(table string, fields ...string) SqlQueryBuilder {
		this.Query = this.Query.InsertInto(table, fields...)
		return this
}

func (this *SqlQueryBuilderWrapper) Values(values ...string) SqlQueryBuilder {
		this.Query = this.Query.Values(values...)
		return this
}

func (this *SqlQueryBuilderWrapper) Subquery(sub string, alias string) string {
		return this.Subquery(sub, alias)
}

func (this *SqlQueryBuilderWrapper) String() string {
		return this.Query.String()
}

func (this *SqlQueryBuilderWrapper) First(v ...interface{}) interface{} {
		sql := this.Asc().Limit(1).Offset(0).String()
		if this.Model == nil || this.Exec == nil {
				return sql
		}
		if this.Model == nil {
				return this.Exec.Raw(sql, v...)
		}
		this.Error = this.Exec.Raw(sql, v...).QueryRow(this.Model)
		if this.Error == nil {
				return this.Model
		}
		return nil
}

func (this *SqlQueryBuilderWrapper) Last(v ...interface{}) interface{} {
		sql := this.Desc().Limit(1).Offset(0).String()
		if this.Model == nil || this.Exec == nil {
				return sql
		}
		if this.Model == nil {
				return this.Exec.Raw(sql, v...)
		}
		this.Error = this.Exec.Raw(sql, v...).QueryRow(this.Model)
		if this.Error == nil {
				return this.Model
		}
		return nil
}

func (this *SqlQueryBuilderWrapper) One(v ...interface{}) interface{} {
		sql := this.Limit(1).Offset(0).String()
		if this.Model == nil || this.Exec == nil {
				return sql
		}
		if this.Model == nil {
				return this.Exec.Raw(sql, v...)
		}
		this.Error = this.Exec.Raw(sql, v...).QueryRow(this.Model)
		if this.Error == nil {
				return this.Model
		}
		return nil
}

func (this *SqlQueryBuilderWrapper) Get(v ...interface{}) interface{} {
		sql := this.String()
		if sql == "" {
				this.Error = errors.New("empty sql")
				return nil
		}
		if this.Model == nil || this.Exec == nil {
				return sql
		}
		if this.Model == nil {
				return this.Exec.Raw(sql, v...)
		}
		models := this.NewModels()
		if models == nil {
				return this.Exec.Raw(sql, v...)
		}
		_, this.Error = this.Exec.Raw(sql, v...).QueryRows(models)
		if this.Error == nil {
				return models
		}
		return nil
}

func (this *SqlQueryBuilderWrapper) GetError(clean ...bool) error {
		if len(clean) == 0 {
				clean = append(clean, true)
		}
		if clean[0] {
				return this.Error
		}
		var err = this.Error
		this.Error = nil
		return err
}

func (this *SqlQueryBuilderWrapper) SetModel(m interface{}) SqlQueryBuilder {
		if m != nil {
				this.Model = m
		}
		return this
}

func (this *SqlQueryBuilderWrapper) NewModels() interface{} {
		if this.Model == nil {
				return nil
		}
		original := reflects.RealValue(reflect.ValueOf(this.Model))
		if original.Kind() == reflect.Slice {
				return this.Model
		}
		return reflect.MakeSlice(original.Type(), 1, 1).Interface()
}

func (this *SqlQueryBuilderWrapper) Paginator(page int, count int, columns []string, args ...interface{}) (interface{}, bool) {
		if len(columns) != 0 {
				this.Select(columns...)
		}
		rows := this.Limit(count).Offset((page - 1) * count).Get(args...)
		return rows, this.Error == nil
}
