package practiceinterview

import (
	"fmt"
	"sort"
)

func Test() {
	numSlice := []int{-1, 0, 1, 2, -1, -4, -2}
	fmt.Println(threeSumClosest(numSlice, 0))
}

func threeSumClosest(numSlice []int, targetSum int) [][3]int {
	threeNumSlice := make([][3]int, 1)
	sort.Ints(numSlice)
	for i := 0; i < len(numSlice); i++ {
		lIndex := i + 1
		rIndex := len(numSlice) - 1
		currentSum := numSlice[i] + numSlice[lIndex] + numSlice[rIndex]
		for lIndex < rIndex {
			if currentSum < targetSum {
				lIndex++
			} else if currentSum == targetSum {
				return append(threeNumSlice, [3]int{numSlice[i], numSlice[lIndex], numSlice[rIndex]})
			} else {
				rIndex--
			}
		}
	}

	return threeNumSlice
}
