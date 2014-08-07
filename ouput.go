package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/wsxiaoys/terminal"
	"github.com/wsxiaoys/terminal/color"
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
	fmt.Println("\tSource:\t", message.SrcIp)
	fmt.Println("\tDestination:\t", message.DstIp)
	fmt.Println("\tTimestamp:\t", message.Timestamp)
	fmt.Println("\tQuestion:\t", message.Dns.Question)
	fmt.Println("\n")
	color.Printf("\t@gAnswers (%d):\n", len(message.Dns.Answers))
	for i := range message.Dns.Answers {

		color.Printf("\t%s", message.Dns.Answers[i].Record)
		color.Printf("\t%s", message.Dns.Answers[i].Name)
		color.Printf("\t@r%s\n", message.Dns.Answers[i].Data)

	}
	fmt.Println("\n\n")
}

func (message *Message) ToJSON() ([]byte, error) {
	return json.Marshal(message)
}
