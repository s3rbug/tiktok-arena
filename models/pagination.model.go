package models

type PaginationQueries struct {
	Page     int    `query:"page"`
	Count    int    `query:"count"`
	SortName string `query:"sort_name"`
	SortSize string `query:"sort_size"`
}

type SearchBody struct {
	SearchText string
}

func ValidatePaginationQueries(queries *PaginationQueries) {
	if queries.Page <= 0 {
		queries.Page = 1
	}

	switch {
	case queries.Count > 50:
		queries.Count = 50
	case queries.Count <= 0:
		queries.Count = 20
	}

	switch queries.SortName {
	case "asc":
		queries.SortName = "asc"
	case "desc":
		queries.SortName = "desc"
	default:
		queries.SortName = ""
	}

	switch queries.SortSize {
	case "asc":
		queries.SortSize = "asc"
	case "desc":
		queries.SortSize = "desc"
	default:
		queries.SortSize = ""
	}
}
