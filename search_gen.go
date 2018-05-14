package main

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Search) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Items":
			var zb0002 uint32
			zb0002, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.Items == nil && zb0002 > 0 {
				z.Items = make(map[string]map[string]string, zb0002)
			} else if len(z.Items) > 0 {
				for key := range z.Items {
					delete(z.Items, key)
				}
			}
			for zb0002 > 0 {
				zb0002--
				var za0001 string
				var za0002 map[string]string
				za0001, err = dc.ReadString()
				if err != nil {
					return
				}
				var zb0003 uint32
				zb0003, err = dc.ReadMapHeader()
				if err != nil {
					return
				}
				if za0002 == nil && zb0003 > 0 {
					za0002 = make(map[string]string, zb0003)
				} else if len(za0002) > 0 {
					for key := range za0002 {
						delete(za0002, key)
					}
				}
				for zb0003 > 0 {
					zb0003--
					var za0003 string
					var za0004 string
					za0003, err = dc.ReadString()
					if err != nil {
						return
					}
					za0004, err = dc.ReadString()
					if err != nil {
						return
					}
					za0002[za0003] = za0004
				}
				z.Items[za0001] = za0002
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *Search) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "Items"
	err = en.Append(0x81, 0xa5, 0x49, 0x74, 0x65, 0x6d, 0x73)
	if err != nil {
		return
	}
	err = en.WriteMapHeader(uint32(len(z.Items)))
	if err != nil {
		return
	}
	for za0001, za0002 := range z.Items {
		err = en.WriteString(za0001)
		if err != nil {
			return
		}
		err = en.WriteMapHeader(uint32(len(za0002)))
		if err != nil {
			return
		}
		for za0003, za0004 := range za0002 {
			err = en.WriteString(za0003)
			if err != nil {
				return
			}
			err = en.WriteString(za0004)
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *Search) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "Items"
	o = append(o, 0x81, 0xa5, 0x49, 0x74, 0x65, 0x6d, 0x73)
	o = msgp.AppendMapHeader(o, uint32(len(z.Items)))
	for za0001, za0002 := range z.Items {
		o = msgp.AppendString(o, za0001)
		o = msgp.AppendMapHeader(o, uint32(len(za0002)))
		for za0003, za0004 := range za0002 {
			o = msgp.AppendString(o, za0003)
			o = msgp.AppendString(o, za0004)
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Search) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Items":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.Items == nil && zb0002 > 0 {
				z.Items = make(map[string]map[string]string, zb0002)
			} else if len(z.Items) > 0 {
				for key := range z.Items {
					delete(z.Items, key)
				}
			}
			for zb0002 > 0 {
				var za0001 string
				var za0002 map[string]string
				zb0002--
				za0001, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				var zb0003 uint32
				zb0003, bts, err = msgp.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				if za0002 == nil && zb0003 > 0 {
					za0002 = make(map[string]string, zb0003)
				} else if len(za0002) > 0 {
					for key := range za0002 {
						delete(za0002, key)
					}
				}
				for zb0003 > 0 {
					var za0003 string
					var za0004 string
					zb0003--
					za0003, bts, err = msgp.ReadStringBytes(bts)
					if err != nil {
						return
					}
					za0004, bts, err = msgp.ReadStringBytes(bts)
					if err != nil {
						return
					}
					za0002[za0003] = za0004
				}
				z.Items[za0001] = za0002
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *Search) Msgsize() (s int) {
	s = 1 + 6 + msgp.MapHeaderSize
	if z.Items != nil {
		for za0001, za0002 := range z.Items {
			_ = za0002
			s += msgp.StringPrefixSize + len(za0001) + msgp.MapHeaderSize
			if za0002 != nil {
				for za0003, za0004 := range za0002 {
					_ = za0004
					s += msgp.StringPrefixSize + len(za0003) + msgp.StringPrefixSize + len(za0004)
				}
			}
		}
	}
	return
}
