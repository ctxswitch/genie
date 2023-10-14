package ipaddr

const (
	// DefaultIPAddrUniques is the default number of unique IP addresses to
	// generate.
	DefaultIPAddrUniques = 0
)

var (
	// DefaultIPAddrCidrs is the default list of CIDRs to generate IP addresses
	// from.
	DefaultIPAddrCidrs = []string{"192.168.0.0/16"}
)
