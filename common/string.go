package common

func StringIsAllNumber(str string) bool {
	anyReplaceList := []int32{'.', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	for strIndex, strInt := range str {
		if strIndex == 0 {
			if strInt == '-' {
				continue
			}
		}
		if isHaveReplaceList(strInt, anyReplaceList) {
			continue
		}
		return false
	}
	return true
}

func isHaveReplaceList(strInt int32, replaceList []int32) bool {
	for _, chStrInt := range replaceList {
		if chStrInt == strInt {
			return true
		}
	}
	return false
}
