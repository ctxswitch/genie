package ipaddr

import (
	"math/rand"
	"net"
)

type IPAddr struct {
	cidrs   []string
	uniques uint32
	cache   []string
}

func New(cfg Config) *IPAddr {
	return &IPAddr{
		cidrs:   cfg.Cidrs,
		uniques: cfg.Uniques,
	}
}

func (ip *IPAddr) Cache() []string {
	addrs := make([]string, 0)

	for _, cidr := range ip.cidrs {
		addrs = appendRange(addrs, cidr)
	}

	return addrs
}

func appendRange(ips []string, cidr string) []string {
	// The cidr is validated when the config is loaded.
	ip, ipnet, _ := net.ParseCIDR(cidr)

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); ipInc(ip) {
		ips = append(ips, ip.String())
	}

	return ips[1 : len(ips)-1]
}

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
