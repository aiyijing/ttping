package tt

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ContainsPort(address string) bool {
	// Split the address into host and port
	hostPort := strings.Split(address, ":")
	if len(hostPort) < 2 {
		return false
	}

	// Check if the port is a valid integer
	_, err := strconv.Atoi(hostPort[1])
	if err != nil {
		return false
	}

	return true
}

func resolveAddress(address string) (string, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		// 解析地址失败，返回空字符串和错误信息
		return "", fmt.Errorf("解析地址失败：%s", err)
	}

	// 尝试解析 host 为 IP 地址
	ip, err := resolveHost(host)
	if err != nil {
		return "", err
	}

	// 返回第一个解析到的 IP 和端口
	return net.JoinHostPort(ip, port), nil
}

func resolveHost(host string) (string, error) {
	// 尝试解析 host 为 IP 地址
	ip := net.ParseIP(host)
	if ip != nil {
		return ip.String(), nil
	}

	// 尝试解析 host 为域名
	ips, err := net.LookupIP(host)
	if err != nil {
		// 域名解析失败，返回空字符串和错误信息
		return "", fmt.Errorf("域名解析失败：%s", err)
	}

	// 返回第一个解析到的 IP 和端口
	return ips[0].String(), nil
}
