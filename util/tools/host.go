package tools

import "net"

// GetCurrentIP 获取本机IP地址
func GetCurrentIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
