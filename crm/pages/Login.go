package pages

import (
	"crm/lib"
	"crm/models/user"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/awesome-gocui/gocui"
)

type Login struct {
	lib.Layout
	lib.Renderable

	Email    string
	Password string
}

func (l *Login) doLogin(g *gocui.Gui, v *gocui.View) error {
	// Récupérer email et mot de passe
	if ev, err := g.View("email"); err == nil {
		l.Email = strings.TrimSpace(ev.Buffer())
	}
	if pv, err := g.View("password"); err == nil {
		l.Password = strings.TrimSpace(pv.Buffer())
	}

	// Affichage du résultat
	g.Update(func(g *gocui.Gui) error {
		_ = g.DeleteView("result")
		maxX, maxY := g.Size()
		db, err := sql.Open("mysql", "nchoquet:nchoquet@tcp(127.0.0.1:3306)/go_semigraphical_app")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if v, err := g.SetView("result", maxX/2-20, maxY/2+7, maxX/2+20, maxY/2+9, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
			v.Clear()
			v.Title = "Result"

			row := db.QueryRow("SELECT id, first_name, last_name, email FROM user WHERE email=? AND password=SHA1(?)", l.Email, l.Password)
			var user = user.User{}
			if err = row.Scan(
				&user.Id, &user.FirstName,
				&user.LastName, &user.Email,
			); err != nil {
				_ = g.DeleteView("login")
				_ = l.createLoginButton(g, maxX, maxY)
				_, _ = fmt.Fprintf(v, "Wrong password")
			} else {
				v.Clear()
				var layout = &Main{
					User: &user,
				}

				g.SetManagerFunc(layout.Render)

				if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
					log.Panicln(err)
				}
			}
		}
		return nil
	})

	return nil
}

func (l *Login) createLoginButton(g *gocui.Gui, maxX, maxY int) error {
	if v, err := g.SetView("login", maxX/2-10, maxY/2+3, maxX/2+10, maxY/2+5, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Clear()
		_, _ = fmt.Fprintln(v, "Login (Press Enter)")
	}

	if err := g.SetKeybinding("login", gocui.KeyEnter, gocui.ModNone, l.doLogin); err != nil {
		return err
	}

	return nil
}

func (l *Login) createPasswordField(g *gocui.Gui, maxX, maxY int) error {
	if v, err := g.SetView("password", maxX/2-20, maxY/2-1, maxX/2+20, maxY/2+1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Password"
		v.Editable = true
		v.Wrap = false
		v.Mask = '*'

		if l.Password != "" {
			_, _ = fmt.Fprintf(v, l.Password)
			_ = v.SetCursor(0, len(l.Password))
		}
	}

	return nil
}

func (l *Login) createEmailField(g *gocui.Gui, maxX, maxY int) error {
	if v, err := g.SetView("email", maxX/2-20, maxY/2-4, maxX/2+20, maxY/2-2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Email"
		v.Editable = true
		v.Wrap = false

		if l.Email != "" {
			_, _ = fmt.Fprintf(v, l.Email)
			_ = v.SetCursor(0, len(l.Email))
		}
	}

	return nil
}

func (l *Login) Render(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	g.Cursor = true

	if v, err := g.SetView("login-frame", 1, 1, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Frame = true
		v.Title = "Login"
	}

	// Champ email
	if err := l.createEmailField(g, maxX, maxY); err != nil {
		return err
	}

	// Champ mot de passe
	if err := l.createPasswordField(g, maxX, maxY); err != nil {
		return err
	}

	// Bouton login
	if err := l.createLoginButton(g, maxX, maxY); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(l.Layout.Widgets[l.Layout.CurrentIndex]); err != nil {
		return err
	}
	_ = l.Layout.Keybindings(g)
	l.Layout.UpdateCursor(g)

	return nil
}
