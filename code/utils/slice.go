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
