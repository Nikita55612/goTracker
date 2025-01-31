package tracker

import "fmt"

type SearchResult []SearchItem

// Параметы найденного элемента
type SearchItem struct {
	Id,
	Seeds,
	Leechs,
	Downloads *int
	Mark,
	Forum,
	Topic,
	Author,
	Size,
	Date *string
}

func (s *SearchItem) String() string {
	return fmt.Sprintf(
		"Id[%v] %v (%v)\nForum: %v, Author: %v, [S%v/L%v], Downloads: %v, Size: %v, Date: %v",
		*s.Id,
		*s.Topic,
		*s.Mark,
		*s.Forum,
		*s.Author,
		*s.Seeds,
		*s.Leechs,
		*s.Downloads,
		*s.Size,
		*s.Date,
	)
}

func (s *SearchItem) TorrentFile() ([]byte, error) {
	return ItemTorrentFile(*s.Id)
}

// Параметры для составление ссылки запроса к трекеру
type TrackerReqParam struct {
	BaseUrl  string
	InfoHash InfoHash
	PeerId   PeerId
	Port     uint16
	Uploaded,
	Downloaded,
	Compact,
	Left int
}

// Параметры запроса к трекеру по умолчанию
func defaultTrackerReqParam() *TrackerReqParam {
	return &TrackerReqParam{
		Compact: 1,
	}
}

type OptionTrackerReqParam func(p *TrackerReqParam) *TrackerReqParam

// Set false Compact param for TrackerReqParam
func TRPWithoutCompact(p *TrackerReqParam) *TrackerReqParam {
	p.Compact = 0
	return p
}

// Set Uploaded param for TrackerReqParam
func TRPWithUploaded(p *TrackerReqParam, u int) OptionTrackerReqParam {
	return func(p *TrackerReqParam) *TrackerReqParam {
		p.Uploaded = u
		return p
	}
}

// Set Downloaded param for TrackerReqParam
func TRPWithDownloaded(p *TrackerReqParam, d int) OptionTrackerReqParam {
	return func(p *TrackerReqParam) *TrackerReqParam {
		p.Downloaded = d
		return p
	}
}

// Составление параметров запроса пиров к трекеру
func BuildTrackerReqParam(
	baseUrl string,
	infoHash InfoHash,
	peerId PeerId,
	port uint16,
	left int,
	f ...OptionTrackerReqParam,
) *TrackerReqParam {
	trackerReqParam := defaultTrackerReqParam()
	trackerReqParam.BaseUrl = baseUrl
	copy(trackerReqParam.InfoHash[:], infoHash[:])
	copy(trackerReqParam.PeerId[:], peerId[:])
	trackerReqParam.Left = left
	trackerReqParam.Port = port
	for _, do := range f {
		do(trackerReqParam)
	}
	return trackerReqParam
}
