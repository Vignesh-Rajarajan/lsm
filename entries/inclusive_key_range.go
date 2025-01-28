package entries

type LessOrEqual interface {
	IsLessThanOrEqualTo(other LessOrEqual) bool
}

type InclusiveKeyRange[T LessOrEqual] struct {
	start T
	end   T
}

func NewInclusiveKeyRange[T LessOrEqual](start, end T) InclusiveKeyRange[T] {
	if !start.IsLessThanOrEqualTo(end) {
		panic("start key must be less than or equal to end key")
	}
	return InclusiveKeyRange[T]{start: start, end: end}
}

func (r InclusiveKeyRange[T]) Start() T {
	return r.start
}

func (r InclusiveKeyRange[T]) End() T {
	return r.end
}
