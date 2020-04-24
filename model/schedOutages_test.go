package model

import (
	"testing"

	"gotest.tools/assert"
)

func TestIsValid(t *testing.T) {
	so := ScheduledOutage{}
	assert.ErrorContains(t, so.IsValid(), "Name")
	so.Name = "Test"
	assert.ErrorContains(t, so.IsValid(), "Type")
	so.Type = "wrong"
	assert.ErrorContains(t, so.IsValid(), "Invalid scheduled type")
}

func TestHours(t *testing.T) {
	st := ScheduledTime{
		Begins: "06:00:00",
		Ends:   "07:59:59",
	}
	assert.NilError(t, st.IsValid("daily"))
	// Validate Begin
	st.Begins = "3:01pm"
	assert.ErrorContains(t, st.IsValid("daily"), "Invalid begin hour")
	// Validate End
	st.Begins = "00:00:00"
	st.Ends = "5:01am"
	assert.ErrorContains(t, st.IsValid("daily"), "Invalid end hour")
}

func TestIsSpecificValid(t *testing.T) {
	so := ScheduledOutage{
		Name: "test",
		Type: "specific",
		Times: []ScheduledTime{
			{
				Begins: "01-Jan-2019 00:00:00",
				Ends:   "01-Feb-2019 23:59:59",
			},
		},
	}
	assert.NilError(t, so.IsValid())
}

func TestIsDailyValid(t *testing.T) {
	so := ScheduledOutage{
		Name: "test",
		Type: "daily",
		Times: []ScheduledTime{
			{
				Begins: "06:00:00",
				Ends:   "07:59:59",
			},
		},
	}
	assert.NilError(t, so.IsValid())
	// Force error
	so.Times[0].Day = "1"
	assert.ErrorContains(t, so.IsValid(), "Daily schedule")
}

func TestIsWeeklyValid(t *testing.T) {
	so := ScheduledOutage{
		Name: "test",
		Type: "weekly",
		Times: []ScheduledTime{
			{
				Day:    "sunday",
				Begins: "00:00:00",
				Ends:   "07:59:59",
			},
		},
	}
	assert.NilError(t, so.IsValid())
	// Force error
	so.Times[0].Day = "1"
	assert.ErrorContains(t, so.IsValid(), "Invalid day")
}

func TestIsMonthlyValid(t *testing.T) {
	so := ScheduledOutage{
		Name: "test",
		Type: "monthly",
		Times: []ScheduledTime{
			{
				Day:    "5",
				Begins: "00:00:00",
				Ends:   "07:59:59",
			},
		},
	}
	assert.NilError(t, so.IsValid())
	// Force error
	so.Times[0].Day = "monday"
	assert.ErrorContains(t, so.IsValid(), "Invalid monthly day")
}
