package main

import (
	"fmt"
	"github.com/growse/pcap"
	"github.com/miekg/dns"
	"strings"
	"time"
)

const (
	TYPE_IP  = 0x0800
	TYPE_ARP = 0x0806
	TYPE_IP6 = 0x86DD

	IP_ICMP = 1
	IP_INIP = 4
	IP_TCP  = 6
	IP_UDP  = 17
)

type Answer struct {
	Class  string `json:"class"`
	Name   string `json:"name"`
	Record string `json:"record"`
	Data   string `json:"data"`
	TTL    string `json:"ttl"`
}

type Message struct {
	Dns struct {
		Answers  []Answer `json:"answers"`
		Question string   `json:"question"`
	} `json:"dns"`
	DstIp     string    `json:"dstip"`
	Protocol  string    `json:"protocol"`
	SrcIp     string    `json:"srcip"`
	Timestamp time.Time `json:"timestamp"`
}

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

	ip4hdr, ip4ok := pkt.Headers[0].(*pcap.Iphdr)

	if ip4ok {

		switch ip4hdr.Protocol {
		case IP_ICMP:
			// protohdr := pkt.Headers[1].(*pcap.Icmphdr)
			message.Protocol = "TCMP"
		case IP_TCP:
			// protohdr := pkt.Headers[1].(*pcap.Tcphdr)
			message.Protocol = "TCP"
		case IP_UDP:
			// protohdr := pkt.Headers[1].(*pcap.Udphdr)
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
		fmt.Println(ip6hdr)
	}

	// fmt.Println(pkt, msg, ggerr)

	// fmt.Println(msg.Ns)
	// fmt.Println(len(msg.Answer))

	for i := range msg.Question {
		message.Dns.Question = msg.Question[i].Name
	}

	for i := range msg.Answer {
		split := strings.Split(msg.Answer[i].String(), "\t")
		answer := Answer{
			Name:   split[0],
			TTL:    split[1],
			Class:  split[2],
			Record: split[3],
			Data:   split[4],
		}
		message.Dns.Answers = append(message.Dns.Answers, answer)
	}

	return message, nil
}
