package service

import "math"

func NormalizePage(filter PageFilter) PageFilter {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	return filter
}

func TotalPages(total int, limit int) int {
	if limit < 1 {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(limit)))
}
