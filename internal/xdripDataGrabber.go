package internal

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type XDripDataGrabber struct {
	baseURL string
	client  *http.Client
	logger  *zap.Logger
}

type EntriesResponse []struct {
	ID         string    `json:"_id"`
	Sgv        int       `json:"sgv"`
	SysTime    time.Time `json:"sysTime"`
	DateString time.Time `json:"dateString"`
	Device     string    `json:"device"`
	Direction  string    `json:"direction"`
	Noise      int       `json:"noise"`
	Filtered   int       `json:"filtered"`
	Type       string    `json:"type"`
	Unfiltered int       `json:"unfiltered"`
	Date       int64     `json:"date"`
	UtcOffset  int       `json:"utcOffset"`
	Mills      int64     `json:"mills"`
}

func NewXDropDataGrabber(url string) *XDripDataGrabber {
	return &XDripDataGrabber{
		baseURL: url,
		client:  &http.Client{},
		logger:  zap.L(),
	}
}

func MgdlToMmol(input int) float64 {
	return float64(input) / 18.0
}

func (x XDripDataGrabber) Entries(count int) EntriesResponse {
	x.logger.Info("Grabbing entries from xdrip", zap.Int("count", count))
	path := "/entries/sgv?count="

	req, err := http.NewRequest("GET", x.baseURL+path+strconv.Itoa(count), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("accept", "application/json")
	r, err := x.client.Do(req)
	defer r.Body.Close()

	entries := EntriesResponse{}

	err = json.NewDecoder(r.Body).Decode(&entries)
	if err != nil {
		panic(err)
	}

	return entries

}
