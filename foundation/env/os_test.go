// Author: peng.lee577 (peng.lee577@gmail.com)
// Date: 2021/11/29

package env

import (
	"net"
	"strconv"
	"testing"
)

func TestPID(t *testing.T) {
	pid := PID()
	if pid <= 0 {
		t.Fatalf("pid wrong,pid=%v", pid)
	}

	for i := 0; i < 10; i++ {
		p2 := PID()
		if pid != p2 {
			t.Fatalf("pid not eq,want=%v,got=%v", pid, p2)
		}
	}

	pidStr := PIDString()
	ps := strconv.Itoa(pid)
	if ps != pidStr {
		t.Fatalf("PIDString() error, want=%v,got=%v", pid, pidStr)
	}
}

func TestLocalIP(t *testing.T) {
	ip := LocalIP()
	if len(ip) == 0 {
		t.Fatalf("LocalIP() return empty")
	}
	ipObj := net.ParseIP(ip)
	ip4 := ipObj.To4()
	if ip4 == nil {
		t.Fatalf("not ipv4")
	}

	ipv6 := LocalIPv6()
	if len(ipv6) == 0 {
		t.Fatalf("LocalIPv6() return empty")
	}
	ipObjv6 := net.ParseIP(ipv6)
	ip6 := ipObjv6.To16()
	if ip6 == nil {
		t.Fatalf("not ipv6")
	}
}
