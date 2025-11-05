package tox

func ArrayToMapBool[T comparable](arr []T) map[T]bool {
	ret := map[T]bool{}
	for _, i := range arr {
		ret[i] = true
	}
	return ret
}
