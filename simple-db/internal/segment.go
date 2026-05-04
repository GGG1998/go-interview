package internal

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"time"
)

type Segment struct {
	file      *os.File
	increment int
}

func (s *Segment) open(index int) error {
	file, err := os.OpenFile(fmt.Sprintf("segment-%d.bin", s.increment), os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	s.file = file
	s.increment = index
	return nil
}

func (s *Segment) append(key string, record []byte) error {
	rec := NewRecord(time.Now().UnixNano(), []byte(key), record)
	_, err := s.file.Write(rec.Encode())
	if err != nil {
		return err
	}
	return nil
}

func (s *Segment) readOne() (*Record, error) {
	record := &Record{}
	container := make([]byte, 8)
	_, err := s.file.ReadAt(container, 0)
	if err != nil {
		return nil, errors.New("xyz")
	}

	keyLength := binary.LittleEndian.Uint32(container[0:4])
	dataLength := binary.LittleEndian.Uint32(container[4:8])
	fullSize := 8 + keyLength + dataLength

	recordRaw := make([]byte, fullSize)
	_, err = s.file.ReadAt(recordRaw, 0)
	record.Decode(recordRaw)
	return record, nil

}

func (s *Segment) close() error {
	return s.file.Close()
}
