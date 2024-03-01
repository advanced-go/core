package runtime

import (
	"strconv"
	"time"
)

// FmtRFC3339Millis - format time
// https://datatracker.ietf.org/doc/html/rfc3339
func FmtRFC3339Millis(t time.Time) string {
	// Format according to time.RFC3339Nano since it is highly optimized,
	// but truncate it to use millisecond resolution.
	// Unfortunately, that format trims trailing 0s, so add 1/10 millisecond
	// to guarantee that there are exactly 4 digits after the period.
	var b []byte
	const prefixLen = len("2006-01-02T15:04:05.000")
	n := len(b)
	t = t.Truncate(time.Millisecond).Add(time.Millisecond / 10)
	b = t.AppendFormat(b, time.RFC3339Nano)
	b = append(b[:n+prefixLen], b[n+prefixLen+1:]...) // drop the 4th digit
	return string(b)
}

// ParseRFC3339Millis - parse a string into a time.Time, using the following string : 2023-04-14T14:14:45.522Z
func ParseRFC3339Millis(ts string) (time.Time, error) {
	if len(ts) == 0 {
		return time.Now().UTC(), nil
	}
	year, month, day, err := parseYMD(ts)
	if err != nil {
		return time.Now().UTC(), err
	}
	hour, min, sec, ms, err2 := parseHMSM(ts)
	if err2 != nil {
		return time.Now().UTC(), err2
	}
	return time.Date(year, time.Month(month), day, hour, min, sec, ms*1e6, time.UTC), nil
}

// ParseTimestamp2 - parse a string into a time.Time, using the following string : 2023-04-14 14:14:45.522460
func ParseTimestamp2(s string) (time.Time, error) {
	if len(s) == 0 {
		return time.Now().UTC(), nil
	}
	year, err := strconv.Atoi(s[0:4])
	if err != nil {
		return time.Now().UTC(), err
	}
	month, er1 := strconv.Atoi(s[5:7])
	if er1 != nil {
		return time.Now().UTC(), er1
	}
	day, er2 := strconv.Atoi(s[8:10])
	if er2 != nil {
		return time.Now().UTC(), er2
	}
	hour, er3 := strconv.Atoi(s[11:13])
	if er3 != nil {
		return time.Now().UTC(), er3
	}
	min, er4 := strconv.Atoi(s[14:16])
	if er4 != nil {
		return time.Now().UTC(), er4
	}
	sec, er5 := strconv.Atoi(s[17:19])
	if er5 != nil {
		return time.Now().UTC(), er5
	}
	ns, er6 := strconv.Atoi(s[20:26])
	if er6 != nil {
		return time.Now().UTC(), er6
	}
	return time.Date(year, time.Month(month), day, hour, min, sec, ns*1000, time.UTC), nil
}

func parseYMD(s string) (y int, m int, d int, err error) {
	y, err = strconv.Atoi(s[0:4])
	if err != nil {
		return
	}
	m, err = strconv.Atoi(s[5:7])
	if err != nil {
		return
	}
	d, err = strconv.Atoi(s[8:10])
	return
}

func parseHMSM(s string) (h int, m int, sec int, ms int, err error) {
	h, err = strconv.Atoi(s[11:13])
	if err != nil {
		return
	}
	m, err = strconv.Atoi(s[14:16])
	if err != nil {
		return
	}
	sec, err = strconv.Atoi(s[17:19])
	if err != nil {
		return
	}
	ms, err = strconv.Atoi(s[20:23])
	return
}
