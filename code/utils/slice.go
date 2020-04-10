package utils

/*
查看slice中一个字段是否存在
字符串
*/

func InSliceString(str string, slice []string) bool {
	isset := false
	for _, v := range slice {
		if v == str {
			isset = true
			break
		}
	}
	return isset
}

/*
取差集
*/
func SliceDiff(slice1 []string, slice2 []string) []string {
	newSlice := []string{}
	isset := false
	for _, v := range slice1 {
		isset = false
		for _, v2 := range slice2 {
			if v2 == v {
				isset = true
				break
			}
		}
		if isset == false {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}
