package generator

import (
	"bytes"
	"encoding/binary"
	"net"
)

func IDByIP(ip string) uint32 {
	var id uint32
	err := binary.Read(bytes.NewBuffer(net.ParseIP(ip).To4()), binary.BigEndian, &id)
	if err != nil {
		return 0
	}

	return id
}
