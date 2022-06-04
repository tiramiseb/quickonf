package datastores

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/ini.v1"
)

// TODO Use nmcli instead of reading the ini files

const networkManagerConnectionsPath = "/etc/NetworkManager/system-connections"

var NetworkManagerConnections = networkManagerConnectionsList{
	connections: map[string]NetworkManagerConnection{},
}

type NetworkManagerConnectionConnection struct {
	UUID          string `ini:"uuid"`
	Type          string `ini:"type"`
	InterfaceName string `ini:"interface-name"`
	Autoconnect   string `ini:"autoconnect"`
}

type NetworkManagerConnectionWifi struct {
	Mode string `ini:"mode"`
	SSID string `ini:"ssid"`
}

type NetworkManagerConnectionWifiSecurity struct {
	PSK string `ini:"psk"`
}

type NetworkManagerConnection struct {
	NetworkManagerConnectionConnection   `ini:"connection"`
	NetworkManagerConnectionWifi         `ini:"wifi"`
	NetworkManagerConnectionWifiSecurity `ini:"wifi-security"`
}

func (n NetworkManagerConnection) String() string {
	var str strings.Builder
	str.WriteString("UUID: ")
	str.WriteString(n.UUID)
	str.WriteString("\nType: ")
	str.WriteString(n.Type)
	if n.InterfaceName != "" {
		str.WriteString("\nInterface: ")
		str.WriteString(n.InterfaceName)
	}
	if n.Autoconnect != "" {
		str.WriteString("\nAutoconnect: ")
		str.WriteString(n.Autoconnect)
	}
	if n.Mode != "" {
		str.WriteString("\nMode: ")
		str.WriteString(n.Mode)
	}
	if n.SSID != "" {
		str.WriteString("\nSSID: ")
		str.WriteString(n.SSID)
	}
	if n.PSK != "" {
		str.WriteString("\nPSK:")
		str.WriteString(n.PSK)
	}
	str.WriteRune('\n')
	return str.String()
}

type networkManagerConnectionsList struct {
	mutex       sync.Mutex
	connections map[string]NetworkManagerConnection
}

func (n *networkManagerConnectionsList) Get(name string) (conn NetworkManagerConnection, exists bool, err error) {
	n.mutex.Lock()
	conn, exists = n.connections[name]
	n.mutex.Unlock()
	if !exists {
		conn, err = loadNetworkManagerConnection(name)
		if err == nil {
			exists = true
			n.mutex.Lock()
			n.connections[name] = conn
			n.mutex.Unlock()
		} else if errors.Is(err, fs.ErrNotExist) {
			err = nil
		}
	}
	return
}

func (n *networkManagerConnectionsList) Reset() {
	n.mutex.Lock()
	n.connections = map[string]NetworkManagerConnection{}
	n.mutex.Unlock()
}

func loadNetworkManagerConnection(name string) (NetworkManagerConnection, error) {
	n := NetworkManagerConnection{}
	cfg, err := ini.Load(filepath.Join(networkManagerConnectionsPath, name+".nmconnection"))
	if err != nil {
		return n, err
	}
	err = cfg.MapTo(&n)
	return n, err
}
