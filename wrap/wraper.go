package wrap

import (
	"github.com/golang/protobuf/proto"
	"fmt"
)

/*
在包里加入bagId之后在反序列化的时候可以直接判断bagId
protobuf 不支持多类型传输只能 加一个bagId用于解包的时候区分
既然用了晓光的多种类型传输 bagId确实就得加了
*/
const (
	BAG_ID_MSG byte = iota //bag id in for message 0
	BAG_ID_CMD             // 1
)

func CreateCmd(tp Cmd_Type, txt string, tick int64) *Cmd {
	cmd := &Cmd{
		Tp:  &tp, //不知为什么全都是指针……
		Ct:  proto.Int64(tick),
		Txt: proto.String(txt),
	}
	return cmd
}

func CreateMeta(tp Meta_Type, txt string) *Meta {
	meta := &Meta{
		Tp:  &tp,
		Txt: proto.String(txt),
	}
	return meta
}

func CreateMessage(tp Msg_Type, event string) *Msg {
	msg := &Msg{
		Tp:    &tp,
		Event: proto.String(event),
	}
	return msg
}

func CreateUserCountMsg(roomName string, userCount int) *Msg {
	msgType := Msg_USER_COUNT
	msg := &Msg{
		Tp:        &msgType,
		Event:     proto.String(roomName),
		UserCount: proto.Int(userCount),
	}
	return msg
}

func Boxing(m proto.Message) ([]byte, error) {
	b, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	var bagId byte
	switch m.(type) {
	case *Msg:
		bagId = BAG_ID_MSG
	case *Cmd:
		bagId = BAG_ID_CMD
	default:
		return nil, fmt.Errorf("unknown bagId")
	}
	return append([]byte{bagId}, b...), nil
}

func Unboxing(data []byte) (interface{}, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("warning: bad frame")
	}
	bagId := data[0]
	data = data[1:]

	switch bagId {
	case BAG_ID_MSG:
		if len(data) < 1 {
			return nil, fmt.Errorf("warning: bad message frame")
		}
		upMsg := new(Msg)
		if err := proto.Unmarshal(data, upMsg); err != nil {
			return nil, fmt.Errorf("parse msg err %v", err)
		}
		return upMsg, nil
	case BAG_ID_CMD:
		upCmd := new(Cmd)
		if err := proto.Unmarshal(data, upCmd); err != nil {
			return nil, fmt.Errorf("parse cmd err %v", err)
		}
		return upCmd, nil
	default:
		return nil, fmt.Errorf("warning: unknown bag id: %v", bagId)
	}
}
