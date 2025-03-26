// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"errors"
	"fmt"
	"net"
)

func GetLocalIPv4() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != net.FlagLoopback && iface.Flags&net.FlagUp != 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}

			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					return ipNet.IP.String(), nil
				}
			}
		}
	}
	return "", errors.New("get local IP error")
}

func MustGetLocalIPv4() string {
	ipv4, err := GetLocalIPv4()
	if err != nil {
		panic("get local IP error")
	}
	return ipv4
}

func FindIp() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("获取网络接口时出错: %v\n", err)
		return nil
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Printf("获取网络接口 %s 的地址时出错: %v\n", i.Name, err)
			continue
		}
		if i.Name == "WLAN" {
			var ips []string
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				if ip != nil {
					ips = append(ips, ip.String())
				}
			}
			return ips
		}
	}
	return nil
}
