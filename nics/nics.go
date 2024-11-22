package nics

import (
	"context"
	"net"

	util "github.com/takanoriyanagitani/go-nic-info/util"
)

type NetworkInterfacesSource util.Io[[]net.Interface]

func NicsSource(_ context.Context) ([]net.Interface, error){
	return net.Interfaces()
}

var NetworkInterfacesSourceDefault NetworkInterfacesSource = NicsSource
