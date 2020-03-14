package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"

	"rdbviewer/shared"
)

func formatFirstSeen(ts *timestamp.Timestamp, format string) string {
	date := shared.PbtimeToTime(ts)

	ret := ""
	cmp := time.Date(2018, time.September, 6, 24, 59, 59, 0, time.UTC)
	if date.Before(cmp) == true {
		ret += "before "
	}
	ret += date.Format(format)
	return ret
}

func formatDuration(duration float64) string {
	d := time.Duration(duration*1000.0) * time.Millisecond
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	var ret []string
	if h > 0 {
		ret = append(ret, fmt.Sprintf("%dh", h))
	}
	if m > 0 {
		ret = append(ret, fmt.Sprintf("%dm", m))
	}
	ret = append(ret, fmt.Sprintf("%ds", s))
	return strings.Join(ret, " ")
}
