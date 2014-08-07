package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/wsxiaoys/terminal"
	_ "github.com/wsxiaoys/terminal/color"
	"os"
	"text/tabwriter"
)

var (
	w     = new(tabwriter.Writer)
	count = 0
)

type OutputJSON struct {
	Dns struct {
		Answers   []struct{} `json:"answers"`
		Questions []struct{} `json:"questions"`
	} `json:"dns"`
	Dstip          string   `json:"dstip"`
	IpHeader       struct{} `json:"ip_header"`
	ProtocolHeader struct{} `json:"protocol_header"`
	Srcip          string   `json:"srcip"`
	Timestamp      string   `json:"timestamp"`
}

func (message *Message) ToStdout() {
	w.Init(os.Stdout, 1, 2, 2, ' ', 0)

	count++
	fmt.Fprintf(w, "\t---\t%d\t---\n", count)
	fmt.Fprintf(w, "\n")

	fmt.Fprintf(w, "\tQuestion:\t\033[0;31;49m%s\033[0m\n", message.Dns.Question)
	fmt.Fprintf(w, "\tTimestamp:\t%s\n", message.Timestamp)
	fmt.Fprintf(w, "\tSource:\t%s\n", message.SrcIp)
	fmt.Fprintf(w, "\tDestination:\t%s\n\n", message.DstIp)

	fmt.Fprintf(w, "\t\033[0;32;49mAnswers (%d):\033[0m\t\n\n", len(message.Dns.Answers))

	fmt.Fprintf(w, "\tRR\tName\tData\n")
	fmt.Fprintf(w, "\t----\t----\t----\n")

	for i := range message.Dns.Answers {
		fmt.Fprintf(w, "\t%s", message.Dns.Answers[i].Record)
		fmt.Fprintf(w, "\t%s", message.Dns.Answers[i].Name)
		fmt.Fprintf(w, "\t\033[0;32;49m%s\033[0m\n", message.Dns.Answers[i].Data)
	}

	fmt.Fprintf(w, "\n")
}

func (message *Message) ToJSON() ([]byte, error) {
	return json.Marshal(message)
}
