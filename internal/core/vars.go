package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	anaclog "github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/iplist"
	"github.com/anacrolix/torrent/mse"
	"golang.org/x/time/rate"
)

var (
	Engine     Eng
	Info       = log.New(os.Stderr, "[INFO] ", log.LstdFlags) // Info Logger
	Warn       = log.New(os.Stderr, "[WARN] ", log.LstdFlags) // Logger for Warnings
	Err        = log.New(os.Stderr, "[ERR ] ", log.LstdFlags) // Error Logger
	Flagconfig = struct {                                     // Configuration for HTTP Handlers
		ListenAddress string
		UnixSocket    string
		TLSKeyPath    string
		TLSCertPath   string
	}{}
	Dirconfig = struct {
		DirPath   string
		ConfigDir string
		CacheDir  string
		DataDir   string
		TrntDir   string
	}{}
	Configmu sync.Mutex
)

type TorConfig struct {
	ListenHost                        *string
	ListenPort                        *int
	NoDefaultPortForwarding           *bool
	UpnpID                            *string
	DisableTrackers                   *bool
	DisablePEX                        *bool
	NoDHT                             *bool
	PeriodicallyAnnounceTorrentsToDht *bool
	NoUpload                          *bool
	DisableAggressiveUpload           *bool
	Seed                              *bool
	UploadLimiterLimit                *float64
	UploadLimiterBurst                *int
	DownloadLimiterLimit              *float64
	DownloadLimiterBurst              *int
	MaxUnverifiedBytes                *int64
	PeerID                            *string
	DisableUTP                        *bool
	DisableTCP                        *bool
	HeaderObfuscationPolicy           *string
	CryptoProvides                    *uint32
	IPBlocklist                       *bool
	DisableIPv6                       *bool
	DisableIPv4                       *bool
	DisableIPv4Peers                  *bool
	Debug                             *bool
	Logger                            *bool
	HTTPUserAgent                     *string
	ExtendedHandshakeClientVersion    *string
	Bep20                             *string
	NominalDialTimeout                *int64
	MinDialTimeout                    *int64
	EstablishedConnsPerTorrent        *int
	HalfOpenConnsPerTorrent           *int
	TotalHalfOpenConns                *int
	TorrentPeersHighWater             *int
	TorrentPeersLowWater              *int
	HandshakesTimeout                 *int64
	PublicIP4                         *string
	PublicIP6                         *string
	DisableAcceptRateLimiting         *bool
	DropDuplicatePeerIds              *bool
	DropMutuallyCompletePeers         *bool
	AcceptPeerConnections             *bool
	DisableWebtorrent                 *bool
	DisableWebseeds                   *bool
}

func (t *TorConfig) ToTorrentConfig() (tc *torrent.ClientConfig) {
	tc = torrent.NewDefaultClientConfig()

	tc.Logger = anaclog.Logger{}
	tc.Logger.Handlers = []anaclog.Handler{anaclog.DiscardHandler} // Discard Logging of Torrent Client by Default

	tc.HTTPProxy = http.ProxyFromEnvironment // Use Proxy Variables from Environment

	if t.ListenHost != nil {
		tc.ListenHost = func(string) string { return *t.ListenHost }
		Info.Println("ListenHost of Torrent Client has been set to ", t.ListenHost)
	}

	if t.ListenPort != nil {
		tc.ListenPort = *t.ListenPort
		Info.Println("ListenPort of Torrent Client has been set to ", tc.ListenPort)
	}

	if t.NoDefaultPortForwarding != nil {
		tc.NoDefaultPortForwarding = *t.NoDefaultPortForwarding
		Info.Println("NoDefaultPortForwarding of Torrent Client has been set to ", tc.NoDefaultPortForwarding)
	}

	if t.UpnpID != nil {
		tc.UpnpID = *t.UpnpID
		Info.Println("UpnpID of Torrent Client has been set to ", tc.UpnpID)
	}

	if t.DisableTrackers != nil {
		tc.DisableTrackers = *t.DisableTrackers
		Info.Println("DisableTrackers of Torrent Client has been set to ", tc.DisableTrackers)
	}

	if t.DisablePEX != nil {
		tc.DisablePEX = *t.DisablePEX
		Info.Println("DisablePEX of Torrent Client has been set to ", tc.DisablePEX)
	}

	if t.NoDHT != nil {
		tc.NoDHT = *t.NoDHT
		Info.Println("NoDHT of Torrent Client has been set to ", tc.NoDHT)
	}

	if t.PeriodicallyAnnounceTorrentsToDht != nil {
		tc.PeriodicallyAnnounceTorrentsToDht = *t.PeriodicallyAnnounceTorrentsToDht
		Info.Println("PeriodicallyAnnounceTorrentsToDht has been set to", tc.PeriodicallyAnnounceTorrentsToDht)
	}

	if t.NoUpload != nil {
		tc.NoUpload = *t.NoUpload
		Info.Println("NoUpload of Torrent Client has been set to ", tc.NoUpload)
	}

	if t.DisableAggressiveUpload != nil {
		tc.DisableAggressiveUpload = *t.DisableAggressiveUpload
		Info.Println("DisableAggressiveUpload of Torrent Client has been set to ", tc.DisableAggressiveUpload)
	}

	if t.Seed != nil {
		tc.Seed = *t.Seed
		Info.Println("Seed of Torrent Client has been set to ", tc.Seed)
	}

	if t.UploadLimiterLimit != nil && t.UploadLimiterBurst != nil {
		tc.UploadRateLimiter = rate.NewLimiter(rate.Limit(*t.UploadLimiterLimit), *t.UploadLimiterBurst)
		Info.Println("Upload Rate Limiter is now Active with  ", t.UploadLimiterLimit, " as limit  and ", t.UploadLimiterBurst, " as Burst")
	}

	if t.DownloadLimiterLimit != nil && t.DownloadLimiterBurst != nil {
		tc.DownloadRateLimiter = rate.NewLimiter(rate.Limit(*t.DownloadLimiterLimit), *t.DownloadLimiterBurst)
		Info.Println("Download Rate Limiter is now Active with  ", t.DownloadLimiterLimit, " as limit  and ", t.DownloadLimiterBurst, " as Burst")
	}

	if t.MaxUnverifiedBytes != nil {
		tc.MaxUnverifiedBytes = *t.MaxUnverifiedBytes
		Info.Println("MaxUnverifiedBytes of Torrent Client has been set to ", tc.MaxUnverifiedBytes)
	}

	if t.PeerID != nil {
		tc.PeerID = *t.PeerID
		Info.Println("PeerID of Torrent Client has been set to ", tc.PeerID)
	}

	if t.DisableUTP != nil {
		tc.DisableUTP = *t.DisableUTP
		Info.Println("DisableUTP of Torrent Client has been set to ", tc.DisableUTP)
	}

	if t.DisableTCP != nil {
		tc.DisableTCP = *t.DisableTCP
		Info.Println("DisableTCP of Torrent Client has been set to ", tc.DisableTCP)
	}

	if t.HeaderObfuscationPolicy != nil {
		if *t.HeaderObfuscationPolicy == "notpreferred" {
			tc.HeaderObfuscationPolicy = torrent.HeaderObfuscationPolicy{
				RequirePreferred: false,
				Preferred:        false,
			}
		} else if *t.HeaderObfuscationPolicy == "preferred" {
			tc.HeaderObfuscationPolicy = torrent.HeaderObfuscationPolicy{
				RequirePreferred: false,
				Preferred:        true,
			}
		} else if *t.HeaderObfuscationPolicy == "requirepreferred" {
			tc.HeaderObfuscationPolicy = torrent.HeaderObfuscationPolicy{
				RequirePreferred: true,
				Preferred:        true,
			}
		}
		Info.Println("HeaderObfuscationPolicy of Torrent Client has been set")
	}

	if t.CryptoProvides != nil {
		if *t.CryptoProvides == 1 {
			tc.CryptoProvides = mse.CryptoMethodPlaintext
		} else if *t.CryptoProvides == 2 {
			tc.CryptoProvides = mse.CryptoMethodRC4
		} else if *t.CryptoProvides == 3 {
			tc.CryptoProvides = mse.AllSupportedCrypto
		}
		Info.Println("CryptoProvides of Torrent Client has been set")
	}

	if t.IPBlocklist != nil {
		if *t.IPBlocklist {
			Info.Println("Trying to Read Torrent Client Blocklist from", filepath.Join(Dirconfig.ConfigDir, "blocklist"))
			blockfile, err := os.Open(filepath.Join(Dirconfig.ConfigDir, "blocklist"))
			if err != nil {
				Err.Println("Please put your Blocklist at", filepath.Join(Dirconfig.ConfigDir, "blocklist"))
				Err.Fatalln("Error Opening Blocklist: ", err)
			}
			defer blockfile.Close()

			// Read blocklist
			tc.IPBlocklist, err = iplist.NewFromReader(blockfile)
			if err != nil {
				Err.Fatalln("Invalid Blocklist: ", err)
			}

			Info.Println("Loading blocklist of ", tc.IPBlocklist.NumRanges(), " Ranges")
		}
	}

	if t.DisableIPv6 != nil {
		tc.DisableIPv6 = *t.DisableIPv6
		Info.Println("DisableIPv6 of Torrent Client has been set to ", tc.DisableIPv6)
	}

	if t.DisableIPv4 != nil {
		tc.DisableIPv4 = *t.DisableIPv4
		Info.Println("DisableIPv4 of Torrent Client has been set to ", tc.DisableIPv4)
	}

	if t.DisableIPv4Peers != nil {
		tc.DisableIPv4Peers = *t.DisableIPv4Peers
		Info.Println("DisableIPv4Peers of Torrent Client has been set to ", tc.DisableIPv4Peers)
	}

	if t.Debug != nil {
		tc.Debug = *t.Debug
		Info.Println("Debug has been set to ", tc.Debug)
	}

	if t.Logger != nil {
		if *t.Logger {
			tc.Logger = anaclog.Logger{}
			tc.Logger.Handlers = []anaclog.Handler{anaclog.StreamHandler{
				W: os.Stderr,
				Fmt: func(msg anaclog.Record) []byte {
					var pc [1]uintptr
					msg.Callers(1, pc[:])
					return []byte(fmt.Sprintf("[TORC] %s %s\n", time.Now().Format("2006/01/02 03:04:05"), msg.Text()))
				},
			}}
		}
	}

	if t.HTTPUserAgent != nil {
		tc.HTTPUserAgent = *t.HTTPUserAgent
		Info.Println("HTTP User Agent of Torrent Client has been set to ", tc.HTTPUserAgent)
	}

	if t.ExtendedHandshakeClientVersion != nil {
		tc.ExtendedHandshakeClientVersion = *t.ExtendedHandshakeClientVersion
		Info.Println("ExtendedHandshakeClientVersion of Torrent Client has been set to ", tc.ExtendedHandshakeClientVersion)
	}
	if t.NominalDialTimeout != nil {
		tc.NominalDialTimeout = time.Duration(*t.NominalDialTimeout)
		Info.Println("NominalDialTimeout of Torrent Client has been set to ", tc.NominalDialTimeout)
	}
	if t.MinDialTimeout != nil {
		tc.MinDialTimeout = time.Duration(*t.MinDialTimeout)
		Info.Println("MinDialTimeout of Torrent Client has been set to ", tc.MinDialTimeout)
	}
	if t.EstablishedConnsPerTorrent != nil {
		tc.EstablishedConnsPerTorrent = *t.EstablishedConnsPerTorrent
		Info.Println("EstablishedConnsPerTorrent of Torrent Client has been set to ", tc.EstablishedConnsPerTorrent)
	}
	if t.HalfOpenConnsPerTorrent != nil {
		tc.HalfOpenConnsPerTorrent = *t.HalfOpenConnsPerTorrent
		Info.Println("HalfOpenConnsPerTorrent of Torrent Client has been set to ", tc.HalfOpenConnsPerTorrent)
	}

	if t.TotalHalfOpenConns != nil {
		tc.TotalHalfOpenConns = *t.TotalHalfOpenConns
		Info.Println("TotalHalfOpenConns of Torrent Client has been set to ", tc.TotalHalfOpenConns)
	}
	if t.TorrentPeersHighWater != nil {
		tc.TorrentPeersHighWater = *t.TorrentPeersHighWater
		Info.Println("TorrentPeersHighWater of Torrent Client has been set to ", tc.TorrentPeersHighWater)
	}
	if t.TorrentPeersLowWater != nil {
		tc.TorrentPeersLowWater = *t.TorrentPeersLowWater
		Info.Println("TorrentPeersLowWater of Torrent Client has been set to ", tc.TorrentPeersLowWater)
	}

	if t.HandshakesTimeout != nil {
		tc.HandshakesTimeout = time.Duration(*t.HandshakesTimeout)
		Info.Println("HandshakesTimeout of Torrent Client has been set to ", tc.HandshakesTimeout)
	}
	if t.PublicIP4 != nil {
		tc.PublicIp4 = net.ParseIP(*t.PublicIP4)
		Info.Println("PublicIpv4 of Torrent Client has been set to ", *t.PublicIP4)
	}

	if t.PublicIP6 != nil {
		tc.PublicIp6 = net.ParseIP(*t.PublicIP6)
		Info.Println("PublicIpv6 of Torrent Client has been set to ", *t.PublicIP6)
	}

	if t.DisableAcceptRateLimiting != nil {
		tc.DisableAcceptRateLimiting = *t.DisableAcceptRateLimiting
		Info.Println("DisableAcceptRateLimiting of Torrent Client has been set to ", tc.DisableAcceptRateLimiting)
	}
	if t.DropDuplicatePeerIds != nil {
		tc.DropDuplicatePeerIds = *t.DropDuplicatePeerIds
		Info.Println("DropDuplicatePeerIds of Torrent Client has been set to ", tc.DropDuplicatePeerIds)
	}
	if t.DropMutuallyCompletePeers != nil {
		tc.DropMutuallyCompletePeers = *t.DropMutuallyCompletePeers
		Info.Println("DropMutuallyCompletePeers of Torrent Client has been set to ", tc.DropMutuallyCompletePeers)
	}
	if t.AcceptPeerConnections != nil {
		tc.AcceptPeerConnections = *t.AcceptPeerConnections
		Info.Println("AcceptPeerConnections of Torrent Client has been set to ", tc.AcceptPeerConnections)
	}

	if t.DisableWebtorrent != nil {
		tc.DisableWebtorrent = *t.DisableWebtorrent
		Info.Println("DisableWebtorrent of Torrent Client has been set to ", tc.DisableWebtorrent)
	}

	if t.DisableWebseeds != nil {
		tc.DisableWebseeds = *t.DisableWebseeds
		Info.Println("DisableWebseeds of Torrent Client has been set to ", tc.DisableWebseeds)
	}
	return
}

// EngConfig is Engine Configuration Structure which doesn't require restart of Torrent Client
type EngConfig struct {
	DisableLocalCache bool   `json:"disableonlinecache"` // Disables Local Torrent Storage
	OnlineCacheURL    string `json:"onlinecacheurl"`     // Default is https://itorrents.org/torrent/%s.torrent , Setting Empty Disables OnlineCache

	TrackerRefresh  int64    `json:"trackerrefreshinterval"` // In Minutes
	TrackerListURLs []string `json:"trackerlisturls"`        // Default List is []string{"https://ngosang.github.io/trackerslist/trackers_best.txt"}

	DisAllowTrackersUser  bool `json:"disallowtrackersforuser"`  // If set to true , Remove all Trackers that is Added by User to magnet/torrent file. Also disallow adding trackers to torrent
	DisAllowTrackersCache bool `json:"disallowtrackersforcache"` // If set to true , Remove all Trackers from Torrent File fetched from Online/Local Cache

	GlobalSeedRatio     float64 `json:"globalseedratio"`        // Stops Torrent on Reaching Provided SeedRatio
	SRRefresh           int64   `json:"seedratiocheckinterval"` // In Minutes
	DontRemoveCacheInfo bool    `json:"dontremovecacheinfo"`    // When Torrent is Deleted from Storage, it's cache file(.torrent) from Local Cache is not Deleted

	LockbyDefault bool `json:"lockbydefault"` // If set to true , locks every torrent on Add

	ListenCompletion bool   `json:"listencompletion"`
	HookPostURL      string `json:"hookposturl"`
	NotifyOnComplete bool   `json:"notifyoncomplete"`
}

func (ec *EngConfig) GetDTU() (ret bool) {
	Configmu.Lock()
	ret = ec.DisAllowTrackersUser
	Configmu.Unlock()
	return
}

func (ec *EngConfig) GetDTC() (ret bool) {
	Configmu.Lock()
	ret = ec.DisAllowTrackersCache
	Configmu.Unlock()
	return
}

func (ec *EngConfig) GetDLC() (ret bool) {
	Configmu.Lock()
	ret = ec.DisableLocalCache
	Configmu.Unlock()
	return
}
func (ec *EngConfig) GetOCU() (ret string) {
	Configmu.Lock()
	ret = ec.OnlineCacheURL
	Configmu.Unlock()
	return
}
func (ec *EngConfig) GetTR() (ret int64) {
	Configmu.Lock()
	ret = ec.TrackerRefresh
	if ret < 0 {
		ret = 0
	}
	Configmu.Unlock()
	return
}
func (ec *EngConfig) GetTLU() (ret []string) {
	Configmu.Lock()
	ret = ec.TrackerListURLs
	Configmu.Unlock()
	return
}

func (ec *EngConfig) GetGSR() (ret float64) {
	Configmu.Lock()
	ret = ec.GlobalSeedRatio
	Configmu.Unlock()
	return
}
func (ec *EngConfig) GetSRR() (ret int64) {
	Configmu.Lock()
	ret = ec.SRRefresh
	if ret < 0 {
		ret = 0
	}
	Configmu.Unlock()
	return
}

func (ec *EngConfig) GetLBD() (ret bool) {
	Configmu.Lock()
	ret = ec.LockbyDefault
	Configmu.Unlock()
	return
}

func (ec *EngConfig) GetListenC() (ret bool) {
	Configmu.Lock()
	ret = ec.ListenCompletion
	Configmu.Unlock()
	return
}

func (ec *EngConfig) GetHPU() (ret string) {
	Configmu.Lock()
	ret = ec.HookPostURL
	Configmu.Unlock()
	return
}

func (ec *EngConfig) GetNOC() (ret bool) {
	Configmu.Lock()
	ret = ec.NotifyOnComplete
	Configmu.Unlock()
	return
}

func (ec *EngConfig) WriteConfig() (err error) {
	_ = os.Remove(filepath.Join(Dirconfig.ConfigDir, "engconfig.json"))
	f, err := os.OpenFile(filepath.Join(Dirconfig.ConfigDir, "engconfig.json"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	jfile, _ := json.MarshalIndent(ec, "", "\t")
	_, _ = f.Write(jfile)
	if err = f.Close(); err != nil {
		return
	}
	return
}

func (ec *EngConfig) DRCI() (ret bool) {
	Configmu.Lock()
	ret = ec.DontRemoveCacheInfo
	Configmu.Unlock()
	return
}

type ConReq struct {
	Command string `json:"command"`
	Data1   string `json:"data1"`
	Data2   string `json:"data2"`
	Data3   string `json:"data3"`
	Aop     int    `json:"aop"`
}

type Resp struct {
	Type     string `json:"type"`
	State    string `json:"state"`
	Infohash string `json:"infohash,omitempty"`
	Msg      string `json:"message"`
}

type DataMsg struct {
	Type     string      `json:"type"`
	Data     interface{} `json:"data,omitempty"`
	Infohash string      `json:"infohash,omitempty"`
}

type ConnectionMsg struct {
	Type    string `json:"usertype"`
	Session string `json:"session"`
}

type DiskUsageStat struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

type UserConnMsg struct {
	Username string    `json:"username"`
	IsAdmin  bool      `json:"isadmin"`
	Time     time.Time `json:"contime"`
}

type Mutbool struct {
	sync.Mutex
	val bool
}

func (M *Mutbool) Set(v bool) {
	M.Lock()
	M.val = v
	M.Unlock()
}

func (M *Mutbool) Get() (ret bool) {
	M.Lock()
	ret = M.val
	M.Unlock()
	return
}

type MutInt struct {
	sync.Mutex
	val int
}

func (M *MutInt) Set(v int) {
	M.Lock()
	M.val = v
	M.Unlock()
}

func (M *MutInt) Get() (ret int) {
	M.Lock()
	ret = M.val
	M.Unlock()
	return
}

func (M *MutInt) Inc() {
	M.Lock()
	M.val++
	M.Unlock()
}

func (M *MutInt) Dec() {
	M.Lock()
	M.val--
	M.Unlock()
}
