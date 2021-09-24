package helpers

import (
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

type Download struct {
	URL      string
	Output   string
	Segments [][2]int
}

func NewDownload(url, output string, segmentsCount int) *Download {
	segments := make([][2]int, segmentsCount)
	return &Download{
		URL:      url,
		Output:   output,
		Segments: segments,
	}
}

func (d *Download) getDownloadSize() (int64, error) {
	resp, err := d.checkStatus()
	if err != nil {
		return -1, err
	}
	downloadSize, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		errors.Wrap(err, "error in conversion")
	}
	return int64(downloadSize), nil
}

func (d *Download) checkStatus() (*http.Response, error) {
	req, err := http.NewRequest("HEAD", d.URL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not fetch header")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("response status %s:", resp.Status)
	}
	return resp, nil
}
