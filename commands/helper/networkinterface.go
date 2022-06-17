package helper

import "net"

func NetworkInterfaceExists(iface string) (bool, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return false, err
	}
	for _, obj := range ifaces {
		if obj.Name == iface {
			return true, nil
		}
	}
	return false, nil
}
