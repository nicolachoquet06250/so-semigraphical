package main

import (
	"crm/pages"
	"errors"
	"log"

	. "github.com/awesome-gocui/gocui"

	_ "github.com/go-sql-driver/mysql"
)

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
