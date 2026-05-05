package internal

import "iter"

type Log struct {
	maxSize int64
	index   map[string]int64
	active  *Segment
}

func (l *Log) next() iter.Seq2[string, int64] {
	return func(yield func(K string, V int64) bool) {
		// We don't know how many files exist
		// We know that loop will be break by an error
		startSegmentIndex := 0
		var startOffset int64
		segment := &Segment{}

		for {
			if err := segment.open(startSegmentIndex); err != nil {
				return
			}
			for {
				record, nextOffset, err := segment.ReadAt(startOffset)
				if err != nil {
					break
				}
				if !yield(string(record.Key), startOffset) {
					segment.close()
					return
				}
				startOffset = nextOffset
			}
			startSegmentIndex += 1
			startOffset = 0
			segment.close()
		}
	}
}

func (l *Log) buildIndex() {
	for k, v := range l.next() {
		l.index[k] = v
	}
}

func (l *Log) Open(path string) {

}
