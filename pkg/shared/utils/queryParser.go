package utils

import (
	"net/http"
	"strconv"

	"github.com/juan-carvajal/go-api/pkg/models/shared"
)

const DEFAUL_PAGE_SIZE = 100

func ParseCommonQueryParams(r *http.Request) shared.SearchParams {
	queryParams := r.URL.Query()

	params := shared.SearchParams{}

	rawOffset := queryParams.Get("offset")

	offset, err := strconv.Atoi(rawOffset)

	if err != nil {
		params.Offset = 0
	} else {
		params.Offset = offset
	}

	rawPageSize := queryParams.Get("page_size")

	pageSize, err := strconv.Atoi(rawPageSize)

	if err != nil {
		params.PageSize = DEFAUL_PAGE_SIZE
	} else {
		params.Offset = pageSize
	}

	search := queryParams.Get("search")

	params.Search = search

	return params
}
