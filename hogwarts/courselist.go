// +build !solution

package hogwarts

var used = make(map[string]bool)
var color = make(map[string]int)
var edge = make(map[string][]string)
var prev = make(map[string]string)
var topSort = make([]string, 0)

func dfsFindCycle(subject string) bool {
	color[subject] = 1
	for _, value := range edge[subject] {
		if color[value] == 0 {
			prev[value] = subject
			if dfsFindCycle(value) {
				return true
			}
		} else if color[value] == 1 {
			return true
		}
	}
	color[subject] = 2
	return false
}

func dfsTopSort(subject string) {
	if used[subject] {
		return
	}
	used[subject] = true
	for _, value := range edge[subject] {
		dfsTopSort(value)
	}
	topSort = append(topSort, subject)
}

func GetCourseList(prereqs map[string][]string) []string {
	var color = make(map[string]int)
	for key, value := range prereqs {
		for _, subject := range value {
			edge[subject] = append(edge[subject], key)
		}
	}
	for _, value := range edge {
		for _, subject := range value {
			if color[subject] == 0 {
				if dfsFindCycle(subject) {
					panic("Cycle!")
				}
			}
		}
	}
	for key := range edge {
		if !used[key] {
			dfsTopSort(key)
		}
	}
	for i, j := 0, len(topSort)-1; i < j; i, j = i+1, j-1 {
		topSort[i], topSort[j] = topSort[j], topSort[i]
	}
	return topSort
}
