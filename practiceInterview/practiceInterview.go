package practiceinterview

func Test() {

	// 堆排序
	heapSort(tempSlice)
	// 归并排序
	// start_1 := time.Now() // 获取当前时间
	// segmentation(tempSlice)
	// fmt.Println("未优化/n", time.Since(start_1))

	// start_2 := time.Now() // 获取当前时间
	// pointerSplitMergeSort(tempSlice, 0, len(tempSlice))
	// fmt.Println("优化后/n", time.Since(start_2))
}

func heapSort(slice_arg []int) {

}

func adjustHeap(slice_arg []int, index_arg, sliceLength_arg int) {

}

// 优化方法1---------------------------------------------------------------------
func pointerSplitMergeSort(slice_arg []int, startIndex_arg, endInde_arg int) {

	if endInde_arg-startIndex_arg < 2 {
		return
	}

	mideleIndex := (startIndex_arg + endInde_arg) >> 1

	pointerSplitMergeSort(slice_arg, startIndex_arg, mideleIndex)
	pointerSplitMergeSort(slice_arg, mideleIndex, endInde_arg)

	pointerSplitMerge(slice_arg, startIndex_arg, mideleIndex, endInde_arg)

}

func pointerSplitMerge(slice_arg []int, startIndex_arg, mideleIndex_arg, endInde_arg int) {

	leftSliceLen := mideleIndex_arg - startIndex_arg
	leftSlice := make([]int, leftSliceLen)

	copy(leftSlice, slice_arg[startIndex_arg:mideleIndex_arg])

	leftIndex, startIndex, rightIndex := 0, startIndex_arg, mideleIndex_arg
	for leftIndex < leftSliceLen {
		if rightIndex < endInde_arg && slice_arg[rightIndex] < leftSlice[leftIndex] {
			slice_arg[startIndex] = slice_arg[rightIndex]
			startIndex++
			rightIndex++
		} else {
			slice_arg[startIndex] = leftSlice[leftIndex]
			startIndex++
			leftIndex++
		}
	}

}

// 普通方法2-------------------------------------------------------------------------
// 分割切片
func segmentation(_tempSlice []int) []int {
	if len(_tempSlice) < 2 {
		return _tempSlice
	}
	// 找到中间下标
	middleIndeex := len(_tempSlice) / 2
	// 左半部分切片只剩两个
	leftSlice := segmentation(_tempSlice[:middleIndeex])
	// 右半部分切片只剩两个
	rightSlice := segmentation(_tempSlice[middleIndeex:])
	return mergeSlices(leftSlice, rightSlice)
}

// 合并切片
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

	resultSlice = append(resultSlice, _rightSlice[rightIndex:]...)
	resultSlice = append(resultSlice, _leftSlice[leftIndex:]...)

	return resultSlice
}
