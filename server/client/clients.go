package client

import "tiangong/common/errors"

var (
	Clients = make(map[string]*Client, 128)
)

func AddClient(name string, client *Client) error {
	if _, f := Clients[name]; f {
		return errors.NewError("Unable to add existing client, name: "+name, nil)
	}
	Clients[name] = client
	return nil
}
