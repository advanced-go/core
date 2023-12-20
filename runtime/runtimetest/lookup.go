package runtimetest

//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces

const (
	stringValueError = "error: stringFromType() value parameter is nil"
	listValueError   = "error: listFromType() value parameter is nil"
)

// LookupResultConstraints - lookup function constraints
type LookupResultConstraints interface {
	string | []string
}

// LookupFunctionConstraints - lookup function constraints
type LookupFunctionConstraints interface {
	func(string) string | func(string) []string
}

type Lookup[T LookupResultConstraints] interface {
	Resolve(key string) (string, bool)
	SetOverride(t any)
}

type lookup[F LookupFunctionConstraints] struct {
	overrideFn F
	defaultFn  F
}

func (l *lookup[F]) SetOverride(t any) {
	if t == nil {
		l.overrideFn = nil
	}
}

func (l *lookup[F]) Resolve(key string) (string, bool) {
	return key, false
}

func NewLookup[T LookupResultConstraints, F LookupFunctionConstraints](defaultFn F) Lookup[T] {
	l := new(lookup[F])
	l.defaultFn = defaultFn
	return l
}
