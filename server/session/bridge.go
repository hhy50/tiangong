package session

import (
	"bufio"
	"net/http"

	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/server/client"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

var HTTP_CLIENT = &http.Client{}

type Bridge interface {
	Transport(protocol.PacketHeader, buf.Buffer) error
}

// WirelessBridging point to point
type WirelessBridging struct {
	dst *client.Client
}

func (w *WirelessBridging) Transport(h protocol.PacketHeader, buffer buf.Buffer) error {
	var httpReq *http.Request
	var httpResp *http.Response
	var err error
	defer func() {
		if httpReq != nil {
			httpReq.Body.Close()
		}
		if httpResp != nil {
			httpResp.Body.Close()
		}
	}()

	httpReq, err = http.ReadRequest(bufio.NewReader(buffer))
	if err != nil {
		return err
	}

	httpResp, err = HTTP_CLIENT.Do(httpReq)
	if err != nil {
		return err
	}

	bytes, _ := buf.ReadAll(buffer)
	log.Info("[%s]", bytes[:httpResp.ContentLength])
	return nil
}
