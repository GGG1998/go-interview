package internal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Record struct {
	keyLength    uint32
	recordLength uint32
	Key          []byte
	Record       []byte
}

func NewRecord(timestamp int64, key []byte, record []byte) *Record {
	return &Record{
		Key:    key,
		Record: record,
	}
}

func (r *Record) Header() int64 { return 8 /* uint32 + uint32 */ }

func (r *Record) Encode() []byte {
	headerSize := len(r.Key) + len(r.Record)
	buf := bytes.NewBuffer(make([]byte, 0, headerSize))

	binary.Write(buf, binary.LittleEndian, uint32(len(r.Key)))
	binary.Write(buf, binary.LittleEndian, uint32(len(r.Record)))

	buf.Write(r.Key)
	buf.Write(r.Record)

	return buf.Bytes()
}

func (r *Record) Decode(data []byte) (int, error) {
	/*
		TODO:
		- Change NewBuffer to NewReader, learn why
		- Set limit for record/key for security reason, learn why
	*/
	fmt.Println("Decode enter")
	if len(data) < 8 {
		return 0, io.ErrUnexpectedEOF
	}

	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.LittleEndian, &r.keyLength); err != nil {
		return 0, err
	}
	fmt.Println("Decode bug1")
	if err := binary.Read(buf, binary.LittleEndian, &r.recordLength); err != nil {
		return 0, err
	}

	r.Key = make([]byte, r.keyLength)
	r.Record = make([]byte, r.recordLength)
	fmt.Println("Decode bug2 %d", r.recordLength)

	if _, err := buf.Read(r.Key); err != nil {
		return 0, err
	}
	fmt.Println("Decode almost")
	if _, err := buf.Read(r.Record); err != nil {
		return 0, err
	}

	// Return - buf.Len(), learn why
	fmt.Println("Decode Finish")
	return len(data), nil
}
