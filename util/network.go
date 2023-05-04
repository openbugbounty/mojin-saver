package util

import (
	"encoding/binary"
	"net"
	"net/url"
)

func SliceIn(item int64, list []int64) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func StringSliceIn(item string, list []string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func IPAddresses(cidr string) ([]string, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return []string{}, err
	}
	return IPAddressesIPnet(ipnet), nil
}

func IPAddressesIPnet(ipnet *net.IPNet) (ips []string) {
	// convert IPNet struct mask and address to uint32
	mask := binary.BigEndian.Uint32(ipnet.Mask)
	start := binary.BigEndian.Uint32(ipnet.IP)

	// find the final address
	finish := (start & mask) | (mask ^ 0xffffffff)

	// loop through addresses as uint32
	for i := start; i <= finish; i++ {
		// convert back to net.IP
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ips = append(ips, ip.String())
	}
	return ips
}

func IPInCIDR(ip, cidr string) bool {
	cidrIps, err := IPAddresses(cidr)
	if err != nil {
		return true
	}
	return StringSliceIn(ip, cidrIps)
}

func URL2Host(target string) (string, error) {
	u, err := url.Parse(target)
	if err != nil {
		return "", err
	}
	return u.Hostname(), nil
}
