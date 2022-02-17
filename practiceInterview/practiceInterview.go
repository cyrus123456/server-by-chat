package practiceinterview

import (
	"fmt"
)

func Test() {
	str := "dsgj"
	reverStr, Ok := reverString(str)
	if Ok {
		fmt.Println("翻转成功", reverStr)
		return
	}
	fmt.Println("暗转失败")
}

func reverString(str string) (string, bool) {
	Str := []rune(str)
	lengh := len(Str)
	if lengh > 5000 {
		fmt.Println("字符串长度超过5000")
		return str, false
	}
	for i := 0; i < lengh; i++ {
		Str[i], Str[lengh-1-i] = Str[lengh-1-i], Str[i]
	}
	return string(Str), true
}
