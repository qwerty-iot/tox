package tox

func ToPtr[T any](a T) *T {
	p := new(T)
	*p = a
	return p
}

func ItemInArray[T comparable](s T, arr []T, comp func(a T, b T) bool) bool {
	if arr == nil {
		return false
	}
	for _, v := range arr {
		if comp != nil {
			if comp(v, s) {
				return true
			}
		} else {
			if v == s {
				return true
			}
		}
	}
	return false
}
