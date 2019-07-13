package events

import (
	"fmt"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestSetTime(t *testing.T) {
	e := Event{}
	dateTxt := "Mon Jan 2 15:04:05 -0700 MST 2006"
	date, _ := time.Parse(dateTxt, dateTxt)
	fmt.Println(date)
	fmt.Println(date.UTC())
	e.SetTime(date)
	fmt.Println(e.Time)
	assert.Equal(t, "Monday, January 2, 2006 10:04:05 PM GMT", e.Time)
}
