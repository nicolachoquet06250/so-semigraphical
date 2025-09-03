package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/awesome-gocui/gocui"

	_ "modernc.org/sqlite"
)

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

var (
	widgets      = []string{"email", "password", "login"}
	currentIndex = 0
	email        string
	password     string
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	g.Cursor = true

	if v, err := g.SetView("login-frame", 1, 1, maxX-1, maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = true
		v.Title = "Login"
	}

	// Champ email
	if v, err := g.SetView("email", maxX/2-20, maxY/2-4, maxX/2+20, maxY/2-2, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Email"
		v.Editable = true
		v.Wrap = false
	}

	// Champ mot de passe
	if v, err := g.SetView("password", maxX/2-20, maxY/2-1, maxX/2+20, maxY/2+1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Password"
		v.Editable = true
		v.Wrap = false
		v.Mask = '*'
	}

	// Bouton login
	if v, err := g.SetView("login", maxX/2-10, maxY/2+3, maxX/2+10, maxY/2+5, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, "Login (Press Enter)")
	}

	if _, err := g.SetCurrentView(widgets[currentIndex]); err != nil {
		return err
	}
	updateCursor(g)

	return nil
}

func layout2(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("main-frame", 1, 1, maxX-1, maxY-1, 0); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = true
		v.Title = "Main"
	}

	return nil
}

func updateCursor(g *gocui.Gui) {
	name := widgets[currentIndex]
	if name == "email" || name == "password" {
		g.Cursor = true
	} else {
		g.Cursor = false
	}
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	currentIndex = (currentIndex + 1) % len(widgets)
	_, err := g.SetCurrentView(widgets[currentIndex])
	updateCursor(g)
	return err
}

func doLogin(g *gocui.Gui, v *gocui.View) error {
	// Récupérer email et mot de passe
	if ev, err := g.View("email"); err == nil {
		email = strings.TrimSpace(ev.Buffer())
	}
	if pv, err := g.View("password"); err == nil {
		password = strings.TrimSpace(pv.Buffer())
	}

	// Affichage du résultat
	g.Update(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()
		db, err := sql.Open("sqlite", "sqlite.db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if v, err := g.SetView("result", maxX/2-20, maxY/2+7, maxX/2+20, maxY/2+9, 0); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Result"
			v.Clear()

			row := db.QueryRow("SELECT id, first_name, last_name, email, password FROM user WHERE email=? AND password=?", email, password)
			var user = User{}
			if err = row.Scan(
				&user.Id, &user.FirstName,
				&user.LastName, &user.Email,
				&user.Password,
			); err != nil {
				fmt.Fprintf(v, "Wrong password")
			} else {
				fmt.Fprintf(v, "FirstName: %s, LastName: %s", user.FirstName, user.LastName)
				g.SetManagerFunc(layout2)

				if err := keybindings(g); err != nil {
					log.Panicln(err)
				}

				if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
					log.Panicln(err)
				}
			}
		}
		return nil
	})

	return nil
}

func keybindings(g *gocui.Gui) error {
	// Tabulation pour naviguer
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}
	// Entrée sur le bouton login
	if err := g.SetKeybinding("login", gocui.KeyEnter, gocui.ModNone, doLogin); err != nil {
		return err
	}
	// Quitter
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
