package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

func AddMetaCache(ih metainfo.Hash, mi metainfo.MetaInfo) {
	// create .torrent file
	if w, err := os.Stat(Dirconfig.CacheDir); err == nil && w.IsDir() {
		cacheFilePath := filepath.Join(Dirconfig.CacheDir, fmt.Sprintf("%s.torrent", ih.HexString()))
		// only create the cache file if not exists
		if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
			cf, err := os.Create(cacheFilePath)
			if err == nil {
				werr := mi.Write(cf)
				if werr != nil {
					Warn.Println("failed to create torrent file ", err)
				}
				Info.Println("created torrent cache file", ih.HexString())
			} else {
				Warn.Println("failed to create torrent file ", err)
			}
			_ = cf.Close()
		}
	}

}

func RemMetaCache(ih metainfo.Hash) {
	_ = os.Remove(filepath.Join(Dirconfig.CacheDir, fmt.Sprintf("%s.torrent", ih.HexString())))
}

func GetMetaCache(ih metainfo.Hash) (spec *torrent.TorrentSpec, reterr error) {
	return SpecfromPath(filepath.Join(Dirconfig.CacheDir, fmt.Sprintf("%s.torrent", ih.HexString())))
}

func EmptyMetaCache() {
	_ = os.Remove(Dirconfig.CacheDir)
	checkDir(Dirconfig.CacheDir)
}
