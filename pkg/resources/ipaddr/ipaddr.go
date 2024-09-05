package ipaddr

import (
	"math/rand"
	"net"
)

// IPAddr is a resource that generates a random IP address.
type IPAddr struct {
	cidrs   []string
	uniques uint32
	cache   []string
}

// New returns a new IPAddr resource.
func New(cfg Config) *IPAddr {
	return &IPAddr{
		cidrs:   cfg.Cidrs,
		uniques: cfg.Uniques,
	}
}

// Cache sets the internal cache of IP addresses.  This is currently
// not very efficient and can be improved or removed.  If the cidr is
// extremely large, it will take some time (and copious amounts of memory
// to generate and store the ips in the cache).
func (ip *IPAddr) Cache() []string {
	addrs := make([]string, 0)

	for _, cidr := range ip.cidrs {
		addrs = appendRange(addrs, cidr)
	}

	return addrs
}

// appendRange appends a range of IP addresses to the given slice.
func appendRange(ips []string, cidr string) []string {
	// The cidr is validated when the config is loaded.
	addr, ipnet, _ := net.ParseCIDR(cidr)

	for ip := addr.Mask(ipnet.Mask); ipnet.Contains(ip); ipInc(ip) {
		ips = append(ips, ip.String())
	}

	return ips[1 : len(ips)-1]
}

// ipInc increments the given IP address.
func ipInc(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

func (ip *IPAddr) Get() string {
	if ip.cache == nil {
		ip.cache = ip.Cache()
	}

	return ip.cache[rand.Intn(len(ip.cache))] //nolint:gosec
}
