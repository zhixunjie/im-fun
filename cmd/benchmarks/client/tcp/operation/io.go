package operation

import (
	"bufio"
	"encoding/binary"
	"github.com/zhixunjie/im-fun/benchmarks/client/tcp/model"
)

func WriteProto(wr *bufio.Writer, proto *model.Proto) (err error) {
	rawHeaderLen := model.RawHeaderLen
	// write header
	if err = binary.Write(wr, binary.BigEndian, uint32(rawHeaderLen)+uint32(len(proto.Body))); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, rawHeaderLen); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, proto.Ver); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, proto.Op); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, proto.Seq); err != nil {
		return
	}

	// write body
	if proto.Body != nil {
		if err = binary.Write(wr, binary.BigEndian, proto.Body); err != nil {
			return
		}
	}
	err = wr.Flush()
	return
}

func ReadProto(rd *bufio.Reader, proto *model.Proto) (err error) {
	// ready header
	if err = binary.Read(rd, binary.BigEndian, &proto.PackLen); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &proto.HeaderLen); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &proto.Ver); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &proto.Op); err != nil {
		return
	}
	if err = binary.Read(rd, binary.BigEndian, &proto.Seq); err != nil {
		return
	}

	// read body
	var bodyLen = int(proto.PackLen - int32(proto.HeaderLen))
	if bodyLen > 0 {
		proto.BodyLen = int32(bodyLen)
		proto.Body = make([]byte, bodyLen)

		// begin to read
		var n, t int
		for {
			if t, err = rd.Read(proto.Body[n:]); err != nil {
				return
			}
			if n += t; n == bodyLen {
				break
			}
		}
	} else {
		proto.Body = nil
	}
	return
}
