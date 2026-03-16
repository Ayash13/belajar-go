package mathutil

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}

// Subtract returns the difference of two integers.
func Subtract(a, b int) int {
	return a - b
}

// Average returns the average of a slice of float64.
func Average(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	var total float64
	for _, n := range nums {
		total += n
	}
	return total / float64(len(nums))
}
