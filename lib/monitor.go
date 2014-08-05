package dnas

import (
	"fmt"
	"github.com/growse/pcap"
	"github.com/mephux/dnas/lib"
	"log"
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

var (
	snaplen = 65536
)

func Monitor(options interface{}) {
	fmt.Println(options)

	expr := fmt.Sprintf("port %d", options.Port)

	h, err := pcap.OpenLive(options.Interface, int32(snaplen), true, 500)

	if h == nil {
		log.Fatal(fmt.Sprintf("%s fatal error: %s", dnsa.NAME, err))
	}
}
