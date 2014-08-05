package main

import (
	"fmt"
	"github.com/mephux/dnas/lib"
	"github.com/visionmedia/go-flags"
	"os"
	// "strings"
)

type Options struct {
	Interface string `short:"i" long:"interface" description:"Interface to monitor"`
	Daemon    bool   `short:"d" long:"daemon" description:"Run DNAS in daemon mode" default:"false"`
	Verbose   []bool `short:"V" long:"verbose" description:"Show verbose logging" default:"false"`
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
		dnas.NAME,
		dnas.DESCRIPTION,
		dnas.VERSION,
	)

	os.Exit(1)
}

func main() {
	options := &Options{}

	var parser = flags.NewParser(options, flags.Default)
	_, err := parser.Parse()

	if err != nil {
		Usage(parser)
	}

	if options.Version {
		Version()
	}

	if options.Interface == "" {
		Usage(parser)
	} else {
		dnas.Monitor(options)
	}

}
