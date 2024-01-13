package net

import "tiangong/common/buf"

func (c Conn) ReadFrom(buffer buf.Buffer) error {
	bytes, err := buf.ReadAll(buffer)
	if err != nil {
		return err
	}
	if _, err = c.Write(bytes); err != nil {
		return err
	}
	return nil
}
