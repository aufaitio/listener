package apis

import (
	"strconv"

	"github.com/aufaitio/listener/util"
	"github.com/go-ozzo/ozzo-routing"
)

const (
	defaultPageSize int = 100
	maxPageSize     int = 1000
)

func getPaginatedListFromRequest(c *routing.Context, count int64) *util.PaginatedList {
	page := parseInt(c.Query("page"), 1)
	perPage := parseInt(c.Query("perPage"), defaultPageSize)
	if perPage <= 0 {
		perPage = defaultPageSize
	}
	if perPage > maxPageSize {
		perPage = maxPageSize
	}
	return util.NewPaginatedList(page, perPage, int(count))
}

func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return defaultValue
}
