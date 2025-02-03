package tracker

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/zeebo/bencode"
)

type Files []BencodeFile

type InfoHash [20]byte

type PieceHashes [][20]byte

type rawInfo struct {
	Info map[string]any `bencode:"info"`
}

type BencodeTorrent struct {
	// Ссылка обрацения к трекеру для получения пиров
	Announce string `bencode:"announce"`
	// Информация торрента
	Info BencodeInfo `bencode:"info"`
	// Хеш информации торрента
	InfoHash InfoHash
}

type BencodeInfo struct {
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	PieceHashes PieceHashes
	Length      int   `bencode:"length"`
	Files       Files `bencode:"files"`
}

type BencodeFile struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}

// Представление InfoHash в Hex кодировке
func (i *InfoHash) Hex() string {
	return hex.EncodeToString(i[:])
}

// Десериализация файла .torrent в структуру BencodeTorrent
func TorrentFromFile(f []byte) (*BencodeTorrent, error) {
	rawInfo := new(rawInfo)
	if err := bencode.DecodeBytes(f, rawInfo); err != nil {
		return nil, err
	}
	rawPieces, ok := rawInfo.Info["pieces"]
	if !ok {
		return nil, fmt.Errorf("pieces is not exists")
	}
	pieces, ok := rawPieces.(string)
	if !ok {
		return nil, fmt.Errorf("rawPieces incorrect type")
	}
	encodeInfo, err := bencode.EncodeBytes(rawInfo.Info)
	if err != nil {
		return nil, err
	}
	torrent := new(BencodeTorrent)
	torrent.InfoHash = sha1.Sum(encodeInfo)
	torrent.Info.PieceHashes = HashesFromPieces([]byte(pieces))
	if err := bencode.DecodeBytes(f, torrent); err != nil {
		return nil, err
	}
	return torrent, nil
}

type BencodeTrackerResp struct {
	// Указывает как часто мы можем делать запрос на сервер для обновления peers списка (sec)
	Interval int `bencode:"interval"`
	// Первые 4 байта - это IP адрес узла, последние 2 байта - порт(uint16 в big-endian кодировке) (Последовательно)
	Peers string `bencode:"peers"`
}

// Чтение пиров в байтовой последовательности [:4]ip [4:6]port
func (t *BencodeTrackerResp) ReadPeers() []byte {
	return []byte(t.Peers)
}

// Десериализация тела ответа от трекера в структуру BencodeTrackerResp
func TrackerRespFromBytes(b []byte) (*BencodeTrackerResp, error) {
	trackerResp := new(BencodeTrackerResp)
	if err := bencode.DecodeBytes(b, trackerResp); err != nil {
		return nil, err
	}
	return trackerResp, nil
}
