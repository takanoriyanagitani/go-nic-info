package nic

import (
	"net"
)

type NicType string

const (
	NicTypeUnknown NicType = "unknown"
	NicTypeVirtual         = "virtual"
)

type NicInfo struct {
	net.Interface `json:"interface"`
	Addrs         []net.Addr `json:"addresses"`
	NicType       `json:"nic_type"`
}

type NetworkInterfaceToNicType func(net.Interface) NicType

func InterfaceToNicType(vnics map[string]struct{}) NetworkInterfaceToNicType {
	return func(ni net.Interface) NicType {
		var iname string = ni.Name
		_, found := vnics[iname]
		switch found {
		case true:
			return NicTypeVirtual
		default:
			return NicTypeUnknown
		}
	}
}
