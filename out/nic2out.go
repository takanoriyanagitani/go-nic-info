package out

import (
	"context"
	"bufio"
	"io"
	"iter"

	ni "github.com/takanoriyanagitani/go-nic-info"
	util "github.com/takanoriyanagitani/go-nic-info/util"
)

type WriterToNicOutput func(io.Writer) func(ni.NicInfo) util.Io[util.Void]

func (o WriterToNicOutput) ToNicSourceToErrors(
	w io.Writer,
) func(util.Io[iter.Seq2[ni.NicInfo, error]]) util.Io[iter.Seq[error]] {
	var bw *bufio.Writer = bufio.NewWriter(w)
	var nicout func(ni.NicInfo) util.Io[util.Void] = o(bw)
	return func(
		inics util.Io[iter.Seq2[ni.NicInfo, error]],
	) util.Io[iter.Seq[error]] {
		return func(ctx context.Context)(iter.Seq[error], error){
			rnics, e := inics(ctx)
			if nil != e {
				return nil, e
			}
			return func(yield func(error) bool){
				defer bw.Flush()

				for nic, e := range rnics {
					if nil == e {
						_, e = nicout(nic)(ctx)
					}

					if !yield(e){
						return
					}
				}
			}, nil
		}
	}
}
