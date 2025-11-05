package tox

func ToPtr[T any](a T) *T {
	p := new(T)
	*p = a
	return p
}
