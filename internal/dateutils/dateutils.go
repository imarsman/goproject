package dateutils

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rickb777/date"
	"github.com/rickb777/date/timespan"

	"github.com/rickb777/date/period"
)

var timeFormats = []string{
	// RFC7232
	"Mon, 02 Jan 2006 15:04:05 GMT",
	// Short date time with numerical zone offsets
	"20060102T150405-0700",
	"20060102T150405.000-0700",
	"20060102T150405.000000-0700",
	"20060102T150405.999999999-0700",

	// Long date time with numerical zone offsets
	"2006-01-02T15:04:05-07:00",
	"2006-01-02T15:04:05.000-07:00",
	"2006-01-02T15:04:05.000000-07:00",
	"2006-01-02T15:04:05.999999999-07:00",

	// Short date time with zulu zone offsets
	"20060102T150405Z",
	"20060102T150405.000Z",
	"20060102T150405.000000Z",
	"20060102T150405.999999999Z",

	// Long date time with zulu zone offsets
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05.000Z",
	"2006-01-02T15:04:05.000000Z",
	"2006-01-02T15:04:05.999999999Z",

	// Just in case, Opta
	// https://sw5feed.xmlteam.com/sportsml/files/2021/01/28/soccer/l.premierleague.com/schedules/sn.opta.c8.sched-20210128220305-results.f1.xml
	"2006-01-02 15-04-05",

	// Short date time with no zone offset. Assume UTC.
	"20060102T150405",
	"20060102T150405.000",
	"20060102T150405.000000",
	"20060102T150405.999999999",

	// Hopefully less likely to be found. Assume UTC.
	// NHL does this for event-start-date values
	"20060102",
	"2006-01-02",
	"2006/01/02",
	"01/02/2006",
	"1/2/2006",

	// Weird ones with improper separators
	"2006-01-02T15-04-05Z",
	"2006-01-02T15-04-05.000Z",
	"2006-01-02T15-04-05.000000Z",
	"2006-01-02T15-04-05.999999999Z",

	// Weird ones with improper separators
	"2006-01-02T15-04-05-0700",
	"2006-01-02T15-04-05.000-0700",
	"2006-01-02T15-04-05.000000-0700",
	"2006-01-02T15-04-05.999999999-0700",

	"2006-01-02T15-04-05-07:00",
	"2006-01-02T15-04-05.000-07:00",
	"2006-01-02T15-04-05.000000-07:00",
	"2006-01-02T15-04-05.999999999-07:00",

	// time.RFC3339,
	"2006-01-02T15:04:05Z07:00",
}

func makeTimestamp(ts int64) int64 {
	return ts / 1e6
}

// TimesForDateRange parse an ISO8601 date range which normally is written like
// yyyy-mm-dd/yyyy-mm-dd, indicating start/end. The returned values
// will be the first nanosecond of the start date and the first nannosecond
// following the end of the range.
func TimesForDateRange(start, end string) (*time.Time, *time.Time, error) {
	d1, err := date.Parse("2006-01-02", start)
	if err != nil {
		return nil, nil, fmt.Errorf("error when parsing start %v", err)
	}
	d2, err := date.Parse("2006-01-02", end)
	if err != nil {
		return nil, nil, fmt.Errorf("error when parsing end %v", err)
	}
	ts := timespan.NewTimeSpan(d1.In(time.UTC), d2.In(time.UTC))
	s := ts.Start()
	e := ts.End()
	return &s, &e, nil
}

// RangeDate returns a date range function over start date to end date inclusive.
// After the end of the range, the range function returns a zero date,
// date.IsZero() is true.
//
// Sample usage:
//
// for rd := dateutils.RangeDate(start, end); ; {
// 	 date := rd()
// 	 if date.IsZero() {
// 	   break
// 	 }
// 	 indicesForDays[getIndexForDate(*date)] = ""
// }
func RangeDate(start, end time.Time) func() *time.Time {
	start = start.In(time.UTC)
	end = end.In(time.UTC)

	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() *time.Time {
		if start.After(end) {
			return &time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)

		return &date
	}
}

// DateZeroTime get date with zero time values
func DateZeroTime(t time.Time) *time.Time {
	t = t.In(time.UTC)

	newT := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return &newT
}

// ParseGetUTC parse time string and return UTC zoned time
func ParseGetUTC(timeStr string) (*time.Time, error) {
	return ParseWithLocation(timeStr, *time.UTC)
}

// parseTimestamp returns seconds and nanoseconds from a timestamp that has the
// format "%d.%09d", time.Unix(), int64(time.Nanosecond()))
// if the incoming nanosecond portion is longer or shorter than 9 digits it is
// converted to nanoseconds.  The expectation is that the seconds and
// seconds will be used to create a time variable.  For example:
//     seconds, nanoseconds, err := ParseTimestamp("1136073600.000000001",0)
//     if err == nil since := time.Unix(seconds, nanoseconds)
// 	   returns seconds as def(aultSeconds) if value == ""
// From https://github.com/moby/moby/blob/master/api/types/time/timestamp.go
// Part of Docker, under Apache licence.
// https://github.com/docker/engine-api/blob/master/LICENSE
// The Apache Licence for this timeestamp parsing code is included with this package.
func parseTimestamp(value string) (int64, int64, error) {
	sa := strings.SplitN(value, ".", 2)
	s, err := strconv.ParseInt(sa[0], 10, 64)
	if err != nil {
		return s, 0, err
	}
	if len(sa) != 2 {
		return s, 0, nil
	}
	n, err := strconv.ParseInt(sa[1], 10, 64)
	if err != nil {
		return s, n, err
	}
	// should already be in nanoseconds but just in case convert n to nanoseconds
	n = int64(float64(n) * math.Pow(float64(10), float64(9-len(sa[1]))))
	return s, n, nil
}

// ParseWithLocation and return time with specific location
func ParseWithLocation(timeStr string, loc time.Location) (*time.Time, error) {
	// Deal with oddball unix timestamp
	match, err := regexp.MatchString("^\\d+$", timeStr)
	if err != nil {
		return &time.Time{}, errors.New("Could not parse time")
	}
	if match == true {
		toSend := timeStr
		// Break it into a format that has a period between second and
		// millisecond portions for the function.
		if len(timeStr) > 10 {
			sec, nsec := timeStr[0:10], timeStr[11:len(timeStr)-1]
			toSend = sec + "." + nsec
		}
		// Get seconds, nanoseconds, and error if there was a problem
		s, n, err := parseTimestamp(toSend)
		if err != nil {
			return &time.Time{}, err
		}
		// If it was a unix seconds timestamp n will be zero. If it was a
		// nanoseconds timestamp there will be a nanoseconds portion that is not
		// zero.
		t := time.Unix(s, n).In(&loc)
		return &t, nil
	}

	// Continue on for non unix timestamp patterns
	for _, format := range timeFormats {
		t, err := time.Parse(format, timeStr)
		if err == nil {
			t = t.In(&loc)

			return &t, nil
		}
	}
	return &time.Time{}, errors.New("Could not parse time")
}

// RFC7232 get format used for http headers
//   e.g. Mon, 02 Jan 2006 15:04:05 GMT
// TimeFormat is the time format to use when generating times in HTTP headers.
// It is like time.RFC1123 but hard-codes GMT as the time zone. The time being
// formatted must be in UTC for Format to generate the correct format. This is
// done in the function before the call to format.
func RFC7232(t time.Time) (formatted string) {
	t = t.In(time.UTC)

	return t.Format(http.TimeFormat)
}

// ISO8601Long get ISO8601Long format string result
func ISO8601Long(t time.Time) (formatted string) {
	t = t.In(time.UTC)

	return t.Format("2006-01-02T15:04:05-07:00")
}

// ISO8601LongMsec get ISO8601Long format string result with msec
func ISO8601LongMsec(t time.Time) (formatted string) {
	t = t.In(time.UTC)

	return t.Format("2006-01-02T15:04:05.000-07:00")
}

// ISO8601Short with no seconds
func ISO8601Short(t time.Time) (formatted string) {
	t = t.In(time.UTC)

	return t.Format("20060102T150405-0700")
}

// ISO8601ShortMsec with no seconds
func ISO8601ShortMsec(t time.Time) (formatted string) {
	t = t.In(time.UTC)

	return t.Format("20060102T150405.000-0700")
}

// IsPeriod is input a valid period format
func IsPeriod(input string) bool {
	// Convert to uppercase as this would be an irritating source of errors
	_, err := period.ParseWithNormalise(strings.ToUpper(input), true)
	return err == nil
}

// Period get period string
func Period(input string) (output string, err error) {
	if IsPeriod(input) == false {
		return "", fmt.Errorf("invalid period %s", input)
	}
	// Convert to uppercase as this would be an irritating source of errors
	p, _ := period.ParseWithNormalise(strings.ToUpper(input), true)

	return p.String(), nil
}

// PeriodFromDuration get period string from time.Duration
func PeriodFromDuration(d time.Duration) (output string) {
	p, _ := period.NewOf(d)
	return p.String()
}

// PeriodPositive get period string
func PeriodPositive(input string) (output string, err error) {
	if IsPeriod(input) == false {
		return "", fmt.Errorf("invalid period %s", input)
	}
	p, _ := period.ParseWithNormalise(strings.ToUpper(input), true)
	if p.IsNegative() {
		p = p.Negate()
	}

	return p.String(), nil
}

// PeriodNegative get period string
func PeriodNegative(input string) (output string, err error) {
	if IsPeriod(input) == false {
		return "", fmt.Errorf("invalid period %s", input)
	}
	p, _ := period.ParseWithNormalise(strings.ToUpper(input), true)
	if p.IsPositive() {
		p = p.Negate()
	}

	return p.String(), nil
}

// PeriodAdd add a period to a time
func PeriodAdd(t time.Time, input string) (*time.Time, error) {
	t = t.In(time.UTC)

	if IsPeriod(input) == false {
		fmt.Printf("%s is not a valid period\n", input)
		return &time.Time{}, fmt.Errorf("invalid period %s", input)
	}

	p, err := period.ParseWithNormalise(strings.ToUpper(input), true)
	if err != nil {
		fmt.Printf("Got error parsing period %v\n", err)
		return &time.Time{}, err
	}

	newTime, _ := p.AddTo(t)
	newTime = newTime.In(time.UTC)
	// fmt.Printf("Input %v Period %s New time %v\n", t, input, newTime)

	return &newTime, nil
}

// PeriodSubtract subtract a period from a time. If incoming period is negative
// it will be handled properly (by making sure it is subtracted.
func PeriodSubtract(t time.Time, input string) (*time.Time, error) {
	if IsPeriod(input) == false {
		return &time.Time{}, fmt.Errorf("Invalid period %s", input)
	}

	p, err := period.ParseWithNormalise(strings.ToUpper(input), true)
	if err != nil {
		// fmt.Print(err)
		return &time.Time{}, err
	}

	newTime, _ := p.Negate().AddTo(t)
	if p.IsNegative() == true {
		newTime, _ = p.AddTo(t)
	}

	newTime = newTime.In(time.UTC)
	return &newTime, nil
}

// StartIsBeforeEnd if time 1 is before time 2 return true, else false
func StartIsBeforeEnd(t1 time.Time, t2 time.Time) bool {
	t1 = t1.In(time.UTC)
	t2 = t2.In(time.UTC)

	return t2.Unix()-t1.Unix() > 0
}

// GetPositivePeriodStringBetween get a positive period value between two times
// Can be negated to get negative period.
func GetPositivePeriodStringBetween(t1 time.Time, t2 time.Time) string {
	t1 = t1.In(time.UTC)
	t2 = t2.In(time.UTC)

	p := period.Between(t1, t2)
	if p.IsNegative() == true {
		p = p.Negate()
	}
	p = period.New(p.Years(), p.Months(), p.Days(), p.Hours(), p.Minutes(), p.Seconds())

	return p.String()
}

// GetNegativePeriodStringBetween get a negative period value between two times
// Can be negated to get negative period.
func GetNegativePeriodStringBetween(t1 time.Time, t2 time.Time) string {
	t1 = t1.In(time.UTC)
	t2 = t2.In(time.UTC)

	p := period.Between(t1, t2)
	if p.IsPositive() == true {
		p = p.Negate()
	}
	p = period.New(p.Years(), p.Months(), p.Days(), p.Hours(), p.Minutes(), p.Seconds())

	return p.String()
}

// GetPeriodStringBetween get a period value between two times
func GetPeriodStringBetween(t1 time.Time, t2 time.Time) string {
	t1 = t1.In(time.UTC)
	t2 = t2.In(time.UTC)

	p := period.Between(t1, t2)
	p = period.New(p.Years(), p.Months(), p.Days(), p.Hours(), p.Minutes(), p.Seconds())

	return p.String()
}

// TimeForDate get time for date string
func TimeForDate(ds string) *time.Time {
	d, err := date.Parse(date.ISO8601, ds)
	if err != nil {
		t := time.Now()
		d := date.New(t.Year(), t.Month(), t.Day())
		t2 := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
		return &t2
	}
	t2 := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	return &t2
}

// IsDate is a date string of format 2006-01-02 valid
func IsDate(ds string) bool {
	_, err := date.Parse("2006-01-02", ds)
	return err == nil
}

// DateForTime get date string for time
func DateForTime(t time.Time) string {
	t = t.In(time.UTC)

	d := date.New(t.Year(), t.Month(), t.Day())
	return d.Format(date.ISO8601)
}

// SameDate get date string for time
func SameDate(base string, t time.Time) bool {
	t = t.In(time.UTC)

	ds := date.New(t.Year(), t.Month(), t.Day()).Format(date.ISO8601)
	return base == ds
}

// Period1LTPeriod2 check whether p1 is less than p2
func Period1LTPeriod2(p1S string, p2S string) (bool, error) {
	p1, err := period.ParseWithNormalise(p1S, true)
	if err != nil {
		return false, fmt.Errorf("problem parsing %s %v", p1S, err)
	}

	p2, err := period.ParseWithNormalise(p2S, true)
	if err != nil {
		return false, fmt.Errorf("problem parsing %s %v", p2S, err)
	}
	check := p1.DurationApprox() < p2.DurationApprox()

	return check, nil
}
