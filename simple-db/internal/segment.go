package internal

import (
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"time"
)

const SEGMENT_NAME = "segment-%d.bin"

type Segment struct {
	file      *os.File
	cursor    int
	increment int
}

func NewSegment(index int, basePath string) *Segment {
	segment := Segment{}
	segment.open(index, basePath, os.O_APPEND|os.O_RDWR|os.O_CREATE)
	return &segment
}

func (s *Segment) Open(index int, basePath string) error {
	return s.open(index, basePath, os.O_APPEND|os.O_RDWR)
}

func (s *Segment) open(index int, basePath string, behaviour int) error {
	segment_path := path.Join(basePath, fmt.Sprintf(SEGMENT_NAME, s.increment))
	file, err := os.OpenFile(segment_path, behaviour, 0644)
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

func (s *Segment) ReadAt(offset int64) (*Record, int64, error) {
	record := &Record{}

	container := make([]byte, record.Header())
	_, err := s.file.ReadAt(container, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("Can't read byte %w", err)
	}

	keyLength := binary.LittleEndian.Uint32(container[0:4])
	dataLength := binary.LittleEndian.Uint32(container[4:8])
	fullSize := uint32(record.Header()) + keyLength + dataLength

	recordRaw := make([]byte, fullSize)
	_, err = s.file.ReadAt(recordRaw, offset)
	record.Decode(recordRaw)
	return record, record.Size() + offset, nil

}

func (s *Segment) Close() error {
	if s.file != nil {
		return nil
	}
	return s.file.Close()
}
