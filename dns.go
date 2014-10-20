package main

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/growse/pcap"
	"github.com/mephux/dns"
)

const (
	// IPTcp TCP type code
	IPTcp = 6

	// IPUdp UDP type code
	IPUdp = 17
)

// Answer holds dns answer data
type Answer struct {
	Id         int64
	QuestionId int64
	ClientId   int64
	Class      string `json:"class"`
	Name       string `json:"name"`
	Record     string `json:"record"`
	Data       string `json:"data"`
	Ttl        string `json:"ttl"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
	Active     bool   `json:"active"`
	SeenCount  int64
}

// Message is used to pass and process data for various output options
type Question struct {
	Id        int64
	SeenCount int64
	ClientId  int64
	Answers   []Answer `json:"answers" sql:"-"`
	Question  string   `json:"question"`
	Length    int      `json:"length" sql:"-"`
	DstIp     string   `json:"dstip"`
	Protocol  string   `json:"protocol" sql:"-"`
	SrcIp     string   `json:"srcip"`
	CreatedAt int64    `json:"created_at"`
	UpdatedAt int64    `json:"updated_at"`
	Packet    string   `json:"packet" sql:"-"`
	Type      uint16   `json:"question_type"`
	Class     uint16   `json:"question_class"`
}

type Client struct {
	Id            int64
	LastSeen      int64
	QuestionCount int64  `json:"question_count"`
	Hostname      string `json:"hostname"`
	Interface     string `json:"interface"`
	MacAddr       string `json:"mac_addr"`
	Ip            string `json:"ip"`
}

// DNS process and parse DNS packets
func DNS(pkt *pcap.Packet, options *Options) (*Question, error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("DNS Decode - Recovered:", r)
			fmt.Println("Please report this issue. https://github.com/mephux/dnas")
		}
	}()

	message := &Question{}

	pkt.Decode()
	msg := new(dns.Msg)
	err := msg.Unpack(pkt.Payload)

	if err != nil || len(msg.Answer) <= 0 {
		return message, fmt.Errorf("Error")
	}

	if len(pkt.Headers) <= 0 {
		return message, fmt.Errorf("Error: Missing header information.")
	}

	message.Length = msg.Len()

	packet, _ := msg.Pack()

	message.Packet = hex.EncodeToString(packet)

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

		message.SrcIp = ip4hdr.SrcAddr()
		message.DstIp = ip4hdr.DestAddr()

	} else {
		ip6hdr, _ := pkt.Headers[0].(*pcap.Ip6hdr)

		message.SrcIp = ip6hdr.SrcAddr()
		message.DstIp = ip6hdr.DestAddr()
	}

	for i := range msg.Question {
		message.Question = msg.Question[i].Name
		message.Type = msg.Question[i].Qtype
		message.Class = msg.Question[i].Qclass
	}

	for i := range msg.Answer {
		split := strings.Split(msg.Answer[i].String(), "\t")
		answer := Answer{
			Name:      split[0],
			Ttl:       split[1],
			Class:     split[2],
			Record:    split[3],
			Data:      split[4],
			ClientId:  options.Client.Id,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		message.Answers = append(message.Answers, answer)
	}

	return message, nil
}
