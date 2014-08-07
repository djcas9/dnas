package main

import (
	"fmt"
	"github.com/growse/pcap"
	"os"
)

var (
	snaplen = 65536
)

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

	for pkt, r := h.NextEx(); r >= 0; pkt, r = h.NextEx() {

		if r == 0 {
			// timeout, continue
			continue
		}

		message, err := DNS(pkt)

		if err == nil {
			json, err := message.ToJSON()

			if err != nil {
				panic(err)
			}

			fmt.Println(string(json))

			message.ToStdout()
		}

	}

	fmt.Fprintf(os.Stderr, "%s Error: %s\n", NAME, h.Geterror())
}
