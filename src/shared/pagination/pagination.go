package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Params struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Sort     string `form:"sort"`
	Order    string `form:"order"`
}

type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
	DefaultOrder    = ""
	DefaultSort     = ""
)

func ParseParams(c *gin.Context) Params {
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(DefaultPage)))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(DefaultPageSize)))
	sort := c.DefaultQuery("sort", DefaultSort)
	order := c.DefaultQuery("order", DefaultOrder)

	return Params{
		Page:     validatePage(page),
		PageSize: validatePageSize(pageSize),
		Sort:     sort,
		Order:    validateOrder(order),
	}
}

func validatePage(page int) int {
	if page < 1 {
		return DefaultPage
	}
	return page
}

func validatePageSize(pageSize int) int {
	if pageSize < 1 {
		return DefaultPageSize
	}
	if pageSize > MaxPageSize {
		return MaxPageSize
	}
	return pageSize
}

func validateOrder(order string) string {
	if order != "asc" && order != "desc" {
		return "desc"
	}
	return order
}

func (p *Params) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *Params) GetOrderClause() string {
	if p.Sort == "" {
		return ""
	}
	return p.Sort + " " + p.Order
}

func (p *Params) BuildMeta(totalItems int64) Meta {
	totalPages := int(math.Ceil(float64(totalItems) / float64(p.PageSize)))

	return Meta{
		Page:       p.Page,
		PageSize:   p.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasNext:    p.Page < totalPages,
		HasPrev:    p.Page > 1,
	}
}

func (p *Params) Apply(db *gorm.DB, dest interface{}) (interface{}, Meta, error) {
	var totalItems int64

	if err := db.Model(dest).Count(&totalItems).Error; err != nil {
		return nil, Meta{}, err
	}

	if err := db.Order(p.GetOrderClause()).Limit(p.PageSize).Offset(p.GetOffset()).Find(dest).Error; err != nil {
		return nil, Meta{}, err
	}

	meta := p.BuildMeta(totalItems)

	return dest, meta, nil
}

func (p *Params) ApplyWithQuery(db *gorm.DB, dest interface{}, query interface{}, args ...interface{}) (interface{}, Meta, error) {
	var totalItems int64

	baseQuery := db.Model(dest)
	if query != nil {
		baseQuery = baseQuery.Where(query, args...)
	}

	if err := baseQuery.Count(&totalItems).Error; err != nil {
		return nil, Meta{}, err
	}

	if err := baseQuery.Order(p.GetOrderClause()).Limit(p.PageSize).Offset(p.GetOffset()).Find(dest).Error; err != nil {
		return nil, Meta{}, err
	}

	meta := p.BuildMeta(totalItems)

	return dest, meta, nil
}
