// +build !solution

package reverse

import (
	"unicode/utf8"
)

func Reverse(input string) string {
	length := len(input)
	buffer := make([]byte, length)
	for i := 0; i < length; {
		r, size := utf8.DecodeRuneInString(input[i:])
		i += size
		utf8.EncodeRune(buffer[length-i:], r)
	}
	return string(buffer)
}
