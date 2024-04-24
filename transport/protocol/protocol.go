package protocol

type Protocol = byte

type Status = byte

const (
	Unknown Protocol = iota
	TCP
	UDP
	ICMP
	HTTP
	HTTPS
	WS

	New Status = iota
	Active
	End
)

var ()

func ProtocolToStr(p Protocol) string {
	switch p {
	case TCP:
		return "TCP"
	case UDP:
		return "UDP"
	case ICMP:
		return "ICMP"
	default:
		return "Unknown"
	}
}
