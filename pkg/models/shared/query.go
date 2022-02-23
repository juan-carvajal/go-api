package shared

type SearchParams struct {
	Pagination
	Search string `json:"search"`
}
