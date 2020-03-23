// +build !solution

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args[1:]
	m := make(map[string]int)
	var addString = ""
	for _, fileName := range args {
		b, err := ioutil.ReadFile(fileName)
		checkError(err)
		stringInFile := string(b)
		for _, char := range stringInFile {
			if char == 10 {
				m[addString]++
				addString = ""
			} else {
				addString += string(char)
			}
		}
		m[addString]++
		addString = ""
	}
	for key, value := range m {
		if value >= 2 {
			fmt.Println(strconv.Itoa(value) + "\t" + key)
		}
	}
}
