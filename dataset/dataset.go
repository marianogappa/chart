package dataset

import (
	"bufio"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/marianogappa/chart/format"
)

type Dataset struct {
	FSS      [][]float64
	SSS      [][]string
	TSS      [][]time.Time
	MinFSS   []float64
	MaxFSS   []float64
	StdinLen int
}

func (d Dataset) Len() int {
	if d.FSS == nil {
		return len(d.TSS)
	}
	return len(d.FSS)
}

func MustNew(r io.Reader, f format.LineFormat) *Dataset {
	d, err := New(r, f)
	if err != nil {
		log.WithError(err).Fatal("Could not build dataset.")
	}
	return d
}

func New(r io.Reader, f format.LineFormat) (*Dataset, error) {
	d := &Dataset{
		FSS:    make([][]float64, 0, 500),
		SSS:    make([][]string, 0, 500),
		TSS:    make([][]time.Time, 0, 500),
		MinFSS: make([]float64, 0, 500),
		MaxFSS: make([]float64, 0, 500),
	}
	return d, d.read(r, f)
}

func (d *Dataset) read(r io.Reader, f format.LineFormat) error {
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
			if len(d.MinFSS) == i {
				d.MinFSS = append(d.MinFSS, f)
			}
			if len(d.MaxFSS) == i {
				d.MaxFSS = append(d.MaxFSS, f)
			}
			if f < d.MinFSS[i] {
				d.MinFSS[i] = f
			}
			if f > d.MaxFSS[i] {
				d.MaxFSS[i] = f
			}
		}

		d.FSS = append(d.FSS, fs)
		d.SSS = append(d.SSS, ss)
		d.TSS = append(d.TSS, ts)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	d.StdinLen = stdinLen
	if nilSSS {
		d.SSS = nil
	}
	if nilFSS {
		d.FSS = nil
		d.MinFSS = nil
		d.MaxFSS = nil
	}
	if nilTSS {
		d.TSS = nil
	}

	return nil
}
