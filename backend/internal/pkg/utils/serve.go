package utils

import "net"

/*
GetServerIp 返回服务器的外部IP地址。 如果无法获取外部IP地址，则返回空字符串。

文档：
  - https://www.jianshu.com/p/301aabc06972
  - https://www.jianshu.com/p/301aabc06972
*/
func GetServerIp() string {
	ip, err := externalIP()
	if err != nil {
		return ""
	}
	return ip.String()
}

// externalIP 返回服务器的外部IP地址和一个错误值。
// 如果无法获取外部IP地址，则返回nil和相应的错误。
func externalIP() (net.IP, error) {
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// 遍历所有网络接口
	for _, iface := range interfaces {
		// 跳过已关闭的接口
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		// 跳过环回接口
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		// 获取接口的地址列表
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		// 遍历接口的地址列表
		for _, addr := range addrs {
			// 从地址中提取IP地址
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			// 返回找到的第一个非环回、非本地链接的外部IP地址
			return ip, nil
		}
	}

	// 如果没有找到外部IP地址，则返回nil和相应的错误
	return nil, err
}

// getIpFromAddr 从给定的地址中提取IP地址。
// 如果地址不是IP地址，则返回nil。
func getIpFromAddr(addr net.Addr) net.IP {
	switch v := addr.(type) {
	case *net.IPNet:
		return v.IP
	case *net.IPAddr:
		return v.IP
	}
	return nil
}
