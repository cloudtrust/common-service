// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package fb

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Tuple struct {
	_tab flatbuffers.Table
}

func GetRootAsTuple(buf []byte, offset flatbuffers.UOffsetT) *Tuple {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Tuple{}
	x.Init(buf, n+offset)
	return x
}

func FinishTupleBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.Finish(offset)
}

func GetSizePrefixedRootAsTuple(buf []byte, offset flatbuffers.UOffsetT) *Tuple {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &Tuple{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func FinishSizePrefixedTupleBuffer(builder *flatbuffers.Builder, offset flatbuffers.UOffsetT) {
	builder.FinishSizePrefixed(offset)
}

func (rcv *Tuple) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Tuple) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Tuple) Key() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Tuple) Value() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func TupleStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func TupleAddKey(builder *flatbuffers.Builder, key flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(key), 0)
}
func TupleAddValue(builder *flatbuffers.Builder, value flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(value), 0)
}
func TupleEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
