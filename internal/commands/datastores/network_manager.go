package datastores

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/ini.v1"
)

const networkManagerConnectionsPath = "/etc/NetworkManager/system-connections"

var NetworkManagerConnections = networkManagerConnectionsList{
	connections: map[string]NetworkManagerConnection{},
}

type NetworkManagerConnectionConnection struct {
	UUID string `ini:"uuid"`
	Type string `ini:"type"`
}

type NetworkManagerConnectionWifiSecurity struct {
	PSK string `ini:"psk"`
}

type NetworkManagerConnection struct {
	NetworkManagerConnectionConnection   `ini:"connection"`
	NetworkManagerConnectionWifiSecurity `ini:"wifi-security"`
}

func (n NetworkManagerConnection) String() string {
	var str strings.Builder
	str.WriteString("UUID: ")
	str.WriteString(n.UUID)
	str.WriteString("\nType: ")
	str.WriteString(n.Type)
	str.WriteRune('\n')
	if n.PSK != "" {
		str.WriteString("PSK:")
		str.WriteString(n.PSK)
		str.WriteString("\n")
	}
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
