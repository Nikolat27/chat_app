package utils

import (
	"net/url"
	"strconv"
)

func ParsePageAndLimitQueryParams(url *url.URL) (int64, int64, *ErrorResponse) {
	pageStr := url.Query().Get("page")
	pageLimit := url.Query().Get("limit")

	if pageStr == "" || pageLimit == "" {
		return 1, 6, nil
	}

	pageInt, err := strconv.ParseInt(pageStr, 0, 64)
	if err != nil {
		return 0, 0, &ErrorResponse{Type: "parseInt", Detail: err.Error()}
	}

	limitInt, err := strconv.ParseInt(pageLimit, 0, 64)
	if err != nil {
		return 0, 0, &ErrorResponse{Type: "parseInt", Detail: err.Error()}
	}

	return pageInt, limitInt, nil
}
