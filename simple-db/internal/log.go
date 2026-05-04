package internal

import "iter"

type Log struct {
	maxSize int64
	index   map[string]int
	active  *Segment
}

func (l *Log) next() iter.Seq2[string, int] {
	startSegmentIndex := 0
	segment := &Segment{}

	return func(yield func(K string, V int) bool) {
		segment.open(startSegmentIndex)
		yield("", 0)
		segment.close()
	}
}

func (l *Log) buildIndex() {
	for k, v := range l.next() {
		l.index[k] = v
	}
}

func (l *Log) Open(path string) {

}
