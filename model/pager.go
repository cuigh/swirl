package model

import (
	"math"
	"net/url"
	"strconv"
)

// PageSize is default page size.
const PageSize = 25

// Pager holds pagination info.
type Pager struct {
	Count    int
	Size     int
	Page     int
	Pages    []int
	Previous int
	Next     int
	MaxPage  int
	URL      func(page int) string
}

// NewPager returns a Pager instance which displays max 7 pages.
func NewPager(uri string, count, size, page int) *Pager {
	if page < 1 {
		page = 1
	}

	u, _ := url.Parse(uri)
	q := u.Query()

	p := &Pager{
		Count:    count,
		Size:     size,
		Page:     page,
		Previous: page - 1,
		MaxPage:  int(math.Ceil(float64(count) / float64(size))),
		URL: func(page int) string {
			if page < 2 {
				q.Del("page")
			} else {
				q.Set("page", strconv.Itoa(page))
			}
			u.RawQuery = q.Encode()
			return u.String()
		},
	}
	if p.Page < p.MaxPage {
		p.Next = page + 1
	}

	// calculate pages
	if count == 0 {
		p.Pages = []int{1}
	} else if p.MaxPage <= 7 {
		p.Pages = make([]int, p.MaxPage)
		for i := 0; i < p.MaxPage; i++ {
			p.Pages[i] = i + 1
		}
	} else if page <= 4 {
		p.Pages = []int{1, 2, 3, 4, 5, 6, 0, p.MaxPage}
	} else if page > (p.MaxPage - 4) {
		p.Pages = []int{1, 0}
		for i := p.MaxPage - 5; i <= p.MaxPage; i++ {
			p.Pages = append(p.Pages, i)
		}
	} else {
		p.Pages = []int{1, 0, page - 2, page - 1, page, page + 1, page + 2, 0, p.MaxPage}
	}

	return p
}
