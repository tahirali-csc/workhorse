package util

import (
	"net"
	"net/http"
	"strings"
)

func GetHostIPAddress() string {
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

func GetSenderIP(request *http.Request) string {

	//forwarded := request.Header.Get("X-FORWARDED-FOR")
	//senderIP := ""
	//if forwarded != "" {
	//	senderIP = forwarded
	//} else {
	//	//senderIP = request.RemoteAddr
	//	//senderIP = request.Host
	//	senderIP = "localhost:8080"
	//}
	//[::1]:55021

	if strings.Contains(request.RemoteAddr, "[::1]") {
		return "localhost"
	}
	return strings.Split(request.RemoteAddr, ":")[0]

}