package domain

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

type Paginated[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

func (p Paginated[T]) Items() any {
	return p.Data
}

func (p Paginated[T]) PageInfo() Pagination {
	return p.Pagination
}
