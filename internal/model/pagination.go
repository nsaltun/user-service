package model

// Pagination contains metadata for paginated responses
type Pagination struct {
	TotalRecords int64       `json:"totalRecords"` // Total number of records in the collection
	Limit        int         `json:"limit"`        // Number of items per page
	Offset       int         `json:"offset"`       // Number of items to skip
	HasNext      bool        `json:"hasNext"`      // Indicator if there is a next page
	HasPrevious  bool        `json:"hasPrevious"`  // Indicator if there is a previous page
	Items        interface{} `json:"items"`        // Actual paginated data
}
