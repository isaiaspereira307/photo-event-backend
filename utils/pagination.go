package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Pagination represents pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}

// GetPaginationFromRequest extracts pagination parameters from request
func GetPaginationFromRequest(ctx *gin.Context) Pagination {
	// Default pagination values
	page := 1
	pageSize := 10

	// Parse page number
	pageStr := ctx.DefaultQuery("page", "1")
	if pageVal, err := strconv.Atoi(pageStr); err == nil && pageVal > 0 {
		page = pageVal
	}

	// Parse page size
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	if pageSizeVal, err := strconv.Atoi(pageSizeStr); err == nil && pageSizeVal > 0 {
		pageSize = pageSizeVal
	}

	return Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalRows:  0,
		TotalPages: 0,
	}
}

// Paginate applies pagination to a GORM query
func Paginate(db *gorm.DB, pagination *Pagination, result interface{}) (*gorm.DB, error) {
	var totalRows int64
	countDb := db.Model(result)
	if err := countDb.Count(&totalRows).Error; err != nil {
		return db, err
	}
	pagination.TotalRows = totalRows

	// Calculate total pages
	pagination.TotalPages = int((totalRows + int64(pagination.PageSize) - 1) / int64(pagination.PageSize))

	// Apply offset and limit to the query
	offset := (pagination.Page - 1) * pagination.PageSize
	return db.Offset(offset).Limit(pagination.PageSize), nil
}
