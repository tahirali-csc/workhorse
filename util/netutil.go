package util

import (
	"net"
	"strings"
)

func GetHostIPAddress() string {
	// host, _ := os.Hostname()
	// addrs, _ := net.LookupIP(host)

	// for _, addr := range addrs {
	// 	if ipv4 := addr.To4(); ipv4 != nil {
	// 		fmt.Println("IPv4: ", ipv4)
	// 	}
	// }

	//TODO: Review to find the logic of finding primary outbound network interface
	addr, err := net.InterfaceAddrs()
	if err == nil {
		for _, v := range addr {
			ip := v.String()
			//The ip address is in CIDR format
			if strings.Contains(ip, "192.") {
				return strings.Split(ip, "/")[0]
			}
		}
	}

	return ""
}
