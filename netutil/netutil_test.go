package netutil

import (
    "net"
    "testing"
)

// TestGetLocalIPv4 测试获取本机IPv4地址的基本正确性
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：允许返回为空；若非空则断言均为IPv4且不包含回环地址
func TestGetLocalIPv4(t *testing.T) {
    ips, err := GetLocalIPv4()
    if err != nil { t.Fatalf("获取IPv4失败: %v", err) }
    for _, ip := range ips {
        parsed := net.ParseIP(ip)
        if parsed == nil || parsed.To4() == nil {
            t.Fatalf("应为IPv4地址: %s", ip)
        }
        if parsed.IsLoopback() {
            t.Fatalf("不应包含回环地址: %s", ip)
        }
    }
}