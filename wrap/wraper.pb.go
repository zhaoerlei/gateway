// Code generated by protoc-gen-go.
// source: wraper.proto
// DO NOT EDIT!

package wrap

import proto "github.com/golang/protobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type Msg_Type int32

const (
	Msg_UNKNOWN      Msg_Type = 0
	Msg_CLIENT_EVENT Msg_Type = 1
	Msg_USER_COUNT   Msg_Type = 2
)

var Msg_Type_name = map[int32]string{
	0: "UNKNOWN",
	1: "CLIENT_EVENT",
	2: "USER_COUNT",
}
var Msg_Type_value = map[string]int32{
	"UNKNOWN":      0,
	"CLIENT_EVENT": 1,
	"USER_COUNT":   2,
}

func (x Msg_Type) Enum() *Msg_Type {
	p := new(Msg_Type)
	*p = x
	return p
}
func (x Msg_Type) String() string {
	return proto.EnumName(Msg_Type_name, int32(x))
}
func (x *Msg_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Msg_Type_value, data, "Msg_Type")
	if err != nil {
		return err
	}
	*x = Msg_Type(value)
	return nil
}

type Cmd_Type int32

const (
	Cmd_UNKNOWN   Cmd_Type = 0
	Cmd_HANDSHAKE Cmd_Type = 1
	Cmd_PING      Cmd_Type = 2
	Cmd_PONG      Cmd_Type = 3
)

var Cmd_Type_name = map[int32]string{
	0: "UNKNOWN",
	1: "HANDSHAKE",
	2: "PING",
	3: "PONG",
}
var Cmd_Type_value = map[string]int32{
	"UNKNOWN":   0,
	"HANDSHAKE": 1,
	"PING":      2,
	"PONG":      3,
}

func (x Cmd_Type) Enum() *Cmd_Type {
	p := new(Cmd_Type)
	*p = x
	return p
}
func (x Cmd_Type) String() string {
	return proto.EnumName(Cmd_Type_name, int32(x))
}
func (x *Cmd_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Cmd_Type_value, data, "Cmd_Type")
	if err != nil {
		return err
	}
	*x = Cmd_Type(value)
	return nil
}

type Meta_Type int32

const (
	Meta_UNKNOWN Meta_Type = 0
)

var Meta_Type_name = map[int32]string{
	0: "UNKNOWN",
}
var Meta_Type_value = map[string]int32{
	"UNKNOWN": 0,
}

func (x Meta_Type) Enum() *Meta_Type {
	p := new(Meta_Type)
	*p = x
	return p
}
func (x Meta_Type) String() string {
	return proto.EnumName(Meta_Type_name, int32(x))
}
func (x *Meta_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Meta_Type_value, data, "Meta_Type")
	if err != nil {
		return err
	}
	*x = Meta_Type(value)
	return nil
}

// bag id = 0
type Msg struct {
	Event            *string   `protobuf:"bytes,1,req,name=event" json:"event,omitempty"`
	Tp               *Msg_Type `protobuf:"varint,2,req,name=tp,enum=wrap.Msg_Type,def=0" json:"tp,omitempty"`
	Meta             []*Meta   `protobuf:"bytes,3,rep,name=meta" json:"meta,omitempty"`
	UserCount        *int32    `protobuf:"varint,4,opt,name=userCount" json:"userCount,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *Msg) Reset()         { *m = Msg{} }
func (m *Msg) String() string { return proto.CompactTextString(m) }
func (*Msg) ProtoMessage()    {}

const Default_Msg_Tp Msg_Type = Msg_UNKNOWN

func (m *Msg) GetEvent() string {
	if m != nil && m.Event != nil {
		return *m.Event
	}
	return ""
}

func (m *Msg) GetTp() Msg_Type {
	if m != nil && m.Tp != nil {
		return *m.Tp
	}
	return Default_Msg_Tp
}

func (m *Msg) GetMeta() []*Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

func (m *Msg) GetUserCount() int32 {
	if m != nil && m.UserCount != nil {
		return *m.UserCount
	}
	return 0
}

// bag id = 1
type Cmd struct {
	Tp               *Cmd_Type `protobuf:"varint,1,req,name=tp,enum=wrap.Cmd_Type" json:"tp,omitempty"`
	Ct               *int64    `protobuf:"varint,2,req,name=ct" json:"ct,omitempty"`
	Txt              *string   `protobuf:"bytes,3,opt,name=txt" json:"txt,omitempty"`
	Meta             []*Meta   `protobuf:"bytes,4,rep,name=meta" json:"meta,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}

func (m *Cmd) Reset()         { *m = Cmd{} }
func (m *Cmd) String() string { return proto.CompactTextString(m) }
func (*Cmd) ProtoMessage()    {}

func (m *Cmd) GetTp() Cmd_Type {
	if m != nil && m.Tp != nil {
		return *m.Tp
	}
	return Cmd_UNKNOWN
}

func (m *Cmd) GetCt() int64 {
	if m != nil && m.Ct != nil {
		return *m.Ct
	}
	return 0
}

func (m *Cmd) GetTxt() string {
	if m != nil && m.Txt != nil {
		return *m.Txt
	}
	return ""
}

func (m *Cmd) GetMeta() []*Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

// meta, add invisible payload to message & command
type Meta struct {
	Tp               *Meta_Type `protobuf:"varint,1,req,name=tp,enum=wrap.Meta_Type" json:"tp,omitempty"`
	Txt              *string    `protobuf:"bytes,2,req,name=txt" json:"txt,omitempty"`
	XXX_unrecognized []byte     `json:"-"`
}

func (m *Meta) Reset()         { *m = Meta{} }
func (m *Meta) String() string { return proto.CompactTextString(m) }
func (*Meta) ProtoMessage()    {}

func (m *Meta) GetTp() Meta_Type {
	if m != nil && m.Tp != nil {
		return *m.Tp
	}
	return Meta_UNKNOWN
}

func (m *Meta) GetTxt() string {
	if m != nil && m.Txt != nil {
		return *m.Txt
	}
	return ""
}

func init() {
	proto.RegisterEnum("wrap.Msg_Type", Msg_Type_name, Msg_Type_value)
	proto.RegisterEnum("wrap.Cmd_Type", Cmd_Type_name, Cmd_Type_value)
	proto.RegisterEnum("wrap.Meta_Type", Meta_Type_name, Meta_Type_value)
}