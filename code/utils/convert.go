package utils

import "github.com/axgle/mahonia"


//转换为utf8
func ToUTF8(str string) string {
	if str == ""{
		return str
	}
	srcCoder := mahonia.NewDecoder("gbk")
	srcResult := srcCoder.ConvertString(str)
	tagCoder := mahonia.NewDecoder("utf-8")
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

//检查内容是不是utf8
func IsUTF8(str string) bool {
	buf := []byte(str)
	nBytes := 0
	for i := 0; i < len(buf); i++ {
		if nBytes == 0 {
			if (buf[i] & 0x80) != 0 { //与操作之后不为0，说明首位为1
				for (buf[i] & 0x80) != 0 {
					buf[i] <<= 1 //左移一位
					nBytes++     //记录字符共占几个字节
				}
				if nBytes < 2 || nBytes > 6 { //因为UTF8编码单字符最多不超过6个字节
					return false
				}
				nBytes-- //减掉首字节的一个计数
			}
		} else {                     //处理多字节字符
			if buf[i]&0xc0 != 0x80 { //判断多字节后面的字节是否是10开头
				return false
			}
			nBytes--
		}
	}
	return nBytes == 0
}

//公共转换字符转换
//ConvertCoder(str, "gbk", "utf-8")
func ConvertCoder(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}
