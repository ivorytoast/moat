//go:build !integration

package main

import (
	"testing"
)

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
