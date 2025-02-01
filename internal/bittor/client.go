package bittor

import (
	"gotracker/internal/tracker"
	"os"
	"path/filepath"
)

type Client struct {
	Torrent   Torrent
	Handshake Handshake
}

func ClientFromFile(f []byte) (client *Client, err error) {
	bencodeTorrent, err := tracker.TorrentFromFile(f)
	if err != nil {
		return
	}
	peerId := tracker.GenPeerId()
	param := tracker.BuildTrackerReqParam(
		bencodeTorrent.Announce,
		bencodeTorrent.InfoHash,
		peerId,
		4009,
		bencodeTorrent.Info.Length,
	)
	bencodeTrackerResp, err := tracker.TrackerReq(param)
	if err != nil {
		return
	}
	handshake := NewHandshake(
		bencodeTorrent.InfoHash,
		peerId,
	)
	client = &Client{
		Torrent: Torrent{
			Announce: bencodeTorrent.Announce,
			Info:     bencodeTorrent.Info,
			InfoHash: bencodeTorrent.InfoHash,
			Peers:    PeersFromBytes(bencodeTrackerResp.ReadPeers()),
			Interval: bencodeTrackerResp.Interval,
		},
		Handshake: handshake,
	}
	return
}

func (c *Client) CreateFiles(root string) error {
	if len(c.Torrent.Info.Files) == 0 {
		filePath := filepath.Join(
			root,
			c.Torrent.Info.Name,
		)
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		file.Close()
	}
	for _, f := range c.Torrent.Info.Files {
		filePath := filepath.Join(
			append(
				[]string{root, c.Torrent.Info.Name},
				f.Path...,
			)...,
		)
		dirPath := filepath.Dir(filePath)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}
