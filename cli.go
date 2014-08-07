package main

import (
	"fmt"
	"github.com/visionmedia/go-flags"
	"os"
)

type Options struct {
	Interface string `short:"i" long:"interface" description:"Interface to monitor"`
	Daemon    bool   `short:"d" long:"daemon" description:"Run DNAS in daemon mode" default:"false"`
	Version   bool   `short:"v" long:"version" description:"Show version information" default:"false"`
	Port      int    `short:"p" long:"port" description:"DNS port" default:"53"`
}

func Usage(p *flags.Parser) {
	p.WriteHelp(os.Stdout)
	fmt.Printf("\n")
	os.Exit(1)
}

func Version() {
	fmt.Printf("%s - %s - Version: %s\n",
		NAME,
		DESCRIPTION,
		VERSION,
	)

	os.Exit(1)
}

func CLIRun() *Options {
	options := &Options{}

	var parser = flags.NewParser(options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		Usage(parser)
	}

	if options.Version {
		Version()
	}

	if options.Interface == "" {
		Usage(parser)
	}

	return options
}
