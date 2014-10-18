package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/growse/pcap"
	"github.com/jinzhu/gorm"
)

var (
	snaplen = 65536
)

// OpenFile opens or creates a file for json logging
func OpenFile(path string) *os.File {
	var fo *os.File
	var ferr error

	if _, err := os.Stat(path); err == nil {
		fo, ferr = os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0660)
	} else {
		fo, ferr = os.Create(path)
	}

	if ferr != nil {
		panic(ferr)
	}

	return fo
}

// WriteToFile write json output to file
func WriteToFile(fo *os.File, json []byte) {
	if _, err := fo.WriteString(string(json) + "\n"); err != nil {
		panic(err)
	}
}

// Monitor bind and monitor for DNS packets. This will also
// handle the various output methods.
func Monitor(options *Options) {

	if !options.Quiet {
		fmt.Printf("\n %s (%s) - %s\n",
			Name,
			Version,
			Description,
		)

		hostname, _ := os.Hostname()

		fmt.Printf(" Host: %s\n\n", hostname)
	}

	expr := fmt.Sprintf("port %d", options.Port)

	h, err := pcap.OpenLive(options.Interface, int32(snaplen), true, 500)

	if h == nil {
		fmt.Fprintf(os.Stderr, "%s Error: %s\n", Name, err)
		os.Exit(-1)
	}

	ferr := h.SetFilter(expr)

	if ferr != nil {
		fmt.Fprintf(os.Stderr, "%s Error: %s", Name, ferr)
		os.Exit(-1)
	}

	var file *os.File

	if options.Write != "" {
		file = OpenFile(options.Write)

		defer func() {
			if err := file.Close(); err != nil {
				panic(err)
			}
		}()
	}

	if options.User != "" {
		chuser(options.User)
	}

	var db gorm.DB

	if options.Mysql {
		db, err = DatabaseConnect(options)

		if err != nil {

			fmt.Println("WOORD!!")
			fmt.Println(err.Error())
			panic(err)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			db.Close()
			file.Close()
			fmt.Println(sig)
			os.Exit(1)
		}
	}()

	for pkt, r := h.NextEx(); r >= 0; pkt, r = h.NextEx() {

		if r == 0 {
			continue
		}

		message, err := DNS(pkt)

		if err == nil {

			if options.Write != "" {
				go func() {
					json, err := message.ToJSON()

					if err != nil {
						panic(err)
					}

					WriteToFile(file, json)
				}()
			}

			if options.Mysql {
				go message.ToDatabase(db, options)
			}

			if !options.Quiet {
				message.ToStdout(options)
			}
		}

	}

	fmt.Fprintf(os.Stderr, "%s Error: %s\n", Name, h.Geterror())
}
