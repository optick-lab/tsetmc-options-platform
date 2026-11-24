package dto

// UnderlyingItemDto represents a single underlying asset item in the response.
type UnderlyingItemDto struct {
	ID        int64   `json:"id"`
	Symbol    string  `json:"symbol"`
	InsCode   *string `json:"ins_code"`
	AssetType string  `json:"asset_type"`
}

// PaginatedUnderlyingsDto represents the paginated result of underlying assets.
type PaginatedUnderlyingsDto struct {
	Items    []UnderlyingItemDto `json:"items"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
	Total    int64               `json:"total"`
}
