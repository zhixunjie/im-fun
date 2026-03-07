package model

// Proto proto.
type Proto struct {
	PackLen   int32  // package length
	HeaderLen int16  // header length
	Ver       int16  // protocol version
	Op        int32  // operation for request
	Seq       int32  // sequence number chosen by client
	Body      []byte // body
	BodyLen   int32  // body length
}
