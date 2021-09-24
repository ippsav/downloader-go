package helpers

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Download struct {
	URL      string
	Output   string
	Segments [][2]int64
}

func NewDownload(url, output string, segmentsCount int) *Download {
	segments := make([][2]int64, segmentsCount)
	return &Download{
		URL:      url,
		Output:   output,
		Segments: segments,
	}
}

func (d *Download) Do(file os.File) error {
	downloadSize, err := d.getDownloadSize()
	fragmentSize := downloadSize / int64(len(d.Segments))
	if err != nil {
		return err
	}
	for i, _ := range d.Segments {
		if i == 0 {
			d.Segments[i][0] = 0
		} else {
			d.Segments[i][0] = d.Segments[i-1][1] + 1
		}

		if i < len(d.Segments)-1 {
			d.Segments[i][1] = d.Segments[i][0] + fragmentSize
		} else {
			d.Segments[i][1] = downloadSize - 1
		}
	}
	for i, segment := range d.Segments {
		fmt.Printf("download segement %v", segment)
		req, err := http.NewRequest("GET", d.URL, nil)
		if err != nil {
			return errors.Wrap(err, "could not create request")
		}
		rangeBytes := fmt.Sprintf("bytes=%d,%d", segment[0], segment[1])
		req.Header.Set("Range", rangeBytes)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return errors.Wrapf(err, "could not fetch the current fragment %d", i)
		}
		bodyContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "could not read from body ")
		}
		tmpFile, err := os.OpenFile(fmt.Sprintf("tmp/frag-%d.tmp", i), os.O_WRONLY, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "could not create temporary file")
		}
		tmpFile.Write(bodyContent)
		tmpFile.Close()
	}
	for i, _ := range d.Segments {
		tmpFile, _ := os.OpenFile(fmt.Sprintf("tmp/frag-%d.tmp", i), os.O_WRONLY, os.ModePerm)
		tmpContent, _ := ioutil.ReadAll(tmpFile)
		file.Write(tmpContent)
		tmpFile.Close()
	}
	return nil
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
