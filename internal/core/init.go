package core

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/anacrolix/chansync"
	utp "github.com/anacrolix/go-libutp"

	"github.com/varbhat/exatorrent/internal/db"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"

	anaclog "github.com/anacrolix/log"
)

func checkDir(dir string) {
	fi, err := os.Stat(dir)

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			er := os.MkdirAll(dir, 0755)
			if er != nil {
				Err.Fatalln("Error Creating Directory")
			}
			return
		} else {
			Err.Fatalf("Error Stat Directory %s ( %s ) \n", dir, err.Error())
			return
		}
	}

	if fi != nil {
		if !fi.IsDir() {
			Err.Fatalln("Non-Directory File Present")
			return
		}
	} else {
		Err.Fatalln("Error Stat Directory ", dir)

	}
}

func Initialize() {
	var cfilename string
	var torcc bool
	var psql bool
	var engc bool
	var err error
	var auser string
	var pw bool

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "exatorrent is bittorrent client\n\n")

		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])

		flag.VisitAll(func(f *flag.Flag) {
			_, _ = fmt.Fprintf(flag.CommandLine.Output(), " -%-5v   %v\n", f.Name, f.Usage)
		})
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), " -%-5v   %v\n", "help", "<opt>  Print this Help")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\nVersion: %s", Version)
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\nLicense: GPLv3")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\nSource : https://github.com/varbhat/exatorrent\n")
	}

	flag.StringVar(&Flagconfig.ListenAddress, "addr", ":5000", `<addr> Listen Address (Default: ":5000")`)
	flag.StringVar(&Flagconfig.UnixSocket, "unix", "", `<path> Unix Socket Path`)
	flag.StringVar(&Flagconfig.TLSKeyPath, "key", "", "<path> Path to TLS Key (Required for HTTPS)")
	flag.StringVar(&Flagconfig.TLSCertPath, "cert", "", "<path> Path to TLS Certificate (Required for HTTPS)")
	flag.StringVar(&Dirconfig.DirPath, "dir", "exadir", `<path> exatorrent Directory (Default: "exadir")`)
	flag.StringVar(&auser, "admin", "adminuser", `<user> Default admin username (Default Username: "adminuser" and Default Password: "adminpassword")`)
	flag.BoolVar(&pw, "passw", false, `<opt>  Set Default admin password from "EXAPASSWORD" environment variable`)
	flag.BoolVar(&psql, "psql", false, "<opt>  Generate Sample Postgresql Connection URL")
	flag.BoolVar(&engc, "engc", false, "<opt>  Generate Custom Engine Configuration")
	flag.BoolVar(&torcc, "torc", false, "<opt>  Generate Custom Torrent Client Configuration")
	flag.Parse()

	if len(flag.Args()) != 0 {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Invalid Flags Provided: %s\n\n", flag.Args())
		flag.Usage()
		os.Exit(1)
	}

	// Display All Flag Configurations Provided to exatorrent
	if Flagconfig.UnixSocket != "" {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\nUnix Socket Path => %s", Flagconfig.UnixSocket)
	} else {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\nAddress => %s", Flagconfig.ListenAddress)
	}
	if Flagconfig.TLSKeyPath != "" {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\nTLS Key Path => %s", Flagconfig.TLSKeyPath)
	}
	if Flagconfig.TLSCertPath != "" {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\nTLS Certificate Path => %s", Flagconfig.TLSCertPath)
	}
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\nDirectory => %s\n\n", Dirconfig.DirPath)

	// Create Required SubDirectories if not present
	checkDir(Dirconfig.DirPath)
	Dirconfig.ConfigDir = filepath.Join(Dirconfig.DirPath, "config")
	checkDir(Dirconfig.ConfigDir)
	Dirconfig.CacheDir = filepath.Join(Dirconfig.DirPath, "cache")
	checkDir(Dirconfig.CacheDir)
	Dirconfig.DataDir = filepath.Join(Dirconfig.DirPath, "data")
	checkDir(Dirconfig.DataDir)
	Dirconfig.TrntDir = filepath.Join(Dirconfig.DirPath, "torrents")
	checkDir(Dirconfig.TrntDir)

	// Load Torrent Client Configuration
	cfilename = filepath.Join(Dirconfig.ConfigDir, "clientconfig.json")
	_, cfileerr := os.Stat(cfilename)
	if cfileerr == nil {
		var e error
		cf, e := os.Open(cfilename)
		if e != nil {
			Err.Fatalln("Error Opening ", cfilename)
		}
		if cf != nil {
			e = json.NewDecoder(cf).Decode(&Engine.Tconfig)
			if e != nil {
				Err.Fatalln("Error Decoding ", cfilename)
			}
			Info.Println("Torrent Client Configuration is now loaded from ", cfilename)
			torcc = true
			_ = cf.Close()
		}
	} else if os.IsNotExist(cfileerr) && torcc {
		jfile, _ := json.MarshalIndent(Engine.Tconfig, "", "\t")
		_ = os.WriteFile(cfilename, jfile, 0644)
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\nSample Torrent Client Configuration has been written at %s\n", cfilename)
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Please Customize Torrent Client Configuration File %s if required , and restart\n", cfilename)
		os.Exit(0)
	}

	// Load Custom Engine Configuration
	Engine.Econfig = EngConfig{GlobalSeedRatio: 0, OnlineCacheURL: "", SRRefresh: 150, TrackerRefresh: 60, TrackerListURLs: []string{"https://ngosang.github.io/trackerslist/trackers_best.txt"}}
	// You can also add these "https://newtrackon.com/api/stable" , "https://cdn.jsdelivr.net/gh/XIU2/TrackersListCollection@master/best.txt"
	cfilename = filepath.Join(Dirconfig.ConfigDir, "engconfig.json")
	_, cfileerr = os.Stat(cfilename)
	if cfileerr == nil {
		var e error
		cf, e := os.Open(cfilename)
		if e != nil {
			Err.Fatalln("Error Opening ", cfilename)
		} else {
			if cf != nil {
				e = json.NewDecoder(cf).Decode(&Engine.Econfig)
				if e != nil {
					Err.Fatalln("Error Decoding ", cfilename)
				}
				Info.Printf("Engine Configuration %+v is now loaded\n", Engine.Econfig)
				engc = true
				_ = cf.Close()
			}
		}
	} else if os.IsNotExist(cfileerr) && engc {
		_ = Engine.Econfig.WriteConfig()
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\nSample Engine Configuration has been written at %s\n", cfilename)
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Please Customize Engine Configuration File %s if required , and restart\n", cfilename)
		os.Exit(0)
	}

	// Read Postgresql Secret
	cfilename = filepath.Join(Dirconfig.ConfigDir, "psqlconfig.txt")
	_, cfileerr = os.Stat(cfilename)
	if cfileerr == nil {
		pfile, err := os.Open(cfilename)
		if err != nil {
			Err.Fatalln("Error Reading Postgresql Connection URL: ", err)
		}
		defer pfile.Close()

		scanner := bufio.NewScanner(pfile)
		for scanner.Scan() {
			Engine.PsqlCon = scanner.Text()
			psql = true
			Info.Println("Postgresql Connection URL at", cfilename, " has been Read")
			break
		}

		if err := scanner.Err(); err != nil {
			Err.Fatalln("Error Reading Postgresql Connection URL: ", err)
		}
	} else {
		Engine.PsqlCon = os.Getenv("DATABASE_URL")
		if len(Engine.PsqlCon) != 0 {
			psql = true
			_ = os.Unsetenv("DATABASE_URL")
		}
	}
	if os.IsNotExist(cfileerr) && psql {
		_ = os.WriteFile(cfilename, []byte("postgres://username:password@localhost:5432/database_name"), 0644)
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\nSample Postgresql Connection URL has been written at %s\n", cfilename)
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Please Enter your Postgresql Connection URL at File %s , and restart\n", cfilename)
		os.Exit(0)
	}

	tc := Engine.Tconfig.ToTorrentConfig()

	// Set Different Logger for UTP
	utp.Logger = anaclog.Logger{} // Info Logger
	utp.Logger.Handlers = []anaclog.Handler{anaclog.StreamHandler{
		W: os.Stderr,
		Fmt: func(msg anaclog.Record) []byte {
			var pc [1]uintptr
			msg.Callers(1, pc[:])
			return []byte(fmt.Sprintf("[UTp ] %s %s\n", time.Now().Format("2006/01/02 03:04:05"), msg.Text()))
		},
	}}

	if psql {

		Engine.TorDb = &db.PsqlTrntDb{}
		Engine.TorDb.Open(Engine.PsqlCon)

		Engine.FsDb = &db.PsqlFsDb{}
		Engine.FsDb.Open(Engine.PsqlCon)

		Engine.LsDb = &db.PsqlLsDb{}
		Engine.LsDb.Open(Engine.PsqlCon)

		Engine.UDb = &db.PsqlUserDb{}
		Engine.UDb.Open(Engine.PsqlCon)

		Engine.TUDb = &db.PsqlTrntUserDb{}
		Engine.TUDb.Open(Engine.PsqlCon)

		Engine.TrackerDB = &db.PsqlTDb{}
		Engine.TrackerDB.Open(Engine.PsqlCon)

		Engine.PcDb, err = db.NewPsqlPieceCompletion(Engine.PsqlCon)

		if err != nil {
			Err.Fatalln("Unable to Connect to Postgresql for PieceCompletion")
		}

	} else {
		sqliteSetup(tc)
	}

	_, err = os.Stat(filepath.Join(Dirconfig.DataDir, ".adminadded"))
	if errors.Is(err, os.ErrNotExist) {
		if pw {
			Info.Println(`Adding Admin user with username "` + auser + `" and custom password`)
			er := Engine.UDb.Add(auser, os.Getenv("EXAPASSWORD"), 1)
			if er != nil {
				Err.Fatalln("Unable to add admin user to adminless exatorrent instance :", er)
			}
			_, er = os.Create(filepath.Join(Dirconfig.DataDir, ".adminadded"))
			if er != nil {
				Err.Fatalln(er)
			}
		} else {
			Info.Println(`Adding Admin user with username "` + auser + `" and password "adminpassword"`)
			er := Engine.UDb.Add(auser, "adminpassword", 1)
			if er != nil {
				Err.Fatalln("Unable to add admin user to adminless exatorrent instance :", er)
			}
			_, er = os.Create(filepath.Join(Dirconfig.DataDir, ".adminadded"))
			if er != nil {
				Err.Fatalln(er)
			}
		}
	}

	stor := storage.NewFileOpts(storage.NewFileClientOpts{ClientBaseDir: Dirconfig.TrntDir, FilePathMaker: nil, TorrentDirMaker: func(baseDir string, info *metainfo.Info, infoHash metainfo.Hash) string {
		return filepath.Join(baseDir, infoHash.HexString())
	}, PieceCompletion: Engine.PcDb})

	tc.DefaultStorage = stor

	Engine.Torc, err = torrent.NewClient(tc)
	if err != nil {
		Err.Fatalln("Unable to Create Torrent Client ", err)
	} else {
		Info.Println("Torrent Client Created")
	}

	Engine.onCloseMap = make(map[metainfo.Hash]*chansync.Flag)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				Warn.Println(err)
			}
		}()
		var stoperr error
		stopsignal := make(chan os.Signal, 5)
		signal.Notify(stopsignal, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		sig := <-stopsignal
		fmt.Fprintf(os.Stderr, "\n")
		Warn.Println("Caught Signal:", sig)
		Warn.Println("Closing exatorrent")
		Engine.TorDb.Close()
		Engine.TrackerDB.Close()
		Engine.TUDb.Close()
		stoperr = Engine.PcDb.Close() // Close PcDb at the end
		if stoperr != nil {
			Warn.Println("Error Closing PieceCompletion DB ", stoperr)
		}
		Engine.Torc.Close()    // Close the Torrent Client
		stoperr = stor.Close() // Close the storage.ClientImplCloser
		if stoperr != nil {
			Warn.Println("Error Closing Default Storage ", stoperr)
		}
		os.Exit(1)
	}()

	//Recover Torrents from Database
	torlist, err := Engine.TorDb.GetTorrents()
	if err != nil {
		Err.Fatalln("Error Recovering Torrents")
	}
	for _, eachtrnt := range torlist {
		go func(started bool, infohash metainfo.Hash) {
			AddfromSpec("", &torrent.TorrentSpec{InfoHash: infohash}, true, true)
			if started {
				StartTorrent("", infohash, true)
			}
			flist := Engine.FsDb.Get(infohash)
			if started {
				for _, f := range flist {
					StopFile("", infohash, f)
				}
			} else {
				for _, f := range flist {
					StartFile("", infohash, f)
				}
			}
		}(eachtrnt.Started, eachtrnt.Infohash)
	}

	go UpdateTrackers()
	go TorrentRoutine()

}
