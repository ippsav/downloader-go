package helpers

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Download struct {
	URL      string
	Segments [][2]int64
}

func NewDownload(url string, segmentsCount int) *Download {
	segments := make([][2]int64, segmentsCount)
	return &Download{
		URL:      url,
		Segments: segments,
	}
}

func (d *Download) Do(file *os.File) error {
	wg := sync.WaitGroup{}
	downloadSize, err := d.getDownloadSize()
	if err != nil {
		return err
	}
	d.initSegments(downloadSize)
	wg.Add(len(d.Segments))
	for i, segment := range d.Segments {
		go func(seg [2]int64, index int) {
			defer wg.Done()
			fmt.Printf("download segment %v\n", seg)
			req, _ := http.NewRequest("GET", d.URL, nil)
			//if err != nil {
			//	return errors.Wrap(err, "could not create request")
			//}
			rangeBytes := fmt.Sprintf("bytes=%d,%d", seg[0], seg[1])
			req.Header.Set("Range", rangeBytes)
			resp, _ := http.DefaultClient.Do(req)
			//if err != nil {
			//	return errors.Wrapf(err, "could not fetch the current fragment %d", i)
			//}
			bodyContent, _ := ioutil.ReadAll(resp.Body)
			//if err != nil {
			//	return errors.Wrap(err, "could not read from body ")
			//}
			tmpFile, _ := os.OpenFile(fmt.Sprintf("tmp/frag-%d.tmp", index), os.O_WRONLY|os.O_CREATE, os.ModePerm)
			//if err != nil {
			//	return errors.Wrap(err, "could not create temporary file")
			//}
			_, _ = tmpFile.Write(bodyContent)
			//if err != nil{
			//	return errors.Wrap(err,"could not write body content")
			//}
			tmpFile.Close()
		}(segment, i)
	}
	wg.Wait()
	for i, _ := range d.Segments {
		tmpFile, err := os.OpenFile(fmt.Sprintf("tmp/frag-%d.tmp", i), os.O_RDONLY, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "could not open tmp file to read")
		}
		tmpContent, _ := ioutil.ReadAll(tmpFile)
		file.Write(tmpContent)
		tmpFile.Close()
	}
	return nil
}

func (d *Download) initSegments(size int64) {
	fragmentSize := size / int64(len(d.Segments))
	for i, _ := range d.Segments {
		if i == 0 {
			d.Segments[i][0] = 0
		} else {
			d.Segments[i][0] = d.Segments[i-1][1] + 1
		}

		if i < len(d.Segments)-1 {
			d.Segments[i][1] = d.Segments[i][0] + fragmentSize
		} else {
			d.Segments[i][1] = size - 1
		}
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
