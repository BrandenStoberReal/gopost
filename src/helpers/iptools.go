package helpers

import (
	"net"
)

func RandomIpAddress(min int, max int) string {
	octet1 := byte(RandomRange(min, max))
	octet2 := byte(RandomRange(min, max))
	octet3 := byte(RandomRange(min, max))
	octet4 := byte(RandomRange(min, max))

	return net.IPv4(octet1, octet2, octet3, octet4).String()
}
