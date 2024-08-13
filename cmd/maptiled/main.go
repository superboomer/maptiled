package main

import (
	"log"
	"os"

	"github.com/superboomer/maptiled/internal/commands"
	"github.com/superboomer/maptiled/internal/options"
	"github.com/umputun/go-flags"
)

// Version contains build version
var Version = "dev"

func main() {
	var Opts = &options.Opts{}
	p := flags.NewParser(Opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	p.SubcommandsOptional = true
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			log.Fatalf("flags error: %v", err)
		}
		os.Exit(1)
	}

	log.Print("build version=", Version)

	if err := commands.TUI(Opts); err != nil {
		log.Fatal("fatal error: ", err)
	}
}
