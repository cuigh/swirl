package model

import (
	"github.com/cuigh/auxo/data"
)

// Chart represents a dashboard chart.
type Chart struct {
	Name        string   `json:"name" bson:"_id" valid:"required"` // unique, the name of build-in charts has '$' prefix.
	Title       string   `json:"title" valid:"required"`
	Description string   `json:"desc"`
	Label       string   `json:"label"` // ${name} - ${instance}
	Query       string   `json:"query" valid:"required"`
	Kind        string   `json:"kind"`      // builtin/custom
	Dashboard   string   `json:"dashboard"` // home/service/task...
	Type        string   `json:"type"`      // pie/line...
	Unit        string   `json:"unit"`      // bytes/milliseconds/percent:100...
	Width       int32    `json:"width"`     // 1-12(12 columns total)
	Height      int32    `json:"height"`    // default 50
	Colors      []string `json:"colors"`
	Options     data.Map `json:"options"`
}

func NewChart(name, title, label, query, unit string) *Chart {
	return &Chart{
		Name:        name,
		Title:       title,
		Description: title,
		Label:       label,
		Query:       query,
		Type:        "line",
		Unit:        unit,
		Width:       12,
		Height:      50,
	}
}

type ChartItem struct {
	Name   string   `json:"name"`
	Width  int32    `json:"width"`
	Height int32    `json:"height"`
	Colors []string `json:"colors"`
}

type ChartPanel struct {
	Refresh bool        `json:"refresh"`
	Period  int32       `json:"period"` // minutes
	Charts  []ChartItem `json:"charts"`
}

//type ChartPanel []ChartItem

type ChartPoint struct {
	X int64   `json:"x"`
	Y float64 `json:"y"`
}

type ChartLine struct {
	Label string       `json:"label"`
	Data  []ChartPoint `json:"data"`
}

type ChartValue struct {
	Label string  `json:"label"`
	Data  float64 `json:"data"`
}

type ChartVector struct {
	Data   []float64 `json:"data"`
	Labels []string  `json:"labels"`
}

type ChartInfo struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Label string `json:"label"`
	Query string `json:"query"`
}

func NewChartInfo(name, title, label, query string) ChartInfo {
	return ChartInfo{
		Name:  name,
		Title: title,
		Label: label,
		Query: query,
	}
}
