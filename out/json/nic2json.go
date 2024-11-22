package nic2json

import (
	"io"
	"context"
	"encoding/json"

	ni "github.com/takanoriyanagitani/go-nic-info"
	util "github.com/takanoriyanagitani/go-nic-info/util"

	o "github.com/takanoriyanagitani/go-nic-info/out"
)

func WriterToOutput(w io.Writer) func(ni.NicInfo) util.Io[util.Void]{
	var enc *json.Encoder = json.NewEncoder(w)
	return func(n ni.NicInfo) util.Io[util.Void] {
		return func(_ context.Context) (util.Void, error){
			return util.Empty, enc.Encode(&n)
		}
	}
}

var WriterToNicOutputJson o.WriterToNicOutput = WriterToOutput
