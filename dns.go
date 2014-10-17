package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/growse/pcap"
	"github.com/miekg/dns"
)

const (
	// IPTcp TCP type code
	IPTcp = 6

	// IPUdp UDP type code
	IPUdp = 17
)

// Answer holds dns answer data
type Answer struct {
	Class     string    `json:"class"`
	Name      string    `json:"name"`
	Record    string    `json:"record"`
	Data      string    `json:"data"`
	TTL       string    `json:"ttl"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Active    bool      `json:"active"`
}

// Message is used to pass and process data for various output options
type Message struct {
	DNS struct {
		Answers  []Answer `json:"answers"`
		Question string   `json:"question"`
		Length   int      `json:"length"`
	} `json:"dns"`
	DstIP     string    `json:"dstip"`
	Protocol  string    `json:"protocol"`
	SrcIP     string    `json:"srcip"`
	Timestamp time.Time `json:"timestamp"`
	Packet    []byte    `json:"packet"`
	Bloom     []byte    `json:"bloom"`
}

// DNS process and parse DNS packets
func DNS(pkt *pcap.Packet) (*Message, error) {
	message := &Message{}

	message.Timestamp = time.Now()

	pkt.Decode()
	msg := new(dns.Msg)
	err := msg.Unpack(pkt.Payload)

	if err != nil || len(msg.Answer) <= 0 {
		return message, fmt.Errorf("Error")
	}

	if len(pkt.Headers) <= 0 {
		return message, fmt.Errorf("Error: Missing header information.")
	}

	message.DNS.Length = msg.Len()

	packet, _ := msg.Pack()

	message.Packet = packet

	ip4hdr, ip4ok := pkt.Headers[0].(*pcap.Iphdr)

	if ip4ok {

		switch ip4hdr.Protocol {
		case IPTcp:
			message.Protocol = "TCP"
		case IPUdp:
			message.Protocol = "UDP"
		default:
			message.Protocol = "N/A"
		}

		message.SrcIP = ip4hdr.SrcAddr()
		message.DstIP = ip4hdr.DestAddr()

	} else {
		ip6hdr, _ := pkt.Headers[0].(*pcap.Ip6hdr)

		message.SrcIP = ip6hdr.SrcAddr()
		message.DstIP = ip6hdr.DestAddr()
		fmt.Println(ip6hdr)
	}

	for i := range msg.Question {
		message.DNS.Question = msg.Question[i].Name
	}

	for i := range msg.Answer {
		split := strings.Split(msg.Answer[i].String(), "\t")
		answer := Answer{
			Name:      split[0],
			TTL:       split[1],
			Class:     split[2],
			Record:    split[3],
			Data:      split[4],
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Active:    true,
		}

		message.DNS.Answers = append(message.DNS.Answers, answer)
	}

	return message, nil
}
