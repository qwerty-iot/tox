package tox

func ArrayToMapBool[T comparable](arr []T) map[T]bool {
	ret := map[T]bool{}
	for _, i := range arr {
		ret[i] = true
	}
	return ret
}

func ArrayContains[T comparable](s T, arr []T, comp func(a T, b T) bool) bool {
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
