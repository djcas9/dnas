package main

import (
	"fmt"
	"github.com/growse/pcap"
	"os"
)

var (
	snaplen = 65536
)

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

func WriteToFile(fo *os.File, json []byte) {
	if _, err := fo.WriteString(string(json) + "\n"); err != nil {
		panic(err)
	}
}

func Monitor(options *Options) {

	// hostname, err := os.Hostname()

	expr := fmt.Sprintf("port %d", options.Port)

	h, err := pcap.OpenLive(options.Interface, int32(snaplen), true, 500)

	if h == nil {
		fmt.Fprintf(os.Stderr, "%s Error: %s\n", NAME, err)
		os.Exit(-1)
	}

	ferr := h.SetFilter(expr)

	if ferr != nil {
		fmt.Fprintf(os.Stderr, "%s Error: %s", NAME, ferr)
		os.Exit(-1)
	}

	db, dberr := MakeDB(options.Database)

	fmt.Println(options.Database)

	if dberr != nil {
		panic(dberr)
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

	for pkt, r := h.NextEx(); r >= 0; pkt, r = h.NextEx() {

		if r == 0 {
			continue
		}

		message, err := DNS(pkt, options.Filter)

		if err == nil {

			if options.Write != "" {
				json, err := message.ToJSON()

				if err != nil {
					panic(err)
				}

				WriteToFile(file, json)
			}

			lvlerr := message.ToLevelDB(db, options)
			if lvlerr != nil {
				// panic(lvlerr)
			}

			message.ToStdout(options)
		}

	}

	fmt.Fprintf(os.Stderr, "%s Error: %s\n", NAME, h.Geterror())
}
