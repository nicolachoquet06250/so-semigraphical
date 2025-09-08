package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	. "github.com/awesome-gocui/gocui"

	_ "github.com/go-sql-driver/mysql"
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

func createLoginButton(g *Gui, maxX, maxY int) error {
	if v, err := g.SetView("login", maxX/2-10, maxY/2+3, maxX/2+10, maxY/2+5, 0); err != nil {
		if !errors.Is(err, ErrUnknownView) {
			return err
		}
		v.Clear()
		fmt.Fprintln(v, "Login (Press Enter)")
	}

	if err := g.SetKeybinding("login", KeyEnter, ModNone, doLogin); err != nil {
		return err
	}

	return nil
}

func createPasswordField(g *Gui, maxX, maxY int) error {
	if v, err := g.SetView("password", maxX/2-20, maxY/2-1, maxX/2+20, maxY/2+1, 0); err != nil {
		if !errors.Is(err, ErrUnknownView) {
			return err
		}
		v.Title = "Password"
		v.Editable = true
		v.Wrap = false
		v.Mask = '*'
	}

	return nil
}

func createEmailField(g *Gui, maxX, maxY int) error {
	if v, err := g.SetView("email", maxX/2-20, maxY/2-4, maxX/2+20, maxY/2-2, 0); err != nil {
		if !errors.Is(err, ErrUnknownView) {
			return err
		}
		v.Title = "Login"
		v.Editable = true
		v.Wrap = false
	}

	return nil
}

func layout(g *Gui) error {
	maxX, maxY := g.Size()

	g.Cursor = true

	if v, err := g.SetView("login-frame", 1, 1, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, ErrUnknownView) {
			return err
		}
		v.Frame = true
		v.Title = "Login"
	}

	// Champ email
	if err := createEmailField(g, maxX, maxY); err != nil {
		return err
	}

	// Champ mot de passe
	if err := createPasswordField(g, maxX, maxY); err != nil {
		return err
	}

	// Bouton login
	if err := createLoginButton(g, maxX, maxY); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(widgets[currentIndex]); err != nil {
		return err
	}
	updateCursor(g)

	return nil
}

func layout2(g *Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("main-frame", 1, 1, maxX-1, maxY-1, 0); err != nil {
		if err != ErrUnknownView {
			return err
		}
		v.Frame = true
		v.Title = "Main"
	}

	return nil
}

func updateCursor(g *Gui) {
	name := widgets[currentIndex]
	if name == "email" || name == "password" {
		g.Cursor = true
	} else {
		g.Cursor = false
	}
}

func nextView(g *Gui, v *View) error {
	currentIndex = (currentIndex + 1) % len(widgets)
	_, err := g.SetCurrentView(widgets[currentIndex])
	updateCursor(g)
	return err
}

func doLogin(g *Gui, v *View) error {
	// Récupérer email et mot de passe
	if ev, err := g.View("email"); err == nil {
		email = strings.TrimSpace(ev.Buffer())
	}
	if pv, err := g.View("password"); err == nil {
		password = strings.TrimSpace(pv.Buffer())
	}

	// Affichage du résultat
	g.Update(func(g *Gui) error {
		_ = g.DeleteView("result")
		maxX, maxY := g.Size()
		db, err := sql.Open("mysql", "nchoquet:nchoquet@tcp(127.0.0.1:3306)/go_semigraphical_app")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if v, err := g.SetView("result", maxX/2-20, maxY/2+7, maxX/2+20, maxY/2+9, 0); err != nil {
			if err != ErrUnknownView {
				return err
			}
			v.Clear()
			v.Title = "Result"

			row := db.QueryRow("SELECT id, first_name, last_name, email, password FROM user WHERE email=? AND password=SHA1(?)", email, password)
			var user = User{}
			if err = row.Scan(
				&user.Id, &user.FirstName,
				&user.LastName, &user.Email,
				&user.Password,
			); err != nil {
				_ = g.DeleteView("login")
				_ = createLoginButton(g, maxX, maxY)
				fmt.Fprintf(v, "Wrong password")
			} else {
				fmt.Fprintf(v, "FirstName: %s, LastName: %s", user.FirstName, user.LastName)
				v.Clear()
				g.SetManagerFunc(layout2)

				if err := keybindings(g); err != nil {
					log.Panicln(err)
				}

				if err := g.MainLoop(); err != nil && err != ErrQuit {
					log.Panicln(err)
				}
			}
		}
		return nil
	})

	return nil
}

func keybindings(g *Gui) error {
	// Tabulation pour naviguer
	if err := g.SetKeybinding("", KeyTab, ModNone, nextView); err != nil {
		return err
	}
	// Quitter
	if err := g.SetKeybinding("", KeyCtrlC, ModNone, quit); err != nil {
		return err
	}
	return nil
}

func quit(g *Gui, v *View) error {
	return ErrQuit
}

func main() {
	g, err := NewGui(OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != ErrQuit {
		log.Panicln(err)
	}
}
