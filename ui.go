package main

import (
	"fmt"
	"math"

	"github.com/jroimartin/gocui"
)

type Day struct {
	name string
    x, y int
    w, h int
    body string
}

type Week struct {
	name string
    x, y int
    w, h int
    body string
	days []*Day
}

func newDay(name string, body string) *Day {
    return &Day{name: name, x: 0, y: 0, w: 0, h: 0, body: body}
}

func (w *Day) Layout(g *gocui.Gui) error {
    v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
    if err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        fmt.Fprintln(v, w.body)
    }
    return nil
}

func (d *Day) setPropreties(x, y, w, h int) {
    d.x = x
    d.y = y
    d.w = w
    d.h = h
}

func newWeek(name string, days []*Day, body string) *Week {
    return &Week{name: name, x: 0, y: 0, w: 0, h: 0, body: body, days: days}
}

func (we *Week) setPropreties(x, y, w, h int) {
    we.x = x
    we.y = y
    we.w = w
    we.h = h
}

func (w *Week) getDayDimensions() (width, border int) {
    numberOfDay := len(w.days)
	n := math.Round(float64((w.w) / numberOfDay))
	b := math.Round(n * 0.1)
	wi := n - b

	return int(wi), int(b)
}

func (w *Week) updateDaysView(g *gocui.Gui) {
    width, border := w.getDayDimensions()
    y0 := 2
    x0 := 2
	x := x0
	for _, v := range w.days {
        v.setPropreties(x, y0, width-x0, w.h-y0-1)
        v.Layout(g)
		x += width + border
	}
}

func (w *Week) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
    height := math.Round(float64(maxY)*0.66)

    w.setPropreties(0, 0, maxX-1, int(height))

    view, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        fmt.Fprintln(view, w.body)
	}

    w.updateDaysView(g)

	return nil
}
