package model

import (
	"fmt"

	"github.com/cuigh/auxo/data"
	"github.com/cuigh/auxo/util/cast"
)

// Chart represents a dashboard chart.
type Chart struct {
	Name        string        `json:"name" bson:"_id" valid:"required"` // unique, the name of build-in charts has '$' prefix.
	Title       string        `json:"title" valid:"required"`
	Description string        `json:"desc"`
	Legend      string        `json:"-"`
	Query       string        `json:"-"`
	Metrics     []ChartMetric `json:"metrics" valid:"required"`
	Kind        string        `json:"kind"`      // builtin/custom
	Dashboard   string        `json:"dashboard"` // home/service/task...
	Type        string        `json:"type"`      // pie/line...
	Unit        string        `json:"unit"`      // bytes/milliseconds/percent:100...
	Width       int32         `json:"width"`     // 1-12(12 columns total)
	Height      int32         `json:"height"`    // default 50
	Options     data.Map      `json:"options"`
	//Colors      []string `json:"colors"`
}

func NewChart(dashboard, name, title, legend, query, unit string) *Chart {
	return &Chart{
		Name:        name,
		Title:       title,
		Description: title,
		Metrics: []ChartMetric{
			{Legend: legend, Query: query},
		},
		Dashboard: dashboard,
		Type:      "line",
		Unit:      unit,
		Width:     12,
		Height:    200,
	}
}

type ChartMetric struct {
	Legend string `json:"legend"`
	Query  string `json:"query"`
}

type ChartOption struct {
	Name   string `json:"name"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
	//Colors []string `json:"colors"`
}

type ChartDashboard struct {
	Name            string        `json:"name"`
	Key             string        `json:"key"`
	Period          int32         `json:"period"`           // minutes
	RefreshInterval int32         `json:"refresh_interval"` // seconds, 0 means disabled.
	Charts          []ChartOption `json:"charts"`
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

func (p *ChartPoint) MarshalJSON() ([]byte, error) {
	return cast.StringToBytes(fmt.Sprintf("[%v,%v]", p.X, p.Y)), nil
}

type ChartLine struct {
	Name string       `json:"name"`
	Data []ChartPoint `json:"data"`
}

type ChartMatrixData struct {
	Legend []string    `json:"legend"`
	Series []ChartLine `json:"series"`
}

type ChartValue struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type ChartVectorData struct {
	Legend []string     `json:"legend"`
	Data   []ChartValue `json:"data"`
}
