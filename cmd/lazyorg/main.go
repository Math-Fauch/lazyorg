package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/HubertBel/lazyorg/internal/database"
	"github.com/HubertBel/lazyorg/internal/ics"
	"github.com/HubertBel/lazyorg/internal/ui"
	"github.com/HubertBel/lazyorg/pkg/views"
	"github.com/jroimartin/gocui"
)

func main() {
	var tz int
	flag.IntVar(&tz, "tz", 0, "Offset of the timezone you are in, default chooses the pc timezone")
	file := flag.String("file", "", "File to convert to Lazyorg")
	flag.Parse()
	c, err := os.ReadFile(*file)
	if err != nil {
		panic(err)
	}
	isPCTimeZone := true
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "tz" {
			isPCTimeZone = false
		}
		if f.Name == "file" {
			ics.ConvertIcs2LO(c, tz, isPCTimeZone)
            os.Exit(0)
		}
	})
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dbDirPath := filepath.Join(homeDir, ".local", "share", "lazyorg")
	dbFilePath := filepath.Join(dbDirPath, "data.db")

	err = os.MkdirAll(dbDirPath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	database := &database.Database{}
	err = database.InitDatabase(dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.CloseDatabase()

	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	av := views.NewAppView(g, database)
	g.SetManager(av)

	if err := ui.InitKeybindings(g, av); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
