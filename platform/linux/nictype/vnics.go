package nictype

import (
	"io/fs"
	"errors"
	"os"
	"context"
	"iter"
	"slices"

	util "github.com/takanoriyanagitani/go-nic-info/util"

	vn "github.com/takanoriyanagitani/go-nic-info/vnic"
)

const VirtualNicDevicesDirDefault string = "sys/devices/virtual/net"
const FsRootDefault string = "/"

type IsDirEntryVirtualNic func(fs.DirEntry)(vnic bool)

func IsDirentVnic(dirent fs.DirEntry)(vnic bool){ return dirent.IsDir() }

var IsDirEntryVirtualNicDefault IsDirEntryVirtualNic = IsDirentVnic

func DirentsToVnics(
	checker IsDirEntryVirtualNic,
) func(iter.Seq[fs.DirEntry]) iter.Seq[string] {
	return func(dirents iter.Seq[fs.DirEntry]) iter.Seq[string] {
		return func(yield func(string) bool){
			for dirent := range dirents {
				var isVnic bool = checker(dirent)
				if isVnic {
					var nicname string = dirent.Name()
					if !yield(nicname){
						return
					}
				}
			}
		}
	}
}

type DirentsSource util.Io[[]fs.DirEntry]

func (s DirentsSource) ToVnicNames(
	checker IsDirEntryVirtualNic,
) vn.VnicNameSource {
	return func(ctx context.Context)(iter.Seq[string], error){
		dirents, e := s(ctx)
		if nil != e {
			return nil, e
		}
		var idirent iter.Seq[fs.DirEntry] = slices.Values(dirents)
		return DirentsToVnics(checker)(idirent), nil
	}
}

func DirentsSourceFs(dirname string) func(fs.FS) util.Io[[]fs.DirEntry] {
	return func(f fs.FS) util.Io[[]fs.DirEntry] {
		return func(_ context.Context) ([]fs.DirEntry, error){
			dirents, e := fs.ReadDir(f, dirname)
			if errors.Is(e, fs.ErrNotExist){
				return nil, nil
			}
			return dirents, e
		}
	}
}

var DirentsSourceDefault util.Io[[]fs.DirEntry] = DirentsSourceFs(
	VirtualNicDevicesDirDefault,
)(os.DirFS(FsRootDefault))

var VnicNameSourceDefault vn.VnicNameSource = DirentsSource(
	DirentsSourceDefault,
).
	ToVnicNames(IsDirEntryVirtualNicDefault)
