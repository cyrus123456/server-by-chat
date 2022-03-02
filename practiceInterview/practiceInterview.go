package practiceinterview

import (
	"fmt"
	"sort"
)

func Test() {
	numSlice := []int{-1, 0, 1, 2, 1, -4, -2, 3, -3, 2}
	fmt.Println(threeSumClosest(numSlice, 0))
	// fmt.Println(threeSumClosest2(numSlice, 0))
	// fmt.Println(threeSum(numSlice))
}

func threeSumClosest(numSlice []int, targetSum int) (threeNumSlice [][3]int) {
	sort.Ints(numSlice)
	for i := 0; i < len(numSlice); i++ {
		var currentSum int
		for lIndex, rIndex := i+1, len(numSlice)-1; rIndex > lIndex; {
			fmt.Println(lIndex, rIndex)
			currentSum = numSlice[i] + numSlice[lIndex] + numSlice[rIndex]
			if currentSum < targetSum {
				lIndex++
			} else if currentSum == targetSum {
				threeNumSlice = append(threeNumSlice, [3]int{numSlice[i], numSlice[lIndex], numSlice[rIndex]})
			} else {
				rIndex--
			}
		}
	}
	return
}

func threeSumClosest2(numSlice []int, targetSum int) (threeNumSlice [][3]int) {
	sort.Ints(numSlice)
	for i := 0; i < len(numSlice); i++ {
		for j := 0; j < len(numSlice); j++ {
			for k := 0; k < len(numSlice); k++ {
				if numSlice[i]+numSlice[j]+numSlice[k] == targetSum {
					threeNumSlice = append(threeNumSlice, [3]int{numSlice[i], numSlice[j], numSlice[k]})
				}
			}
		}
	}
	return
}

func threeSum(nums []int) [][]int {
	var res [][]int
	if len(nums) < 3 {
		return res
	}

	sort.Ints(nums)
	for i := 0; i < len(nums)-2; i++ {
		if nums[i] > 0 {
			break
		}

		if i > 0 && nums[i-1] == nums[i] {
			continue
		}

		for j, k := i+1, len(nums)-1; j < k; {
			sum := nums[i] + nums[j] + nums[k]
			if sum == 0 {
				tmp := []int{nums[i], nums[j], nums[k]}
				res = append(res, tmp)
				for j+1 < k && nums[j] == nums[j+1] {
					j++
				}
				for k-1 > j && nums[k] == nums[k-1] {
					k--
				}
				j++
				k--
			} else if sum > 0 {
				k--
			} else {
				j++
			}
		}
	}

	return res
}
