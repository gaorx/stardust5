package sdcsv

import (
	"bytes"
	"encoding/csv"
	"github.com/gaorx/stardust5/sdslices"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/gaorx/stardust5/sderr"
)

type Reader struct {
	reader *csv.Reader
	fields []string
}

type Options struct {
	Header           bool
	Fields           []string
	Comma            rune
	Comment          rune
	FieldsPerRecord  int
	LazyQuotes       bool
	TrimLeadingSpace bool
	ReuseRecord      bool
}

type HandlerResult int

const (
	Stop     HandlerResult = 0
	Continue HandlerResult = 1
)

func NewReader(r io.Reader, opts *Options) (*Reader, error) {
	r1 := csv.NewReader(r)
	var fields []string
	if opts != nil {
		if opts.Comma != 0 {
			r1.Comma = opts.Comma
		}
		if opts.Comma != 0 {
			r1.Comment = opts.Comment
		}
		if opts.FieldsPerRecord > 0 {
			r1.FieldsPerRecord = opts.FieldsPerRecord
		}
		r1.LazyQuotes = opts.LazyQuotes
		r1.TrimLeadingSpace = opts.TrimLeadingSpace
		r1.ReuseRecord = opts.ReuseRecord
	}
	if opts != nil && len(opts.Fields) > 0 {
		fields = slices.Clone(opts.Fields)
	}
	if opts != nil && opts.Header {
		header, err := r1.Read()
		if err != nil {
			return nil, sderr.Wrap(err, "read CSV error")
		}
		header = slices.Clone(header)
		if len(fields) == 0 {
			fields = header
		}
	}
	return &Reader{
		reader: r1,
		fields: fields,
	}, nil
}

func NewReaderBytes(b []byte, opts *Options) (*Reader, error) {
	return NewReader(bytes.NewReader(sdslices.Ensure(b)), opts)
}

func NewReaderText(s string, opts *Options) (*Reader, error) {
	return NewReader(strings.NewReader(s), opts)
}

func NewReaderFile(filename string, opts *Options) (*Reader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, sderr.Wrap(err, "open CSV file error")
	}
	return NewReader(f, opts)
}

func (r *Reader) Fields() []string {
	return r.fields
}

func (r *Reader) ReadRecord() ([]string, error) {
	rec, err := r.reader.Read()
	if err != nil {
		return nil, sderr.Wrap(err, "read CSV record error")
	}
	return rec, nil
}

func (r *Reader) ReadRecords() ([][]string, error) {
	recs, err := r.reader.ReadAll()
	if err != nil {
		return nil, sderr.Wrap(err, "read CSV all record error")
	}
	return recs, nil
}

func (r *Reader) ReadMap() (map[string]string, error) {
	if len(r.fields) <= 0 {
		return nil, sderr.New("no CSV field")
	}
	rec, err := r.reader.Read()
	if err != nil {
		return nil, sderr.Wrap(err, "read CSV record error")
	}
	return makeMap(r.fields, rec), nil
}

func (r *Reader) ReadMaps() ([]map[string]string, error) {
	if len(r.fields) <= 0 {
		return nil, sderr.New("no CSV field")
	}
	recs, err := r.reader.ReadAll()
	if err != nil {
		return nil, sderr.Wrap(err, "read CSV all record error")
	}
	maps := make([]map[string]string, 0, len(recs))
	for _, rec := range recs {
		maps = append(maps, makeMap(r.fields, rec))
	}
	return maps, nil
}

func (r *Reader) ForeachRecord(h func(recNo int, rec []string) HandlerResult) error {
	recNo := 0
	for {
		rec, err := r.ReadRecord()
		if err != nil {
			if sderr.Cause(err) == io.EOF {
				break
			} else {
				return sderr.Wrap(err, "foreach CSV record error")
			}
		}
		hr := h(recNo, rec)
		if hr == Stop {
			break
		}
		recNo++
	}
	return nil
}

func (r *Reader) ForeachMap(h func(recNo int, rec map[string]string) HandlerResult) error {
	if len(r.fields) == 0 {
		return sderr.New("no CSV field")
	}
	recNo := 0
	for {
		rec, err := r.ReadRecord()
		if err != nil {
			if sderr.Cause(err) == io.EOF {
				break
			} else {
				return sderr.Wrap(err, "read CSV record error")
			}
		}
		hr := h(recNo, makeMap(r.fields, rec))
		if hr == Stop {
			break
		}
		recNo++
	}
	return nil
}

func makeMap(fields, record []string) map[string]string {
	fieldNum, valNum := len(fields), len(record)
	m := make(map[string]string, fieldNum)
	for i := 0; i < fieldNum; i++ {
		field := fields[i]
		v := ""
		if i < valNum {
			v = record[i]
		}
		m[field] = v
	}
	return m
}
