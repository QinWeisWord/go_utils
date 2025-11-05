package netutil

import (
    "net"
)

// GetLocalIPv4 获取本机所有IPv4地址（不包含回环127.0.0.1）
// 参数: 无
// 返回值: IPv4地址字符串切片与错误信息（若成功则error为nil）
// 关键步骤：遍历网卡接口并过滤IPv4与非回环地址
func GetLocalIPv4() ([]string, error) {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return nil, err
    }
    ips := make([]string, 0, len(addrs))
    for _, addr := range addrs {
        // 关键步骤：类型断言为*net.IPNet
        if ipNet, ok := addr.(*net.IPNet); ok {
            ip := ipNet.IP
            // 关键步骤：筛选IPv4且非回环
            if ip.To4() != nil && !ip.IsLoopback() {
                ips = append(ips, ip.String())
            }
        }
    }
    return ips, nil
}