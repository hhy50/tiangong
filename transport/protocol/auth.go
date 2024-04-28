package protocol

type AuthType = byte
type AuthStatus = byte

const (
	// AuthType
	AuthSession AuthType = 1
	AuthClient  AuthType = 2

	// AuthStatus
	AuthFail    AuthStatus = 0
	AuthSuccess AuthStatus = 1
)

var (
	AuthResponseLen = PacketHeaderLen
)

type AuthPacketHeader = PacketHeader

// ClientAuthBody
type ClientAuthBody struct {
	Name     string `json:"name"`
	Internal string `json:"internal"`
	Key      string `json:"key"`
	Export   string `json:"export"`
}

// SessionAuthBody
type SessionAuthBody struct {
	Token   string `json:"token"`
	SubHost string `json:"subHost"`
}

func (packet *AuthPacketHeader) Version() byte {
	return packet.reserved[0]
}

func (packet *AuthPacketHeader) AuthType() AuthType {
	return packet.reserved[1]
}

func (packet *AuthPacketHeader) SetVersion(v byte) {
	packet.reserved[0] = v
}

func (packet *AuthPacketHeader) SetType(t byte) {
	packet.reserved[1] = t
}

func (packet *AuthPacketHeader) AuthSuccess() bool {
	return packet.reserved[0] == AuthSuccess
}

func NewAuthResponse(status Status) *AuthPacketHeader {
	return &AuthPacketHeader{
		Len:      0,
		Rid:      0,
		Cmd:      AuthResponse,
		reserved: [5]byte{status},
	}
}
