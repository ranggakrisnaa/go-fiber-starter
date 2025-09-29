package query

import (

	"github.com/Caknoooo/go-pagination"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"strconv"
)

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	TelpNumber string `json:"telp_number"`
	Role       string `json:"role"`
	ImageUrl   string `json:"image_url"`
	IsVerified bool   `json:"is_verified"`
}

type UserFilter struct {
	pagination.BaseFilter
}

func (f *UserFilter) ApplyFilters(query *gorm.DB) *gorm.DB {
	// Apply your filters here
	return query
}

func (f *UserFilter) GetTableName() string {
	return "users"
}

func (f *UserFilter) GetSearchFields() []string {
	return []string{"name"}
}

func (f *UserFilter) GetDefaultSort() string {
	return "id asc"
}

func (f *UserFilter) GetIncludes() []string {
	return f.Includes
}

func (f *UserFilter) GetPagination() pagination.PaginationRequest {
	return f.Pagination
}

func (f *UserFilter) Validate() {
	var validIncludes []string
	allowedIncludes := f.GetAllowedIncludes()
	for _, include := range f.Includes {
		if allowedIncludes[include] {
			validIncludes = append(validIncludes, include)
		}
	}
	f.Includes = validIncludes
}

func (f *UserFilter) GetAllowedIncludes() map[string]bool {
	return map[string]bool{}
}

// BindPaginationForFiber binds pagination parameters from Fiber context
func (f *UserFilter) BindPaginationForFiber(ctx fiber.Ctx) error {
	// Set default values
	if f.Pagination.Page == 0 {
		f.Pagination.Page = 1
	}
	if f.Pagination.PerPage == 0 {
		f.Pagination.PerPage = 10
	}

	// Get query parameters using Fiber v3 API
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	perPage, _ := strconv.Atoi(ctx.Query("per_page", "10"))

	// Bind values to filter - only update pagination fields
	f.Pagination.Page = page
	f.Pagination.PerPage = perPage

	return nil
}
