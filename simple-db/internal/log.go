package internal

import "iter"

const DEFAULT_PATH = "./"
const DEFAULT_SIZE_LOG = 1024 * 1024 * 256

type Log struct {
	basePath string
	maxSize  int64
	index    map[string]int64
	active   *Segment
}

func (l *Log) next() iter.Seq2[string, int64] {
	// It's require error management
	return func(yield func(K string, V int64) bool) {
		// We don't know how many files exist
		// We know that loop will be break by an error
		startSegmentIndex := 0
		var startOffset int64
		segment := &Segment{}

		for {
			if err := segment.Open(startSegmentIndex, l.basePath); err != nil {
				return
			}
			l.active = segment
			for {
				record, nextOffset, err := segment.ReadAt(startOffset)
				if err != nil {
					break
				}
				if !yield(string(record.Key), startOffset) {
					segment.Close()
					return
				}
				startOffset = nextOffset
			}
			startSegmentIndex += 1
			startOffset = 0
			segment.Close()
		}
	}
}

func (l *Log) buildIndex() {
	l.index = make(map[string]int64)
	for k, v := range l.next() {
		l.index[k] = v
	}
}

func (l *Log) Open() {
	l.buildIndex()
	if l.active == nil {
		l.active = NewSegment(0, l.basePath)
	}
}

func (l *Log) Append(key string, record []byte) error {
	if l.active.Size()+int64(len(record)) > l.maxSize {
		l.active = NewSegment(l.active.nextIndex(), l.basePath)
	}

	offset, err := l.active.append(key, record)
	if err != nil {
		return err
	}
	l.index[key] = offset

	return nil
}

func (l *Log) Get(key string) []byte {
	return nil
}

func (l *Log) Delete(key string) error {
	return nil
}

func (l *Log) Close() {
	if l.active != nil {
		l.active.Close()
	}
}
