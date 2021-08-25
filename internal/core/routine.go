package core

import (
	"bufio"
	"net/http"
	"strings"
	"time"
)

// UpdateTrackers Updates the Trackers from TrackerURL
func UpdateTrackers() {
	defer func() {
		if err := recover(); err != nil {
			Warn.Println(err)
		}
	}()
	for {
		updatetrackers()
		// Pause UpdateTrackers() every Engine.Econfig.TrackerRefresh Minutes
		time.Sleep(time.Minute * time.Duration(Engine.Econfig.GetTR()))
	}

}

// TorrentRoutine Stops Torrent on Reaching Global SeedRatio
func TorrentRoutine() {
	defer func() {
		if err := recover(); err != nil {
			Warn.Println(err)
		}
	}()
	for range time.Tick(time.Minute * time.Duration(Engine.Econfig.GetSRR())) {
		stoponseedratio(Engine.Econfig.GetGSR())
	}
}

func updatetrackers() {
	for _, url := range Engine.Econfig.GetTLU() {
		resp, err := http.Get(url)
		if err != nil {
			Warn.Println("TrackerList URL: ", url, " is Invalid")
			continue
		}

		// Add Trackers to txtlines string slice
		scanner := bufio.NewScanner(resp.Body)
		scanner.Split(bufio.ScanLines)

		var trackerpertrurl int
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			Engine.TrackerDB.Add(line)
			trackerpertrurl++
		}

		_ = resp.Body.Close()
		Info.Println("Loaded ", trackerpertrurl, " trackers from ", url)
	}

	Info.Println("Loaded ", Engine.TrackerDB.Count(), " trackers in total , eliminating duplicates")

	// Add Trackers to Every Torrents
	trckrs := [][]string{Engine.TrackerDB.Get()}
	en := Engine.Torc
	if en != nil {
		for _, trnt := range en.Torrents() {
			if trnt != nil {
				trnt.AddTrackers(trckrs)
			}
		}
		Info.Println("Added Loaded Trackers to Torrents")
	}
}

func stoponseedratio(sr float64) {
	if sr != 0 {
		if Engine.Torc != nil {
			trnts := Engine.Torc.Torrents()

			for _, trnt := range trnts {
				if trnt == nil {
					continue
				}
				if trnt.Info() == nil {
					continue
				}

				stats := trnt.Stats()

				seedratio := float64(stats.BytesWrittenData.Int64()) / float64(stats.BytesReadData.Int64())

				// stops task on reaching ratio
				if seedratio >= sr {
					Warn.Printf("Torrent %s Stopped on Reaching Global SeedRatio", trnt.InfoHash())
					go StopTorrent("", trnt.InfoHash())
				}
			}
		}
	}
}
