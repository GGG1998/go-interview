package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"time"

	"simple-db/internal"
)

type Segment struct {
	file  *os.File
	index int
}

func (s *Segment) open(index int) error {
	file, err := os.OpenFile(fmt.Sprintf("segment-%d.bin", s.index), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	s.file = file
	s.index = index
	return nil
}

func (s *Segment) append(key []byte, record []byte) error {
	rec := internal.NewRecord(time.Now().UnixNano(), key, record)
	_, err := s.file.Write(rec.Encode())
	if err != nil {
		return err
	}
	return nil
}

func (s *Segment) readOne() (*internal.Record, error) {
	record := &internal.Record{}
	container := make([]byte, 8)
	_, err := s.file.ReadAt(container, 0)
	if err != nil {
		return nil, errors.New("xyz")
	}
	fmt.Println("XXX")
	keyLength := binary.LittleEndian.Uint32(container[0:4])
	fmt.Println("key", keyLength)
	dataLength := binary.LittleEndian.Uint32(container[4:8])
	fmt.Println("Data", dataLength)
	fullSize := 8 + keyLength + dataLength

	recordRaw := make([]byte, fullSize)
	_, err = s.file.ReadAt(recordRaw, 0)
	record.Decode(recordRaw)
	return record, nil

}

func (s *Segment) close() error {
	return s.file.Close()
}

func main() {
	segment := &Segment{}
	segment.open(0)
	// segment.append([]byte("hello3"), []byte("3world"))
	r, _ := segment.readOne()
	fmt.Println(string(r.Key))
	segment.close()
}
