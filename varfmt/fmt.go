// +build !solution

package varfmt

import (
	"fmt"
	"strconv"
)

func Sprintf(format string, args ...interface{}) string {
	resultString := make([]byte, 0)
	bracketIndex := 0
	convertedMap := make(map[interface{}]string)
	for i := 0; i < len(format); i++ {
		if format[i] == '{' {
			j := i + 1
			for ; format[j] != '}'; j++ {
			}
			index := 0
			if i+1 == j {
				index = bracketIndex
			} else {
				index, _ = strconv.Atoi(format[i+1 : j])
			}
			if convertedMap[args[index]] != "" {
				resultString = append(resultString, []byte(convertedMap[args[index]])...)
			} else {
				convResult := fmt.Sprint(args[index])
				resultString = append(resultString, []byte(convResult)...)
				convertedMap[args[index]] = convResult
			}
			bracketIndex++
			i = j
		} else {
			resultString = append(resultString, format[i])
		}
	}
	return string(resultString)
}
