package po

type SearchList[T any] struct {
	Total int64
	Data  []*T
}

type PageQuery[T any] struct {
	PageNum   int32
	PageSize  int32
	Condition *T
	Sort      map[string]string
}
