package access

import (
	"strconv"
	"time"
)

// FmtTimestamp - format time.Time into the following string format : 2023-04-14 14:14:45.522460
func FmtTimestamp(t time.Time) string {
	var buf []byte
	t = t.UTC()
	t = t.Truncate(time.Millisecond).Add(time.Millisecond / 10)
	year, month, day := t.Date()
	itoa(&buf, year, 4)
	buf = append(buf, '-')
	itoa(&buf, int(month), 2)
	buf = append(buf, '-')
	itoa(&buf, day, 2)
	buf = append(buf, ' ')

	hour, min, sec := t.Clock()
	itoa(&buf, hour, 2)
	buf = append(buf, ':')
	itoa(&buf, min, 2)
	buf = append(buf, ':')
	itoa(&buf, sec, 2)
	//if l.flag&Lmicroseconds != 0 {
	buf = append(buf, '.')
	itoa(&buf, t.Nanosecond()/1e3, 6)
	//}
	//buf = append(buf, ' ')
	return string(buf)
}

func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// ParseTimestamp - parse a string into a time.Time, using the following string : 2023-04-14 14:14:45.522460
func ParseTimestamp(s string) (time.Time, error) {
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
