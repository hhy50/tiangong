package common

func Max[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](nums ...T) T {
	if len(nums) == 1 {
		return nums[0]
	}
	max := nums[0]

	for i := 1; i < len(nums); i++ {
		if max < nums[i] {
			max = nums[i]
		}
	}
	return max
}
