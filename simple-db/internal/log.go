package internal

import "iter"

type Log struct {
	maxSize int64
	index   map[string]int
	active  *Segment
}

func (l *Log) next() iter.Seq2[string, int64] {
	startSegmentIndex := 0
	var startOffset int64
	segment := &Segment{}

	return func(yield func(K string, V int64) bool) {
		err := segment.open(startSegmentIndex)
		defer segment.close()
		if err != nil {
			return
		}

		record, nextOffset, err := segment.ReadAt(startOffset)
		if err != nil || !yield(string(record.Key), startOffset) {
			return
		}
		startSegmentIndex += 1
		startOffset = nextOffset

	}
}

func (l *Log) buildIndex() {
	for k, v := range l.next() {
		l.index[k] = v
	}
}

func (l *Log) Open(path string) {

}
