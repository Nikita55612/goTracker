package gotracker

import (
	"encoding/json"
	"fmt"
	"gotracker/config"
	"gotracker/internal/bittor"
	"gotracker/internal/tracker"
	"os"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestConfig(t *testing.T) {
	var cfg config.Config
	data, err := os.ReadFile("Config.toml")
	if err != nil {
		t.Fatal(err)
	}
	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", cfg)
}

func TestGetConfig(t *testing.T) {
	t.Logf("%#v\n", config.Get())
}

func TestTrackerSearchQuery(t *testing.T) {
	tracker.SearchQuery("Video", 0)
	t.Log("TestTrackerSearchQuery is done")
}

func TestItemTorFile(t *testing.T) {
	_, err := tracker.ItemTorrentFile(6455121)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("TestItemTorFile is done")
}

func TestBencode(t *testing.T) {
	items, err := tracker.SearchQuery("Book", 0)
	if err != nil {
		t.Fatal(err)
	}
	itemId := (*items)[0].Id
	file, err := tracker.ItemTorrentFile(*itemId)
	if err != nil {
		t.Fatal(err)
	}
	tor, err := tracker.TorrentFromFile(file)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("InfoHash:", tor.InfoHash)
	t.Log("TestBencode is done")
}

func TestPeerId(t *testing.T) {
	peerId := tracker.GenPeerId()
	fmt.Println("InfoHash:", string(peerId[:]))
	t.Log("TestBencode is done")
}

func TestTrackerReq(t *testing.T) {
	items, err := tracker.SearchQuery("root", 0)
	if err != nil {
		t.Fatal(err)
	}
	itemId := (*items)[0].Id
	file, err := tracker.ItemTorrentFile(*itemId)
	if err != nil {
		t.Fatal(err)
	}
	tor, err := tracker.TorrentFromFile(file)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("InfoHash:", tor.InfoHash.Hex())
	peerId := tracker.GenPeerId()
	param := tracker.BuildTrackerReqParam(tor.Announce, tor.InfoHash, peerId, 4405, tor.Info.Length)
	resp, err := tracker.TrackerReq(param)
	if err != nil {
		t.Fatal(err)
	}
	peers := resp.ReadPeers()
	fmt.Println("Interval:", resp.Interval)
	fmt.Println("Peers:", bittor.PeersFromBytes(peers))
	t.Log("TestTrackerReq is done")
}

func TestClient(t *testing.T) {
	items, err := tracker.SearchQuery("Duet", 0)
	if err != nil {
		t.Fatal(err)
	}
	itemId := (*items)[0].Id
	file, err := tracker.ItemTorrentFile(*itemId)
	if err != nil {
		t.Fatal(err)
	}
	client, err := bittor.ClientFromFile(file)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", client)
	fmt.Println(client.Torrent.Info.Length)
	jsonTorrent, err := json.Marshal(client.Torrent)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println("\n\n", string(jsonTorrent))
	os.WriteFile("data.json", jsonTorrent, 0644)
	t.Log("TestClient is done")
}
