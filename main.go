package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

var viewNames = []string{"project", "databases", "queries", "queryWindow"}
var currentView = 0

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Project view
	if v, err := g.SetView("project", 0, 0, maxX/3-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Project"
		v.Wrap = true
		fmt.Fprintln(v, "lazy-postgres")
	}

	// Databases view
	if v, err := g.SetView("databases", 0, 3, maxX/3-1, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Databases"
	}

	// Queries view
	if v, err := g.SetView("queries", 0, maxY/2+1, maxX/3-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Queries"
	}

	// Query Window
	if v, err := g.SetView("queryWindow", maxX/3, 0, maxX-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Query Window"
		v.Wrap = true
	}

	// Footer
	if v, err := g.SetView("footer", 0, maxY-2, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		fmt.Fprint(v, "PgUp/PgDn, tab: focus, q: quit, ← ↑ → ↓ : navigate")
	}

	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	currentView = (currentView + 1) % len(viewNames)
	_, err := g.SetCurrentView(viewNames[currentView])
	return err
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
