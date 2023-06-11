package utils

import (
	"bytes"
)

/**
切片合并为字符串
*/
func SlicetoString(glue string, pieces []string) string {
	var buf bytes.Buffer
	l := len(pieces)
	for _, str := range pieces {
		buf.WriteString(str)
		if l--; l > 0 {
			buf.WriteString(glue)
		}
	}
	return buf.String()
}

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

/*
取交集
*/
func SliceIntersect(list1 []string, list2 []string) []string {
	b2 := []string{}
	for _, v1 := range list1 {
		for _, v2 := range list2 {
			if v1 == v2 {
				b2 = append(b2, v1)
			}
		}
	}
	return b2
}

/*
删除最后一个元素
*/
func SliceDeleteLast(s []string) []string {
	if len(s) == 0 {
		return s
	}
	s = append(s[:len(s)-1])
	return s
}

// 去重
func SliceRemoveDuplicate(slc []string) []string {
	result := []string{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

// 去除空数据
func SliceRemoveEmpty(s []string) []string {
	if len(s) == 0 {
		return s
	}
	j := 0
	for _, v := range s {
		if v != "" {
			s[j] = v
			j++
		}
	}
	return s[:j]
}
