package dateutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func notError(input string) bool {
	v, err := ParseGetUTC(input)
	ts := ISO8601LongMsec(*v)
	fmt.Println(ts)
	return err == nil
}

func TestConfig(t *testing.T) {
	// RFC7232
	assert.True(t, notError("Tue, 03 Jan 2006 15:04:05 GMT"))
	// seconds timestamp
	assert.True(t, notError("1611726543"))
	// nanoseconds timestamp
	assert.True(t, notError("1611726543379050000"))

	assert.True(t, notError("19680526T073000-0500"))
	assert.True(t, notError("19680526T073000.567-0500"))
	assert.True(t, notError("19680526T073000.567890-0500"))
	assert.True(t, notError("19680526T073000.012345678-0500"))

	assert.True(t, notError("1968-01-02T03:04:05-05:00"))
	assert.True(t, notError("1968-01-02T03:04:05.567-05:00"))
	assert.True(t, notError("1968-01-02T03:04:05.567890-05:00"))
	assert.True(t, notError("1968-01-02T03:04:05.012345678-05:00"))

	assert.True(t, notError("19680526T073000Z"))
	assert.True(t, notError("19680526T073000.123Z"))
	assert.True(t, notError("19680526T073000.123456Z"))
	assert.True(t, notError("19680526T073000.012345678Z"))

	assert.True(t, notError("1968-01-02T03:04:05Z"))
	assert.True(t, notError("1968-01-02T03:04:05.123Z"))
	assert.True(t, notError("1968-01-02T03:04:05.123456Z"))
	assert.True(t, notError("1968-01-02T03:04:05.012345678Z"))

	assert.True(t, notError("19680526T073000"))
	assert.True(t, notError("19680526T073000.123"))
	assert.True(t, notError("19680526T073000.123456"))
	assert.True(t, notError("19680526T073000.012345678"))

	assert.True(t, notError("19680526"))
	assert.True(t, notError("1968-01-02"))
	assert.True(t, notError("1968/01/02"))
	assert.True(t, notError("01/02/1968"))
	assert.True(t, notError("1/2/1968"))

	assert.True(t, notError("1968-01-02T03-04-05Z"))
	assert.True(t, notError("1968-01-02T03-04-05.123Z"))
	assert.True(t, notError("1968-01-02T03-04-05.123456Z"))
	assert.True(t, notError("1968-01-02T03-04-05.012345678Z"))

	assert.True(t, notError("1968-01-02T03-04-05-0500"))
	assert.True(t, notError("1968-01-02T03-04-05.567-0500"))
	assert.True(t, notError("1968-01-02T03-04-05.567890-0500"))
	assert.True(t, notError("1968-01-02T03-04-05.012345678-0500"))

	assert.True(t, notError("1968-01-02T03-04-05-05:00"))
	assert.True(t, notError("1968-01-02T03-04-05.567-05:00"))
	assert.True(t, notError("1968-01-02T03-04-05.567890-05:00"))
	assert.True(t, notError("1968-01-02T03-04-05.012345678-05:00"))

	assert.True(t, notError("20201210T223900-0500"))
	assert.True(t, notError("20200101T000000-0500"))

	assert.False(t, notError("hello"))
}

// TestIsPeriod check to see if
func TestIsPeriod(t *testing.T) {
	isPeriod := IsPeriod("")
	assert.False(t, isPeriod)

	isPeriod = IsPeriod("P3Y4M5DT6H4M3S")
	assert.True(t, isPeriod)

	isPeriod = IsPeriod("P3y4m5dT6h4m3s")
	assert.True(t, isPeriod)

	p, err := Period("P3Y4M5DT6H4M3S")
	assert.Nil(t, err)

	t.Logf("Period %s", p)
}

func TestParsePeriod(t *testing.T) {
	pStr := "PT10M"
	p, err := Period(pStr)
	assert.Nil(t, err)
	assert.Equal(t, p, pStr)
}

func TestPeriodPositve(t *testing.T) {
	p, err := Period("PT10M")
	assert.Equal(t, err, nil)
	t.Logf("Is period %v", IsPeriod(p))

	pp, err := PeriodPositive(p)
	assert.Equal(t, err, nil)

	pn, err := PeriodNegative(p)
	assert.Equal(t, err, nil)

	assert.Equal(t, pp, p)
	assert.Equal(t, pn, "-PT10M")
}

func TestOrdering(t *testing.T) {
	t1, err1 := ParseGetUTC("20201210T223900-0500")
	if err1 != nil {
		t.Logf("Problem getting first date %v", err1)
	}
	t2, err2 := ParseGetUTC("20201211T223900-0500")
	if err2 != nil {
		t.Logf("Problem getting second date %v", err2)
	}
	if err2 != nil {

	}
	assert.True(t, StartIsBeforeEnd(*t1, *t2) == true)
	assert.True(t, StartIsBeforeEnd(*t2, *t1) == false)
}
