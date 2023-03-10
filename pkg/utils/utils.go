package utils

// GetSortNum 两个数字排序，小的数字在前面，大的数字在后面
func GetSortNum(a, b uint64) (uint64, uint64) {
	if a > b {
		return b, a
	}
	return a, b
}

// GetLargerNum 返回更大的一个数字
func GetLargerNum(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}
