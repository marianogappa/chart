package main

import (
	"bufio"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/marianogappa/chart/format"
)

type dataset struct {
	fss      [][]float64
	sss      [][]string
	tss      [][]time.Time
	minFSS   []float64
	maxFSS   []float64
	stdinLen int
}

func (d dataset) Len() int {
	if d.fss == nil {
		return len(d.tss)
	}
	return len(d.fss)
}

func mustNewDataset(r io.Reader, f format.LineFormat) *dataset {
	d, err := newDataset(r, f)
	if err != nil {
		log.WithError(err).Fatal("Could not build dataset.")
	}
	return d
}

func newDataset(r io.Reader, f format.LineFormat) (*dataset, error) {
	d := &dataset{
		fss:    make([][]float64, 0, 500),
		sss:    make([][]string, 0, 500),
		tss:    make([][]time.Time, 0, 500),
		minFSS: make([]float64, 0, 500),
		maxFSS: make([]float64, 0, 500),
	}
	return d, d.read(r, f)
}

func (d *dataset) read(r io.Reader, f format.LineFormat) error {
	var (
		nilSSS, nilFSS, nilTSS = true, true, true
		scanner                = bufio.NewScanner(r)
		stdinLen               = 0
	)
	for scanner.Scan() {
		stdinLen++
		fs, ss, ts, err := f.ParseLine(scanner.Text())
		if err != nil {
			continue
		}
		if nilSSS && len(ss) > 0 {
			nilSSS = false
		}
		if nilFSS && len(fs) > 0 {
			nilFSS = false
		}
		if nilTSS && len(ts) > 0 {
			nilTSS = false
		}

		for i, f := range fs {
			if len(d.minFSS) == i {
				d.minFSS = append(d.minFSS, f)
			}
			if len(d.maxFSS) == i {
				d.maxFSS = append(d.maxFSS, f)
			}
			if f < d.minFSS[i] {
				d.minFSS[i] = f
			}
			if f > d.maxFSS[i] {
				d.maxFSS[i] = f
			}
		}

		d.fss = append(d.fss, fs)
		d.sss = append(d.sss, ss)
		d.tss = append(d.tss, ts)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	d.stdinLen = stdinLen
	if nilSSS {
		d.sss = nil
	}
	if nilFSS {
		d.fss = nil
		d.minFSS = nil
		d.maxFSS = nil
	}
	if nilTSS {
		d.tss = nil
	}

	return nil
}
