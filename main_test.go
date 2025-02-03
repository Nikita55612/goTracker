package gotracker

import (
	"fmt"
	"gotracker/config"
	"gotracker/internal/bittor"
	"gotracker/internal/tracker"
	"net"
	"os"
	"testing"
	"time"

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
	items, err := tracker.SearchQuery("Red Hot Chili Peppers", 0)
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
	items, err := tracker.SearchQuery("Red Hot Chili Peppers", 0)
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
		t.Fatalf("%v", err)
	}
	// err = client.CreateFiles(`./`)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// jsonTorrent, err := json.Marshal(client.Torrent)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println("\n\n", string(jsonTorrent))
	// os.WriteFile("data.json", jsonTorrent, 0644)
	fmt.Println("Pieces Length:", len(client.Torrent.Info.PieceHashes))
	fmt.Println("PieceLength:", client.Torrent.Info.PieceLength)
	t.Log("TestClient is done")
}

func TestPeerConn(t *testing.T) {
	items, err := tracker.SearchQuery("Red Hot Chili Peppers", 0)
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

	for _, peer := range client.Torrent.Peers {
		fmt.Printf("Connecting to peer: %s\n", peer.String())

		conn, err := net.DialTimeout("tcp", peer.String(), 5*time.Second)
		if err != nil {
			fmt.Printf("Failed to connect: %v\n", err)
			continue
		}
		defer conn.Close()

		// Формируем handshake-сообщение
		handshake := client.Handshake

		fmt.Printf("Sending handshake: %x\n", handshake)

		// Отправляем handshake
		_, err = conn.Write(handshake[:])
		if err != nil {
			fmt.Printf("Failed to write to connection: %v\n", err)
			continue
		}

		// Устанавливаем тайм-аут на чтение
		err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			fmt.Printf("Failed to set read deadline: %v\n", err)
			continue
		}

		// Читаем ответ
		buffer := make([]byte, 4096)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Failed to read from connection: %v\n", err)
			continue
		}

		fmt.Printf("Received %d bytes: %x\n", n, buffer[:n])
	}

	t.Log("TestPeerConn is done")
}

//13426974546f7272656e742070726f746f636f6c0000000000000000211a71cd378eb7901fa43eca2600ee5e8bf992da2d5554323231302d365a4437353543306e553733
// 13426974546f7272656e742070726f746f636f6c0000000000000000211a71cd378eb7901fa43eca2600ee5e8bf992da2d5554323231302d775a4d6c5976526f5a347979
