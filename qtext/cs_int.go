package qtext

func csInt(cs bool) C.int {
	if cs {
		return 0
	}
	return 1
}
