package services

const ListPageLimit = 100

func TotalPages(total int64) int {
	return int((total + int64(ListPageLimit) - 1) / int64(ListPageLimit))
}
