// +build !solution

package hotelbusiness

import (
	"math"
)

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func ComputeLoad(guests []Guest) []Load {
	var guestsMap = make(map[int]int)
	var maxDay = -1
	var minDay = math.MaxInt32
	for _, value := range guests {
		guestsMap[value.CheckInDate]++
		guestsMap[value.CheckOutDate]--
		maxDay = Max(maxDay, value.CheckInDate)
		maxDay = Max(maxDay, value.CheckOutDate)
		minDay = Min(minDay, value.CheckInDate)
		minDay = Min(minDay, value.CheckOutDate)
	}
	var guestsAmount []int
	guestsCount := 0
	for i := minDay; i <= maxDay; i++ {
		guestsCount += guestsMap[i]
		guestsAmount = append(guestsAmount, guestsCount)
	}
	var answer = make([]Load, 0)
	for i, value := range guestsAmount {
		if i > 0 {
			if value == answer[len(answer)-1].GuestCount {
				continue
			}
		}
		answer = append(answer, Load{i + minDay, value})
	}
	return answer
}
