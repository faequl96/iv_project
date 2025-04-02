package utils

import "math"

func CalculateDiscountedPrice(price uint, discountPercentage uint) uint {
	discount := (float64(discountPercentage) / 100) * float64(price)
	finalPrice := math.Round(float64(price) - discount)
	return uint(finalPrice)
}
