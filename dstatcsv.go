package graph

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

func csvReadLine(r *bufio.Reader) ([]string, error) {
	line, _, err := r.ReadLine()
	if err != nil {
		return []string{}, err
	}

	reader := csv.NewReader(bytes.NewBuffer(line))
	return reader.Read()
}

func Parse(r io.Reader) ([]string, []DstatCSVRow, error) {
	csvlines := bufio.NewReader(r)

	log.Printf("debug: parse dstat csv header")
	head1, err := csvReadLine(csvlines)
	if err != nil {
		return nil, nil, err
	}
	if len(head1) < 1 {
		return nil, nil, fmt.Errorf("unknown Dstat 'version' header: %v", head1)
	}
	if strings.HasPrefix(head1[0], "Dstat") != true {
		return nil, nil, fmt.Errorf("unknown Dstat 'version' header: '%s'", head1[0])
	}

	head2, err := csvReadLine(csvlines)
	if err != nil {
		return nil, nil, err
	}
	if len(head2) < 1 {
		return nil, nil, fmt.Errorf("unknown Dstat 'Auther' header: %v", head2)
	}
	if strings.HasPrefix(head2[0], "Author") != true {
		return nil, nil, fmt.Errorf("unknown Dstat 'Auther' header: '%s'", head2[0])
	}

	head3, err := csvReadLine(csvlines)
	if err != nil {
		return nil, nil, err
	}
	if len(head3) < 1 {
		return nil, nil, fmt.Errorf("unknwon Dstat 'Host' header: %v", head3)
	}
	if strings.HasPrefix(head3[0], "Host") != true {
		return nil, nil, fmt.Errorf("unknwon Dstat 'Host' header: '%s'", head3[0])
	}

	head4, err := csvReadLine(csvlines)
	if err != nil {
		return nil, nil, err
	}
	if len(head4) < 1 {
		return nil, nil, fmt.Errorf("unknwon Dstat 'Cmdline' header: %v", head4)
	}
	if strings.HasPrefix(head4[0], "Cmdline") != true {
		return nil, nil, fmt.Errorf("unknwon Dstat 'Cmdline' header: '%s'", head4[0])
	}
	dateHead := head4[len(head4)-2]
	dateValue := head4[len(head4)-1]
	if strings.HasPrefix(dateHead, "Date") != true {
		return nil, nil, fmt.Errorf("unknown Dstat 'Date' header: '%s'", dateHead)
	}
	if dateValue == "" {
		return nil, nil, fmt.Errorf("unknown Dstat 'Date' value: '%s'", dateValue)
	}
	// dstat date value = '05 Aug 2019 14:05:42 JST'
	baseTime, err := time.Parse("02 Jan 2006 15:04:05 MST", dateValue)
	if err != nil {
		return nil, nil, err
	}

	head5, err := csvReadLine(csvlines)
	if err != nil {
		return nil, nil, err
	}
	if len(head5) < 1 {
		return nil, nil, fmt.Errorf("unknwon field header: %v", head5)
	}
	log.Printf("debug: dstat fields: %v", head5)

	lines := make([]DstatCSVRow, 0)
	reader := csv.NewReader(csvlines)
	columns, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}
	if len(columns) < 1 {
		return nil, nil, fmt.Errorf("unknown Dstat columns header: %v", columns)
	}
	for {
		values, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, err
		}

		m := make(DstatRecord)
		t := ""
		for i, k := range columns {
			if k == "time" {
				t = values[i]
				continue // lazy
			}
			floatValue, e := strconv.ParseFloat(values[i], 64)
			if e != nil {
				log.Printf("error: parse fload error: %s", values[i])
				return nil, nil, e
			}
			m[k] = floatValue
		}
		// dstat time value: '05-08 14:07:32'
		recordTime, err := time.Parse("02-01 15:04:05", t)
		if err != nil {
			log.Printf("error: cant parse time value '%s'", t)
			return nil, nil, err
		}
		formatTime := time.Date(
			baseTime.Year(),
			recordTime.Month(),
			recordTime.Day(),
			recordTime.Hour(),
			recordTime.Minute(),
			recordTime.Second(),
			0,
			baseTime.Location(),
		)
		lines = append(lines, DstatCSVRow{
			Time:   formatTime,
			Values: m,
		})
	}
	recordColumns := make([]string, 0)
	for _, c := range columns {
		if c == "time" {
			continue // skip
		}
		recordColumns = append(recordColumns, c)
	}
	return recordColumns, lines, nil
}
