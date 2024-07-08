package models

import (
	"encoding/json"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const DefaultPageSize = 20
const DefaultPage = 1

const MaxPageSize = 100
const MaxPage = 1000

type Page struct {
	CurrentPage int
	PageSize    int
	Sort        map[string]string
	Filter      map[string]interface{}
}

func NewPage(current, size int) *Page {
	return &Page{
		CurrentPage: current,
		PageSize:    size,
		Sort:        map[string]string{},
		Filter:      map[string]interface{}{},
	}
}

func NewPageFromContext(c *TrackDocsContext) *Page {
	currentPage := c.DefaultQuery("current", strconv.Itoa(DefaultPage))
	pageSize := c.DefaultQuery("pageSize", strconv.Itoa(DefaultPageSize))
	sort := c.DefaultQuery("sort", "{}")
	filter := c.DefaultQuery("filter", "{}")

	currentPageInt, err := strconv.Atoi(currentPage)
	if err != nil || currentPageInt > MaxPage || currentPageInt < 1 {
		currentPageInt = DefaultPage
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt > MaxPageSize || pageSizeInt < 1 {
		pageSizeInt = DefaultPageSize
	}

	var sortJson map[string]string
	err = json.Unmarshal([]byte(sort), &sortJson)
	if err != nil {
		sortJson = map[string]string{}
	}
	var filterJson map[string]interface{}
	err = json.Unmarshal([]byte(filter), &filterJson)
	if err != nil {
		filterJson = map[string]interface{}{}
	}

	return &Page{
		CurrentPage: currentPageInt,
		PageSize:    pageSizeInt,
		Sort:        sortJson,
		Filter:      filterJson,
	}
}

func (p *Page) Paginate(db *gorm.DB) *gorm.DB {
	offset := (p.CurrentPage - 1) * p.PageSize
	if len(p.Sort) > 0 {
		done := false
		for k, v := range p.Sort {
			if done {
				break
			}
			db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: k}, Desc: v == "descend"})
			done = true
		}
	}
	if len(p.Filter) > 0 {
		db = db.Where(p.Filter)
	}
	return db.Offset(offset).Limit(p.PageSize)
}

func (p *Page) CountPaginate(db *gorm.DB) *gorm.DB {
	if len(p.Filter) > 0 {
		db = db.Where(p.Filter)
	}
	return db
}
