package common

// PageResult page result
type PageResult struct {
	PageIndex  int
	PageSize   int
	TotalPages int
	TotalCount int
	Data       interface{}
}
