package core

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
)

func TorrentServe(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if !(len(parts) > 4) {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if Engine.LsDb.IsLocked(parts[3]) {
		u, ut, _, err := authHelper(w, r)
		if err != nil {
			Err.Println(err)
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		if ut != 1 {
			if ut == -1 {
				http.Error(w, "404 Not Found", http.StatusNotFound)
				return
			}
			if !Engine.TUDb.HasUser(u, parts[3]) {
				http.Error(w, "404 Not Found", http.StatusNotFound)
				return
			}
		}
	}
	if val := r.URL.Query().Get("dl"); val != "" {
		file := filepath.Join(Dirconfig.TrntDir, filepath.Join(parts[3:]...))
		if !strings.HasPrefix(file, filepath.Join(Dirconfig.TrntDir, parts[3])) {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		info, err := os.Stat(file)
		if err != nil {
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		} else if info.IsDir() {
			if val == "tar" {
				TarDir(file, w, path.Base(r.URL.Path))
				return
			} else if val == "zip" {
				ZipDir(file, w, path.Base(r.URL.Path))
				return
			}
		}
	}
	http.StripPrefix("/api/torrent/", http.FileServer(http.Dir(Dirconfig.TrntDir))).ServeHTTP(w, r)
}

func StreamFile(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if !(len(parts) > 4) {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if Engine.LsDb.IsLocked(parts[3]) {
		u, ut, _, err := authHelper(w, r)
		if err != nil {
			Err.Println(err)
			http.Error(w, "404 Not Found", http.StatusNotFound)
			return
		}
		if ut != 1 {
			if ut == -1 {
				http.Error(w, "404 Not Found", http.StatusNotFound)
				return
			}
			if !Engine.TUDb.HasUser(u, parts[3]) {
				http.Error(w, "404 Not Found", http.StatusNotFound)
				return
			}
		}
	}
	ih, err := MetafromHex(parts[3])
	if err != nil {
		Warn.Println(err)
		return
	}
	t, ok := Engine.Torc.Torrent(ih)
	if !ok {
		Warn.Println("Error fetching torrent of infohash ", ih, " from the client")
		return
	}

	if t.Info() == nil {
		Warn.Println("Metainfo not yet loaded")
		http.Error(w, "Metainfo not yet loaded", http.StatusNotFound)
		return
	}

	// Get File from Given Torrent
	reqpath := strings.Join(parts[4:], "/")
	var f *torrent.File
	for _, file := range t.Files() {
		if file.Path() == reqpath {
			f = file
			break
		}
	}
	if f == nil {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	filereader := f.NewReader()
	defer filereader.Close()
	filereader.SetReadahead(48 << 20)
	http.ServeContent(w, r, reqpath, time.Time{}, filereader)
}

func TarDir(dirpath string, w http.ResponseWriter, name string) {
	w.Header().Set("Content-Type", "application/x-tar")
	w.Header().Set("Content-disposition", `attachment; filename="`+name+`.tar"`)
	w.WriteHeader(http.StatusOK)
	tw := tar.NewWriter(w)
	defer tw.Close()

	_ = filepath.WalkDir(dirpath, func(p string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, ierr := de.Info()
		if ierr != nil {
			return ierr
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		rel, err := filepath.Rel(dirpath, p)
		if err != nil {
			return err
		}
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()

		h, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		h.Name = rel
		if err := tw.WriteHeader(h); err != nil {
			return err
		}
		n, err := io.Copy(tw, f)
		if info.Size() != n {
			return fmt.Errorf("mismatch of size with %s", rel)
		}
		return err
	})
}

func ZipDir(dirpath string, w http.ResponseWriter, name string) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-disposition", `attachment; filename="`+name+`.zip"`)
	w.WriteHeader(http.StatusOK)
	zw := zip.NewWriter(w)
	defer zw.Close()

	_ = filepath.WalkDir(dirpath, func(p string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, ierr := de.Info()
		if ierr != nil {
			return ierr
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		rel, err := filepath.Rel(dirpath, p)
		if err != nil {
			return err
		}
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()

		h, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		h.Name = rel
		//h.Method = zip.Deflate

		zf, err := zw.CreateHeader(h)
		if err != nil {
			return err
		}

		n, err := io.Copy(zf, f)
		if info.Size() != n {
			return fmt.Errorf("mismatch of size with %s", rel)
		}

		return err

	})

}
