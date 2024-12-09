package txn

import "lsm/entries"

type InclusiveRange struct {
	start entries.Key
	end   entries.Key
}

func NewInclusiveRange(start, end entries.Key) InclusiveRange {
	return InclusiveRange{start: start, end: end}
}

func (r InclusiveRange) Start() entries.Key {
	return r.start
}

func (r InclusiveRange) End() entries.Key {
	return r.end
}
