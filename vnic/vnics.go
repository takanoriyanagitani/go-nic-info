package vnic

import (
	"iter"
	"context"
	"maps"

	util "github.com/takanoriyanagitani/go-nic-info/util"
	it "github.com/takanoriyanagitani/go-nic-info/util/iter"
)

type VnicNameSource util.Io[iter.Seq[string]]

func (s VnicNameSource) ToSet() util.Io[map[string]struct{}] {
	return func(ctx context.Context)(map[string]struct{}, error){
		vnics, e := s(ctx)
		if nil != e {
			return nil, e
		}

		var pairs iter.Seq2[string, struct{}] = it.ToSeq2(
			vnics,
			func(_ string) struct{}{ return struct{}{} },
		)

		return maps.Collect(pairs), nil
	}
}
