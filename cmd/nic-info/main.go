package main

import (
	"net"
	"os"
	"iter"
	"context"
	"log"

	ni "github.com/takanoriyanagitani/go-nic-info"
	util "github.com/takanoriyanagitani/go-nic-info/util"

	nics "github.com/takanoriyanagitani/go-nic-info/nics"
	ai "github.com/takanoriyanagitani/go-nic-info/addr"

	out "github.com/takanoriyanagitani/go-nic-info/out"
	oj "github.com/takanoriyanagitani/go-nic-info/out/json"

	vn "github.com/takanoriyanagitani/go-nic-info/vnic"
	ln "github.com/takanoriyanagitani/go-nic-info/platform/linux/nictype"
)

var vnicNameSource vn.VnicNameSource = ln.VnicNameSourceDefault
var vnicNameSet util.Io[map[string]struct{}] = vnicNameSource.ToSet()

var nicsSource nics.NetworkInterfacesSource = nics.
	NetworkInterfacesSourceDefault
var nicsSourceIo util.Io[[]net.Interface] = util.Io[[]net.Interface](
	nicsSource,
)

var nic2typ util.Io[ni.NetworkInterfaceToNicType] = util.Bind(
	vnicNameSet,
	util.Lift(ni.InterfaceToNicType),
)

var ni2addrs ai.NetworkInterfaceToAddrs = ai.NetworkInterfaceToAddrsDefault

var nigen util.Io[ai.NicInfoGen] = util.Bind(
	nic2typ,
	util.Lift(func(i2t ni.NetworkInterfaceToNicType) ai.NicInfoGen {
		return ai.NicInfoGen{
			NetworkInterfaceToAddrs: ni2addrs,
			NetworkInterfaceToNicType: i2t,
		}
	}),
)

var nicsInfo util.Io[iter.Seq2[ni.NicInfo, error]] = util.Bind(
	nigen,
	func(g ai.NicInfoGen) util.Io[iter.Seq2[ni.NicInfo, error]] {
		return g.InterfacesToNicsInfo(nicsSourceIo)
	},
)

var w2no out.WriterToNicOutput = oj.WriterToNicOutputJson
var nics2stdout func(
	util.Io[iter.Seq2[ni.NicInfo, error]],
) util.Io[iter.Seq[error]] = w2no.ToNicSourceToErrors(os.Stdout)

var outErrors util.Io[iter.Seq[error]] = nics2stdout(nicsInfo)

var errors2error func(iter.Seq[error]) util.Io[util.Void] = func(
	i iter.Seq[error],
) util.Io[util.Void] {
	return func(_ context.Context) (util.Void, error){
		for e := range i {
			if nil != e {
				return util.Empty, e
			}
		}
		return util.Empty, nil
	}
}

var errs2err util.Io[util.Void] = util.Bind(
	outErrors,
	errors2error,
)

func main(){
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, e := errs2err(ctx)
	if nil != e {
		log.Printf("%v\n", e)
	}
	cancel()

	<-ctx.Done()
}
