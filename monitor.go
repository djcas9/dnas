package main

import (
	"fmt"
	"os"

	"github.com/growse/pcap"
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

func WriteToDatabase(message *Message, options *Options) {
	kverr := message.ToKVDB(options)

	if kverr != nil {
		// panic(lvlerr)
	}
}

func Monitor(options *Options) {

	fmt.Printf("\n %s (%s) - %s\n",
		NAME,
		VERSION,
		DESCRIPTION,
	)

	hostname, _ := os.Hostname()

	fmt.Printf(" Host: %s\n\n", hostname)

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
				go func() {
					json, err := message.ToJSON()

					if err != nil {
						panic(err)
					}

					WriteToFile(file, json)
				}()
			}

			go WriteToDatabase(message, options)
			message.ToStdout(options)
		}

	}

	fmt.Fprintf(os.Stderr, "%s Error: %s\n", NAME, h.Geterror())
}
