package main

import (
	"math/rand"
)

type Point struct {
	X int
	Y int
}

func reverseRight(numbers []int32) []int32 {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}

func reverseLeft(nums [][]int32) [][]int32 {
	for x := 0; x < len(nums); x++ {
		l := len(nums)
		for x := 0; x < l/2; x++ {
			nums[x], nums[l-x-1] = nums[l-x-1], nums[x]
		}
	}
	return nums
}

func generateBlock() [][]int32 {
	switch i := rand.Int31n(7); i {
	// = block
	case 0:
		return [][]int32{
			[]int32{1, 1},
			[]int32{1, 1},
		}
		// T block
	case 1:
		return [][]int32{
			[]int32{0, 0, 0},
			[]int32{2, 2, 2},
			[]int32{0, 2, 0},
		}
	//L Block
	case 2:
		return [][]int32{
			[]int32{0, 3, 0},
			[]int32{0, 3, 0},
			[]int32{0, 3, 3},
		}
	//J Block
	case 3:
		return [][]int32{
			[]int32{0, 4, 0},
			[]int32{0, 4, 0},
			[]int32{4, 4, 0},
		}
	//S Block
	case 4:
		return [][]int32{
			[]int32{0, 5, 5},
			[]int32{5, 5, 0},
			[]int32{0, 0, 0},
		}
	//Z Block
	case 5:
		return [][]int32{
			[]int32{6, 6, 0},
			[]int32{0, 6, 6},
			[]int32{0, 0, 0},
		}
	//I Block
	default:
		return [][]int32{
			[]int32{0, 7, 0, 0},
			[]int32{0, 7, 0, 0},
			[]int32{0, 7, 0, 0},
			[]int32{0, 7, 0, 0},
		}
	}
}
