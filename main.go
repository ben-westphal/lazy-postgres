package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

var viewNames = []string{"project", "databases", "queries", "queryWindow"}
var currentView = 0
var didInitFocus = false
var showAddDatabase = false

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	cv := g.CurrentView()
	current := ""
	if cv != nil {
		current = cv.Name()
	}

	makeView := func(name string, x0, y0, x1, y1 int, title, content string, frame bool) error {
		v, err := g.SetView(name, x0, y0, x1, y1)
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}
		if err == gocui.ErrUnknownView {
			v.Title = title
			v.Wrap = true
			v.Highlight = (name == current)
			v.Frame = frame
			if content != "" {
				fmt.Fprint(v, content)
			}
		} else {
			v.Highlight = (name == current)
		}
		return nil
	}

	if err := makeView("project", 0, 0, maxX/3-1, 2, "Project", "lazy-postgres", true); err != nil {
		return err
	}
	if err := makeView("databases", 0, 3, maxX/3-1, maxY/2, "Databases", "Database 1\nDatabase 2\nDatabase 3\n", true); err != nil {
		return err
	}
	if err := makeView("queries", 0, maxY/2+1, maxX/3-1, maxY-3, "Queries", "", true); err != nil {
		return err
	}
	if err := makeView("queryWindow", maxX/3, 0, maxX-1, maxY-3, "Query Window", "", true); err != nil {
		return err
	}
	if err := makeView("footer", 0, maxY-2, maxX-1, maxY, "", "tab: change pane, q: quit, ← ↑ → ↓ : navigate, a: add new db", false); err != nil {
		return err
	}

	if !didInitFocus {
		didInitFocus = true
		go func() {
			g.Update(func(g *gocui.Gui) error {
				_, err := g.SetCurrentView(viewNames[0])
				return err
			})
		}()
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

func toggleAddDatabase(g *gocui.Gui, v *gocui.View) error {
	showAddDatabase = !showAddDatabase
	g.Update(func(g *gocui.Gui) error {
		if showAddDatabase {
			maxX, maxY := g.Size()
			width := 60
			height := 7

			x0 := (maxX - width) / 2
			y0 := (maxY - height) / 2
			x1 := x0 + width
			y1 := y0 + height

			if v, err := g.SetView("addDatabase", x0, y0, x1, y1); err != nil {
				if err != gocui.ErrUnknownView {
					return err
				}
				v.Title = "Add Database"
				v.Editable = true
				v.Wrap = true
				fmt.Fprint(v, "Name:\n\nConnection String:\n")
				_, err := g.SetCurrentView("addDatabase")
				return err
			}
		} else {
			g.DeleteView("addDatabase")
			_, err := g.SetCurrentView(viewNames[currentView])
			return err
		}

		return nil
	})

	return nil
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.SelFgColor = gocui.ColorGreen | gocui.AttrBold
	g.SelBgColor = gocui.ColorDefault

	g.SetManagerFunc(layout)
	initKeybindings(g)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func initKeybindings(g *gocui.Gui) error {
	g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView)
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	g.SetKeybinding("", 'q', gocui.ModNone, quit)
	g.SetKeybinding("", 'a', gocui.ModNone, toggleAddDatabase)
	return nil
}

