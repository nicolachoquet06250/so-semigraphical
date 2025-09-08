package main

import (
	"crm/pages"
	"errors"
	"log"

	. "github.com/awesome-gocui/gocui"

	_ "github.com/go-sql-driver/mysql"
)

/*var (
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
		_, _ = fmt.Fprintln(v, "Login (Press Enter)")
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
		v.Title = "Email"
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

func layout2(user user.User) func(g *Gui) error {
	return func(g *Gui) error {
		maxX, maxY := g.Size()

		if v, err := g.SetView("main-frame", 1, 1, maxX-1, maxY-1, 0); err != nil {
			if !errors.Is(err, ErrUnknownView) {
				return err
			}
			v.Frame = true
			v.Title = "Main"
		}

		if v, err := g.SetView("user", maxX-30, 2, maxX-2, 6, 0); err != nil {
			if !errors.Is(err, ErrUnknownView) {
				return err
			}
			v.Frame = true

			_, _ = fmt.Fprintf(v, "FirstName: %s\nLastName: %s\nEmail: %s", user.FirstName, user.LastName, user.Email)
		}

		if err := keybindings(g); err != nil {
			log.Panicln(err)
		}

		return nil
	}
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
			if !errors.Is(err, ErrUnknownView) {
				return err
			}
			v.Clear()
			v.Title = "Result"

			row := db.QueryRow("SELECT id, first_name, last_name, email FROM user WHERE email=? AND password=SHA1(?)", email, password)
			var user = user.User{}
			if err = row.Scan(
				&user.Id, &user.FirstName,
				&user.LastName, &user.Email,
			); err != nil {
				_ = g.DeleteView("login")
				_ = createLoginButton(g, maxX, maxY)
				_, _ = fmt.Fprintf(v, "Wrong password")
			} else {
				v.Clear()
				g.SetManagerFunc(layout2(user))

				if err := g.MainLoop(); err != nil && !errors.Is(err, ErrQuit) {
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
}*/

func main() {
	g, err := NewGui(OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	var layout = &pages.Login{
		Email:    "nchoquet@norsys.fr",
		Password: "12041998Yann?!",
	}

	layout.SetWidgets([]string{"email", "password", "login"})

	g.SetManagerFunc(layout.Render)
	//
	//if err := keybindings(g); err != nil {
	//	log.Panicln(err)
	//}

	if err := g.MainLoop(); err != nil && !errors.Is(err, ErrQuit) {
		log.Panicln(err)
	}
}
