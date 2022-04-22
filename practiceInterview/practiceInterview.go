package practiceinterview

import (
	"fmt"
)

func segmentation(_tempSlice []int) []int {
	if len(_tempSlice) < 2 {
		return _tempSlice
	}
	middleIndeex := len(_tempSlice) / 2
	leftSlice := segmentation(_tempSlice[:middleIndeex])
	rightSlice := segmentation(_tempSlice[middleIndeex:])
	return mergeSlices(leftSlice, rightSlice)
}

func mergeSlices(_leftSlice, _rightSlice []int) []int {

	resultSlice := make([]int, 0)
	leftIndex, rightIndex := 0, 0
	leftLen, rightLen := len(_leftSlice), len(_rightSlice)

	for leftIndex < leftLen && rightIndex < rightLen {

		if _leftSlice[leftIndex] > _rightSlice[rightIndex] {
			resultSlice = append(resultSlice, _rightSlice[rightIndex])
			rightIndex++

		} else {
			resultSlice = append(resultSlice, _leftSlice[leftIndex])
			leftIndex++

		}
	}

	return resultSlice
}

func Test() {
	tempSlice := []int{43, 6, 45, 345, 2, 6, 5, 3, 7, 78, 53, 4, 3456, 84, 3, 7, 4, 6, 3}

	segmentation(tempSlice)

	fmt.Println(tempSlice)
}
