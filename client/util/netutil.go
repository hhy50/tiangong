package netutil

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/songgao/water/waterutil"
)

func GetAddr(b []byte) (srcAddr string, dstAddr string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			srcAddr = ""
			dstAddr = ""
		}
	}()
	if waterutil.IPv4Protocol(b) == waterutil.TCP {
		srcIp := waterutil.IPv4Source(b)
		dstIp := waterutil.IPv4Destination(b)
		srcPort, dstPort := GetPort(b)
		src := fmt.Sprintf("%s:%s", srcIp.To4().String(), srcPort)
		dst := fmt.Sprintf("%s:%s", dstIp.To4().String(), dstPort)
		//log.Printf("%s->%s", src, dst)
		return src, dst
	} else if waterutil.IPv4Protocol(b) == waterutil.UDP {
		srcIp := waterutil.IPv4Source(b)
		dstIp := waterutil.IPv4Destination(b)
		srcPort, dstPort := GetPort(b)
		src := fmt.Sprintf("%s:%s", srcIp.To4().String(), srcPort)
		dst := fmt.Sprintf("%s:%s", dstIp.To4().String(), dstPort)
		//log.Printf("%s->%s", src, dst)
		return src, dst
	} else if waterutil.IPv4Protocol(b) == waterutil.ICMP {
		srcIp := waterutil.IPv4Source(b)
		dstIp := waterutil.IPv4Destination(b)
		return srcIp.To4().String(), dstIp.To4().String()
	}
	return "", ""
}

func GetPort(b []byte) (srcPort string, dstPort string) {
	packet := gopacket.NewPacket(b, layers.LayerTypeIPv4, gopacket.Default)
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		return tcp.SrcPort.String(), tcp.DstPort.String()
	} else if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		return udp.SrcPort.String(), udp.DstPort.String()
	}
	return "", ""
}
