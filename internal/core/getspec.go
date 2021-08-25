package core

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

// SpecfromURL Returns Torrent Spec from HTTP URL
func SpecfromURL(torrentURL string) (spec *torrent.TorrentSpec, reterr error) {
	// TorrentSpecFromMetaInfo may panic if the info is malformed
	defer func() {
		if r := recover(); r != nil {
			reterr = fmt.Errorf("SpecfromURL: error loading spec from URL")
		}
		reterr = nil
	}()

	Info.Println("Adding Torrent from Torrent URL ", torrentURL)

	torrentURL = strings.TrimSpace(torrentURL)
	resp, reterr := http.Get(torrentURL)
	if reterr != nil {
		return
	}

	// Limit Response
	lr := io.LimitReader(resp.Body, 20971520) // 20MB
	info, reterr := metainfo.Load(lr)
	if reterr != nil {
		_ = resp.Body.Close()
		return
	}
	spec = torrent.TorrentSpecFromMetaInfo(info)
	_ = resp.Body.Close()
	return
}

// SpecfromPath Returns Torrent Spec from File Path
func SpecfromPath(path string) (spec *torrent.TorrentSpec, reterr error) {
	// TorrentSpecFromMetaInfo may panic if the info is malformed
	defer func() {
		if r := recover(); r != nil {
			reterr = fmt.Errorf("SpecfromPath: error loading new torrent from file %s: %v+", path, r)
		}
	}()

	fi, err := os.Stat(path)

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file doesn't exist")
	}

	if fi.IsDir() {
		Err.Println("Directory Present instead of file")
		return nil, fmt.Errorf("directory present")
	}

	Info.Println("Getting Torrent Metainfo from File Path", path)

	info, reterr := metainfo.LoadFromFile(path)
	if reterr != nil {
		return
	}
	spec = torrent.TorrentSpecFromMetaInfo(info)
	return
}

// SpecfromBytes Returns Torrent Spec from Bytes
func SpecfromBytes(trnt []byte) (spec *torrent.TorrentSpec, reterr error) {
	// TorrentSpecFromMetaInfo may panic if the info is malformed
	defer func() {
		if r := recover(); r != nil {
			reterr = fmt.Errorf("SpecfromBytes: error loading new torrent from bytes")
		}
	}()
	Info.Println("Getting Torrent Metainfo from Torrent Bytes")
	info, reterr := metainfo.Load(bytes.NewReader(trnt))
	if reterr != nil {
		return nil, reterr
	}
	spec = torrent.TorrentSpecFromMetaInfo(info)
	return
}

// SpecfromB64String Returns Torrent Spec from Base64 Encoded Torrent File
func SpecfromB64String(trnt string) (spec *torrent.TorrentSpec, reterr error) {
	t, err := base64.StdEncoding.DecodeString(trnt)
	if err != nil {
		return nil, err
	}
	return SpecfromBytes(t)
}

// MetafromHex returns metainfo.Hash from given infohash string
func MetafromHex(infohash string) (h metainfo.Hash, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error parsing string to InfoHash")
		}
	}()

	h = metainfo.NewHashFromHex(infohash)

	return h, nil
}

// RemTrackersSpec removes trackers from torrent.Spec
func RemTrackersSpec(spec *torrent.TorrentSpec) {
	if spec == nil {
		return
	}
	spec.Trackers = [][]string{}
}
