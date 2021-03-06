// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package rabbitmq

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	time "time"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson24099d24DecodeGithubComWebGameLinuxCmsLibsRabbitmq(in *jlexer.Lexer, out *WorkPoolConfig) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		case "max_num":
			out.MaxNum = int(in.Int())
		case "interval":
			out.Interval = time.Duration(in.Int64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson24099d24EncodeGithubComWebGameLinuxCmsLibsRabbitmq(out *jwriter.Writer, in WorkPoolConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"max_num\":"
		out.RawString(prefix)
		out.Int(int(in.MaxNum))
	}
	{
		const prefix string = ",\"interval\":"
		out.RawString(prefix)
		out.Int64(int64(in.Interval))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WorkPoolConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson24099d24EncodeGithubComWebGameLinuxCmsLibsRabbitmq(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WorkPoolConfig) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson24099d24EncodeGithubComWebGameLinuxCmsLibsRabbitmq(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WorkPoolConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson24099d24DecodeGithubComWebGameLinuxCmsLibsRabbitmq(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WorkPoolConfig) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson24099d24DecodeGithubComWebGameLinuxCmsLibsRabbitmq(l, v)
}
func easyjson24099d24DecodeGithubComWebGameLinuxCmsLibsRabbitmq1(in *jlexer.Lexer, out *Options) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "queue":
			out.Queue = string(in.String())
		case "consumer":
			out.Consumer = string(in.String())
		case "auto_ack":
			out.AutoAck = bool(in.Bool())
		case "exclusive":
			out.Exclusive = bool(in.Bool())
		case "no_local":
			out.NoLocal = bool(in.Bool())
		case "no_wait":
			out.NoWait = bool(in.Bool())
		case "durable":
			out.Durable = bool(in.Bool())
		case "auto_delete":
			out.AutoDelete = bool(in.Bool())
		case "table":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.Args = make(map[string]interface{})
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 interface{}
					if m, ok := v1.(easyjson.Unmarshaler); ok {
						m.UnmarshalEasyJSON(in)
					} else if m, ok := v1.(json.Unmarshaler); ok {
						_ = m.UnmarshalJSON(in.Raw())
					} else {
						v1 = in.Interface()
					}
					(out.Args)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
			}
		case "mandatory":
			out.Mandatory = bool(in.Bool())
		case "immediate":
			out.Immediate = bool(in.Bool())
		case "internal":
			out.Internal = bool(in.Bool())
		case "work_mode":
			out.WorkMode = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson24099d24EncodeGithubComWebGameLinuxCmsLibsRabbitmq1(out *jwriter.Writer, in Options) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"queue\":"
		out.RawString(prefix[1:])
		out.String(string(in.Queue))
	}
	{
		const prefix string = ",\"consumer\":"
		out.RawString(prefix)
		out.String(string(in.Consumer))
	}
	{
		const prefix string = ",\"auto_ack\":"
		out.RawString(prefix)
		out.Bool(bool(in.AutoAck))
	}
	{
		const prefix string = ",\"exclusive\":"
		out.RawString(prefix)
		out.Bool(bool(in.Exclusive))
	}
	{
		const prefix string = ",\"no_local\":"
		out.RawString(prefix)
		out.Bool(bool(in.NoLocal))
	}
	{
		const prefix string = ",\"no_wait\":"
		out.RawString(prefix)
		out.Bool(bool(in.NoWait))
	}
	{
		const prefix string = ",\"durable\":"
		out.RawString(prefix)
		out.Bool(bool(in.Durable))
	}
	{
		const prefix string = ",\"auto_delete\":"
		out.RawString(prefix)
		out.Bool(bool(in.AutoDelete))
	}
	{
		const prefix string = ",\"table\":"
		out.RawString(prefix)
		if in.Args == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v2First := true
			for v2Name, v2Value := range in.Args {
				if v2First {
					v2First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v2Name))
				out.RawByte(':')
				if m, ok := v2Value.(easyjson.Marshaler); ok {
					m.MarshalEasyJSON(out)
				} else if m, ok := v2Value.(json.Marshaler); ok {
					out.Raw(m.MarshalJSON())
				} else {
					out.Raw(json.Marshal(v2Value))
				}
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"mandatory\":"
		out.RawString(prefix)
		out.Bool(bool(in.Mandatory))
	}
	{
		const prefix string = ",\"immediate\":"
		out.RawString(prefix)
		out.Bool(bool(in.Immediate))
	}
	{
		const prefix string = ",\"internal\":"
		out.RawString(prefix)
		out.Bool(bool(in.Internal))
	}
	{
		const prefix string = ",\"work_mode\":"
		out.RawString(prefix)
		out.String(string(in.WorkMode))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Options) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson24099d24EncodeGithubComWebGameLinuxCmsLibsRabbitmq1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Options) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson24099d24EncodeGithubComWebGameLinuxCmsLibsRabbitmq1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Options) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson24099d24DecodeGithubComWebGameLinuxCmsLibsRabbitmq1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Options) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson24099d24DecodeGithubComWebGameLinuxCmsLibsRabbitmq1(l, v)
}
func easyjson24099d24DecodeGithubComWebGameLinuxCmsLibsRabbitmq2(in *jlexer.Lexer, out *ConnOptions) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "username":
			out.Username = string(in.String())
		case "password":
			out.Password = string(in.String())
		case "host":
			out.Host = string(in.String())
		case "port":
			out.Port = int(in.Int())
		case "virtual_host":
			out.VirtualHost = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson24099d24EncodeGithubComWebGameLinuxCmsLibsRabbitmq2(out *jwriter.Writer, in ConnOptions) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix[1:])
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	{
		const prefix string = ",\"host\":"
		out.RawString(prefix)
		out.String(string(in.Host))
	}
	{
		const prefix string = ",\"port\":"
		out.RawString(prefix)
		out.Int(int(in.Port))
	}
	{
		const prefix string = ",\"virtual_host\":"
		out.RawString(prefix)
		out.String(string(in.VirtualHost))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ConnOptions) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson24099d24EncodeGithubComWebGameLinuxCmsLibsRabbitmq2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ConnOptions) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson24099d24EncodeGithubComWebGameLinuxCmsLibsRabbitmq2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ConnOptions) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson24099d24DecodeGithubComWebGameLinuxCmsLibsRabbitmq2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ConnOptions) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson24099d24DecodeGithubComWebGameLinuxCmsLibsRabbitmq2(l, v)
}
func easyjson24099d24DecodeGithubComWebGameLinuxCmsLibsRabbitmq3(in *jlexer.Lexer, out *ConfigObject) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "mode":
			out.Mode = string(in.String())
		case "conn":
			out.Conn = string(in.String())
		case "options":
			out.Options = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson24099d24EncodeGithubComWebGameLinuxCmsLibsRabbitmq3(out *jwriter.Writer, in ConfigObject) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"mode\":"
		out.RawString(prefix[1:])
		out.String(string(in.Mode))
	}
	{
		const prefix string = ",\"conn\":"
		out.RawString(prefix)
		out.String(string(in.Conn))
	}
	{
		const prefix string = ",\"options\":"
		out.RawString(prefix)
		out.String(string(in.Options))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ConfigObject) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson24099d24EncodeGithubComWebGameLinuxCmsLibsRabbitmq3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ConfigObject) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson24099d24EncodeGithubComWebGameLinuxCmsLibsRabbitmq3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ConfigObject) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson24099d24DecodeGithubComWebGameLinuxCmsLibsRabbitmq3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ConfigObject) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson24099d24DecodeGithubComWebGameLinuxCmsLibsRabbitmq3(l, v)
}
