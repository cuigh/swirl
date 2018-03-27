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

func NewChart(dashboard, name, title, label, query, unit string) *Chart {
	return &Chart{
		Name:        name,
		Title:       title,
		Description: title,
		Label:       label,
		Query:       query,
		Dashboard:   dashboard,
		Type:        "line",
		Unit:        unit,
		Width:       12,
		Height:      150,
	}
}

type ChartItem struct {
	Name   string   `json:"name"`
	Width  int32    `json:"width"`
	Height int32    `json:"height"`
	Colors []string `json:"colors"`
}

type ChartDashboard struct {
	Name            string      `json:"name"`
	Key             string      `json:"key"`
	Period          int32       `json:"period"`           // minutes
	RefreshInterval int32       `json:"refresh_interval"` // seconds, 0 means disabled.
	Charts          []ChartItem `json:"charts"`
}

func (cd *ChartDashboard) ID() string {
	if cd.Key == "" {
		return cd.Name
	}
	return cd.Name + ":" + cd.Key
}

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
