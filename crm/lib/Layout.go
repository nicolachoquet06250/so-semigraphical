package lib

import "github.com/awesome-gocui/gocui"

type Layout struct {
	Widgets      []string
	CurrentIndex int
	BgColor      *gocui.Attribute
	FgColor      *gocui.Attribute
	FrameColor   *gocui.Attribute
}

type Renderable interface {
	Render(g *gocui.Gui) error
}

func (l *Layout) UpdateCursor(g *gocui.Gui) {
	name := l.Widgets[l.CurrentIndex]
	if name == "email" || name == "password" {
		g.Cursor = true
	} else {
		g.Cursor = false
	}
}
func (l *Layout) Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func (l *Layout) NextView(g *gocui.Gui, v *gocui.View) error {
	l.CurrentIndex = (l.CurrentIndex + 1) % len(l.Widgets)
	_, err := g.SetCurrentView(l.Widgets[l.CurrentIndex])
	l.UpdateCursor(g)
	return err
}
func (l *Layout) Keybindings(g *gocui.Gui) error {
	// Tabulation pour naviguer
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, l.NextView); err != nil {
		return err
	}
	// Quitter
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, l.Quit); err != nil {
		return err
	}
	return nil
}

func (l *Layout) SetWidgets(w []string) *Layout {
	l.Widgets = w

	return l
}

func (l *Layout) SetCurrentIndex(i int) *Layout {
	l.CurrentIndex = i

	return l
}

func (l *Layout) SetBgColor(c *gocui.Attribute) *Layout {
	l.BgColor = c

	return l
}

func (l *Layout) SetFgColor(c *gocui.Attribute) *Layout {
	l.FgColor = c

	return l
}

func (l *Layout) SetFrameColor(c *gocui.Attribute) *Layout {
	l.FrameColor = c

	return l
}
