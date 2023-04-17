//go:build !integration

package main

import (
	"moat/state"
	"testing"
)

func Test_FillInEmptyPrices_Invalid_MustHaveAnOpeningPrice(t *testing.T) {
	timestamps := make([]state.TimestampInfo, 0)
	prices := make([]state.PriceInfo, 0)

	timestamps = append(timestamps, state.TimestampInfo{Timestamp: "1", Readable: "Jan 1 2000 12:00:00 AM"})
	timestamps = append(timestamps, state.TimestampInfo{Timestamp: "2", Readable: "Jan 2 2000 12:00:00 AM"})
	timestamps = append(timestamps, state.TimestampInfo{Timestamp: "3", Readable: "Jan 3 2000 12:00:00 AM"})
	timestamps = append(timestamps, state.TimestampInfo{Timestamp: "4", Readable: "Jan 4 2000 12:00:00 AM"})
	timestamps = append(timestamps, state.TimestampInfo{Timestamp: "5", Readable: "Jan 5 2000 12:00:00 AM"})
	timestamps = append(timestamps, state.TimestampInfo{Timestamp: "6", Readable: "Jan 6 2000 12:00:00 AM"})

	prices = append(prices, state.PriceInfo{Timestamp: "2", Readable: "Jan 2 2000 12:00:00 AM", Price: 2.0})
	prices = append(prices, state.PriceInfo{Timestamp: "3", Readable: "Jan 3 2000 12:00:00 AM", Price: 3.0})
	prices = append(prices, state.PriceInfo{Timestamp: "4", Readable: "Jan 4 2000 12:00:00 AM", Price: 4.0})

	_, err := fillInEmptyPrices(timestamps, prices)
	if err == nil {
		t.Errorf("expected error but got nil")
	}
}
func Test_FillInEmptyPrices_Two(t *testing.T) {
	timestamps := make([]state.TimestampInfo, 0)
	prices := make([]state.PriceInfo, 0)

	timestamps = append(timestamps, state.TimestampInfo{Timestamp: "1", Readable: "Jan 1 2000 12:00:00 AM"})
	timestamps = append(timestamps, state.TimestampInfo{Timestamp: "2", Readable: "Jan 2 2000 12:00:00 AM"})
	timestamps = append(timestamps, state.TimestampInfo{Timestamp: "3", Readable: "Jan 3 2000 12:00:00 AM"})
	timestamps = append(timestamps, state.TimestampInfo{Timestamp: "4", Readable: "Jan 4 2000 12:00:00 AM"})
	timestamps = append(timestamps, state.TimestampInfo{Timestamp: "5", Readable: "Jan 5 2000 12:00:00 AM"})
	timestamps = append(timestamps, state.TimestampInfo{Timestamp: "6", Readable: "Jan 6 2000 12:00:00 AM"})

	prices = append(prices, state.PriceInfo{Timestamp: "1", Readable: "Jan 1 2000 12:00:00 AM", Price: 1.0})
	prices = append(prices, state.PriceInfo{Timestamp: "2", Readable: "Jan 2 2000 12:00:00 AM", Price: 2.0})
	prices = append(prices, state.PriceInfo{Timestamp: "4", Readable: "Jan 4 2000 12:00:00 AM", Price: 4.0})

	got, err := fillInEmptyPrices(timestamps, prices)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
	if len(got) != 6 {
		t.Errorf("expected 6 but got %v", len(got))
	}
	if got[0].Price != 1.0 {
		t.Errorf("expected 1.0 but got %v", got[0].Price)
	}
	if got[1].Price != 2.0 {
		t.Errorf("expected 2.0 but got %v", got[1].Price)
	}
	if got[2].Price != 2.0 {
		t.Errorf("expected 2.0 but got %v", got[2].Price)
	}
	if got[3].Price != 4.0 {
		t.Errorf("expected 4.0 but got %v", got[3].Price)
	}
	if got[4].Price != 4.0 {
		t.Errorf("expected 4.0 but got %v", got[4].Price)
	}
	if got[5].Price != 4.0 {
		t.Errorf("expected 4.0 but got %v", got[5].Price)
	}
}

func Test_ConvertToEst(t *testing.T) {
	inputs := []string{
		"1681459200000",
		"1681460460000",
		"1681461240000",
		"1681470360000",
		"1681474440000",
		"1681477200000",
		"1681484520000",
		"1681487940000",
		"1681488000000",
		"1681488060000",
		"1681498500000",
		"1681502340000",
		"1681502400000",
		"1681502460000",
		"1681516740000",
	}
	expectedOutputs := []string{
		"Friday April 14 2023 4:00:00 AM",
		"Friday April 14 2023 4:21:00 AM",
		"Friday April 14 2023 4:34:00 AM",
		"Friday April 14 2023 7:06:00 AM",
		"Friday April 14 2023 8:14:00 AM",
		"Friday April 14 2023 9:00:00 AM",
		"Friday April 14 2023 11:02:00 AM",
		"Friday April 14 2023 11:59:00 AM",
		"Friday April 14 2023 12:00:00 PM",
		"Friday April 14 2023 12:01:00 PM",
		"Friday April 14 2023 2:55:00 PM",
		"Friday April 14 2023 3:59:00 PM",
		"Friday April 14 2023 4:00:00 PM",
		"Friday April 14 2023 4:01:00 PM",
		"Friday April 14 2023 7:59:00 PM",
	}
	for i := 0; i < len(inputs); i++ {
		_, got := convertToEst(inputs[i])
		if got != expectedOutputs[i] {
			t.Errorf("got %v, expected %v", got, expectedOutputs[i])
		}
	}
}
