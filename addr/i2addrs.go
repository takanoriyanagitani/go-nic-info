package addr

import (
	"net"
	"context"
	"iter"

	ni "github.com/takanoriyanagitani/go-nic-info"
	util "github.com/takanoriyanagitani/go-nic-info/util"
)

type NetworkInterfaceToAddrs func(*net.Interface) util.Io[[]net.Addr]

func InterfaceToAddrs(ni *net.Interface) util.Io[[]net.Addr] {
	return func(_ context.Context)([]net.Addr, error){
		return ni.Addrs()
	}
}

var NetworkInterfaceToAddrsDefault NetworkInterfaceToAddrs = InterfaceToAddrs

type NicInfoGen struct {
	NetworkInterfaceToAddrs
	ni.NetworkInterfaceToNicType
}

func (g NicInfoGen) ToNicInfo() func(net.Interface) util.Io[ni.NicInfo] {
	return func(i net.Interface) util.Io[ni.NicInfo] {
		return func(ctx context.Context) (n ni.NicInfo, e error){
			addrs, e := g.NetworkInterfaceToAddrs(&i)(ctx)
			if nil != e {
				return n, e
			}

			var ntyp ni.NicType = g.NetworkInterfaceToNicType(i)

			return ni.NicInfo{
				Interface: i,
				Addrs: addrs,
				NicType: ntyp,
			}, nil
		}
	}
}

func (g NicInfoGen) InterfacesToNicsInfo(
	interfaces util.Io[[]net.Interface],
) util.Io[iter.Seq2[ni.NicInfo, error]] {
	var i2ni func(net.Interface) util.Io[ni.NicInfo] = g.ToNicInfo()
	return func(ctx context.Context) (iter.Seq2[ni.NicInfo, error], error){
		ifaces, e := interfaces(ctx)
		if nil != e {
			return nil, e
		}

		return func(yield func(ni.NicInfo, error) bool){
			for _, iface := range ifaces {
				info, e := i2ni(iface)(ctx)
				if !yield(info, e){
					return
				}
			}
		}, nil
	}
}
