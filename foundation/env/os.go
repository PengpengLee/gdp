// Author: peng.lee577 (peng.lee577@gmail.com)
// Date: 2021/11/29

package env

import (
	"errors"
	"net"
	"os"
	"strconv"
)

// IP 地址版本
const (
	IPv4 = "IPv4"
	IPv6 = "IPv6"
)

var (
	pid       int
	pidString string

	localIPV4 = "unknown"
	localIPV6 = "unknown"
)

// PID 进程ID，即 process id
func PID() int {
	return pid
}

// PIDString 得到 PID 字符串形式
// 如打印日志的场景
func PIDString() string {
	return pidString
}

// LocalIP 本机IP，返回非127域的第一个ipv4 地址
// 极端特殊情况获取失败返回 机器名 或者 unknown
func LocalIP() string {
	return localIPV4
}

// LocalIPv6 本机IP，返回非127域的第一个 ipv6 地址
// 极端特殊情况获取失败返回 机器名 或者 unknown
func LocalIPv6() string {
	return localIPV6
}

func init() {
	pid = os.Getpid()
	pidString = strconv.Itoa(pid)

	if val, err := localIP(IPv4); err == nil {
		localIPV4 = val
	}
	if val, err := localIP(IPv6); err == nil {
		localIPV6 = val
	}
}

func localIP(version string) (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return os.Hostname()
	}

	for _, a := range addrs {
		ipnet, ok := a.(*net.IPNet)
		// 过滤掉本机回环地址
		if ok && !ipnet.IP.IsLoopback() {
			switch version {
			case IPv4:
				if ipnet.IP.To4() != nil {
					return ipnet.IP.String(), nil
				}
			case IPv6:
				if ipnet.IP.To16() != nil {
					return ipnet.IP.String(), nil
				}
			}
		}
	}

	return "", errors.New("fail to get local IP: " + version)
}
