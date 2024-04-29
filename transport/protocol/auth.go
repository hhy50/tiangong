package protocol

import "encoding/json"

type AuthType = byte
type AuthStatus = byte

const (
	AuthSession AuthType = 1
	AuthClient  AuthType = 2

	AuthFail    AuthStatus = 0
	AuthSuccess AuthStatus = 1
)

type AuthRequestPacket = Packet
type AuthResponsePacket = Packet

type ClientAuthBody struct {
	Name     string `json:"name"`
	Internal string `json:"internal"`
	Key      string `json:"key"`
	Export   string `json:"export"`
}

type SessionAuthBody struct {
	Token   string `json:"token"`
	SubHost string `json:"subHost"`
}

func (packet *AuthRequestPacket) Version() byte {
	return packet.Header.reserved[0]
}

func (packet *AuthRequestPacket) AuthType() AuthType {
	return packet.Header.reserved[1]
}

func (packet *AuthResponsePacket) AuthSuccess() bool {
	return packet.Header.reserved[0] == AuthSuccess
}

func NewAuthRequestPacket(version byte, at AuthType, body interface{}) (*AuthRequestPacket, error) {
	bytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return &AuthRequestPacket{
		Header: PacketHeader{
			Len:      uint16(len(bytes)),
			Rid:      0,
			Cmd:      AuthRequest,
			reserved: [5]byte{version, at},
		},
		Body: bytes,
	}, nil
}
func NewAuthResponsePacket(status AuthStatus) *AuthResponsePacket {
	return &AuthResponsePacket{
		Header: PacketHeader{
			Len:      0,
			Rid:      0,
			Cmd:      AuthResponse,
			reserved: [5]byte{status},
		},
		Body: EmptyBody,
	}
}
