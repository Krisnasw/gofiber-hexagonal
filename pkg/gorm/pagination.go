package gorm

import (
	"gorm.io/gorm"
)

// Pagination represents pagination parameters
type Pagination struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Cursor   string `json:"cursor"` // For cursor-based pagination
	OrderBy  string `json:"order_by"`
	Sort     string `json:"sort"` // asc or desc
}

// PaginationResult represents pagination results
type PaginationResult struct {
	CurrentPage    int         `json:"current_page"`
	PageSize       int         `json:"page_size"`
	TotalRecords   int64       `json:"total_records"`
	TotalPages     int         `json:"total_pages"`
	HasNext        bool        `json:"has_next"`
	HasPrevious    bool        `json:"has_previous"`
	NextCursor     string      `json:"next_cursor,omitempty"`
	PreviousCursor string      `json:"previous_cursor,omitempty"`
	Data           interface{} `json:"data"`
}

// OffsetPagination performs offset-based pagination
func OffsetPagination(db *gorm.DB, pagination *Pagination, dest interface{}) (*PaginationResult, error) {
	var totalRecords int64
	result := &PaginationResult{}

	// Get total count
	countDB := db.Session(&gorm.Session{})
	if err := countDB.Count(&totalRecords).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	page := pagination.Page
	if page <= 0 {
		page = 1
	}

	pageSize := pagination.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Apply ordering
	orderQuery := pagination.OrderBy
	if orderQuery == "" {
		orderQuery = "id"
	}

	sortOrder := pagination.Sort
	if sortOrder == "" {
		sortOrder = "asc"
	}

	query := db.Offset(offset).Limit(pageSize).Order(orderQuery + " " + sortOrder)

	if err := query.Find(dest).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalRecords + int64(pageSize) - 1) / int64(pageSize))

	result.CurrentPage = page
	result.PageSize = pageSize
	result.TotalRecords = totalRecords
	result.TotalPages = totalPages
	result.HasNext = page < totalPages
	result.HasPrevious = page > 1
	result.Data = dest

	return result, nil
}

// CursorPagination performs cursor-based pagination
func CursorPagination(db *gorm.DB, pagination *Pagination, cursorField string, dest interface{}) (*PaginationResult, error) {
	result := &PaginationResult{}

	pageSize := pagination.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// Apply cursor condition if cursor is provided
	if pagination.Cursor != "" {
		// Note: In a real implementation, you would decode the cursor
		// and apply the appropriate WHERE condition
		// This is a simplified example
		db = db.Where(cursorField+" > ?", pagination.Cursor)
	}

	// Apply ordering
	orderQuery := pagination.OrderBy
	if orderQuery == "" {
		orderQuery = cursorField
	}

	sortOrder := pagination.Sort
	if sortOrder == "" {
		sortOrder = "asc"
	}

	query := db.Limit(pageSize).Order(orderQuery + " " + sortOrder)

	if err := query.Find(dest).Error; err != nil {
		return nil, err
	}

	// Calculate next cursor (simplified)
	// In a real implementation, you would calculate based on the last record
	result.NextCursor = ""
	result.PreviousCursor = pagination.Cursor

	result.PageSize = pageSize
	result.Data = dest

	return result, nil
}

// NewPagination creates a new pagination instance with default values
func NewPagination() *Pagination {
	return &Pagination{
		Page:     1,
		PageSize: 10,
		Sort:     "asc",
	}
}
