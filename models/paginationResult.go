package models

type PaginationResult[T any] struct {
	Total   uint16
	Count   uint16
	Page    uint16
	Pages   uint16
	Results []T
}
