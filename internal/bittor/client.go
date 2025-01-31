package bittor

import (
	"gotracker/internal/tracker"
)

type Client struct {
	Torrent Torrent
}

func ClientFromFile(f []byte) (client *Client, err error) {
	bencodeTorrent, err := tracker.TorrentFromFile(f)
	if err != nil {
		return
	}
	param := tracker.BuildTrackerReqParam(
		bencodeTorrent.Announce,
		bencodeTorrent.InfoHash,
		tracker.GenPeerId(),
		4009,
		bencodeTorrent.Info.Length,
	)
	bencodeTrackerResp, err := tracker.TrackerReq(param)
	if err != nil {
		return
	}
	client = &Client{
		Torrent: Torrent{
			Announce: bencodeTorrent.Announce,
			Info:     bencodeTorrent.Info,
			InfoHash: bencodeTorrent.InfoHash,
			Peers:    PeersFromBytes(bencodeTrackerResp.ReadPeers()),
			Interval: bencodeTrackerResp.Interval,
		},
	}
	return
}
