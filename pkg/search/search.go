package search

func IsExistString(what string, where *[]string) (exist bool) {
	for _, v := range *where {
		if v == what {
			return true
		}
	}
	return false
}

func IsExistInt(what int, where *[]int) (exist bool) {
	for _, v := range *where {
		if v == what {
			return true
		}
	}
	return false
}
