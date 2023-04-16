//go:build integration

package main

import (
	"github.com/polygon-io/client-go/rest/models"
	"testing"
	"time"
)

func Test_getTimestampsForTradingDayImpl(t *testing.T) {
	openTimestamp := models.Agg{
		Timestamp: models.Millis(time.Date(2023, time.Month(4), 14, 13, 30, 0, 0, time.UTC)),
	}
	closeTimestamp := models.Agg{
		Timestamp: models.Millis(time.Date(2023, time.Month(4), 14, 20, 0, 0, 0, time.UTC)),
	}
	input := &models.GetAggsResponse{
		Results: []models.Agg{openTimestamp, closeTimestamp},
	}
	got := getTimestampsForTradingDayImpl(input)
	if got[0].Timestamp != "1681479000000" {
		t.Errorf("got %v, expected %v", got[0].Timestamp, "1681479000000")
	}
	if got[0].Readable != "Friday April 14 2023 9:30:00 AM" {
		t.Errorf("got %v, expected %v", got[0].Readable, "Friday April 14 2023 9:30:00 AM")
	}
	if got[1].Timestamp != "1681502400000" {
		t.Errorf("got %v, expected %v", got[1].Timestamp, "1681502400000")
	}
	if got[1].Readable != "Friday April 14 2023 4:00:00 PM" {
		t.Errorf("got %v, expected %v", got[1].Readable, "Friday April 14 2023 4:00:00 PM")
	}
}
