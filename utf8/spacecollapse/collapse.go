// +build !solution

package spacecollapse

func CollapseSpaces(input string) string {
	length := len(input)
	wasSpace := false
	result := make([]rune, 0, length)
	for _, value := range input {
		if value == 10 || value == 9 || value == 32 || value == 13 {
			if wasSpace {
				continue
			} else {
				wasSpace = true
				result = append(result, ' ')
				continue
			}
		} else {
			if wasSpace {
				wasSpace = false
			}
		}
		result = append(result, value)
	}
	return string(result)
}
