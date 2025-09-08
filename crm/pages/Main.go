package pages

import (
	"crm/lib"
	"crm/models/user"
	"errors"
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

type Main struct {
	lib.Layout
	lib.Renderable

	User *user.User

	InitialMouseX   int
	InitialMouseY   int
	XOffset         int
	YOffset         int
	MsgMouseDown    bool
	MovingMsg       bool
	GlobalMouseDown bool
}

func (m *Main) mouseUp(g *gocui.Gui, v *gocui.View) error {
	if m.MsgMouseDown {
		m.MsgMouseDown = false
		if m.MovingMsg {
			m.MovingMsg = false
			return nil
		} else {
			_ = g.DeleteView("msg")
		}
	} else if m.GlobalMouseDown {
		m.GlobalMouseDown = false
		_ = g.DeleteView("globalDown")
	}
	return nil
}

func (m *Main) msgDown(g *gocui.Gui, v *gocui.View) error {
	m.InitialMouseX, m.InitialMouseY = g.MousePosition()
	if vx, vy, _, _, err := g.ViewPosition("msg"); err == nil {
		m.XOffset = m.InitialMouseX - vx
		m.YOffset = m.InitialMouseY - vy
		m.MsgMouseDown = true
	}
	return nil
}

func (m *Main) showMsg(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	_, cy := v.Cursor()
	if l, err = v.Line(cy + 2); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-10, maxY/2, maxX/2+10, maxY/2+2, 0); err == nil || errors.Is(err, gocui.ErrUnknownView) {
		v.Clear()
		v.SelBgColor = gocui.ColorCyan
		v.SelFgColor = gocui.ColorBlack
		_, _ = fmt.Fprintln(v, l)
	}
	return nil
}

func (m *Main) Render(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	g.Cursor = true
	g.Mouse = true

	if v, err := g.SetView("main-frame", 1, 1, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Frame = true
		v.Title = "Main"
	}

	if v, err := g.SetView("user", maxX-30, 2, maxX-2, 6, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Frame = true

		_, _ = fmt.Fprintf(v, "FirstName: %s\nLastName: %s\nEmail: %s", m.User.FirstName, m.User.LastName, m.User.Email)
	}

	if v, err := g.SetView("customer-list", 2, 2, 50, maxY-2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Frame = true
		v.Title = "Customer List"
		v.Autoscroll = true

		_, _ = fmt.Fprintf(v, "Nicolas Choquet\nYann Choquet\nGregory Choquet\nKarine Amgar\nLaure Amgar\nMichel Amgar\nNicolas Choquet\nYann Choquet\nGregory Choquet\nKarine Amgar\nLaure Amgar\nMichel Amgar\nNicolas Choquet\nYann Choquet\nGregory Choquet\nKarine Amgar\nLaure Amgar\nMichel Amgar\nNicolas Choquet\nYann Choquet\nGregory Choquet\nKarine Amgar\nLaure Amgar\nMichel Amgar\nNicolas Choquet\nYann Choquet\nGregory Choquet\nKarine Amgar\nLaure Amgar\nMichel Amgar\nNicolas Choquet\nYann Choquet\nGregory Choquet\nKarine Amgar\nLaure Amgar\nMichel Amgar\nNicolas Choquet\nYann Choquet\nGregory Choquet\nKarine Amgar\nLaure Amgar\nMichel Amgar\nNicolas Choquet\nYann Choquet\nGregory Choquet\nKarine Amgar\nLaure Amgar\nMichel Amgar\nNicolas Choquet\nYann Choquet\nGregory Choquet\nKarine Amgar\nLaure Amgar\nMichel Amgar\nNicolas Choquet\nYann Choquet\nGregory Choquet\nKarine Amgar\nLaure Amgar\nMichel Amgar\n")
	}

	if err := g.SetKeybinding("customer-list", gocui.MouseLeft, gocui.ModNone, m.showMsg); err != nil {
		return err
	}
	if err := g.SetKeybinding("customer-list", gocui.MouseRelease, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		_ = g.DeleteView("msg")

		return nil
	}); err != nil {
		return err
	}
	//if err := g.SetKeybinding("", gocui.MouseRelease, gocui.ModNone, m.mouseUp); err != nil {
	//	return err
	//}
	//if err := g.SetKeybinding("msg", gocui.MouseLeft, gocui.ModNone, m.msgDown); err != nil {
	//	return err
	//}

	if err := m.Layout.Keybindings(g); err != nil {
		log.Panicln(err)
	}

	return nil
}
