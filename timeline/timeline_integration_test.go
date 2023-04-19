package timeline

import (
	"moat/state"
	"testing"
)

var tl = CreateTimelineClient()

func Test_Yesterday(t *testing.T) {
	t.Skip("skipping test until I take out the polygon.io api key dependency")
	inputs := [][]int{
		{12, 26, 2022},
		{4, 17, 2023},
		{4, 18, 2023},
	}

	expectedOutputs := []*state.MktDate{
		state.NewMktDate(12, 23, 2022),
		state.NewMktDate(4, 14, 2023),
		state.NewMktDate(4, 17, 2023),
	}

	for idx, input := range inputs {
		mktDate := state.NewMktDate(input[0], input[1], input[2])
		got, _ := tl.Yesterday(mktDate)
		if !got.Equals(expectedOutputs[idx]) {
			t.Errorf("expected %v but got %v", expectedOutputs[idx], got)
		}
	}
}

func Test_IsMarketOpen_True(t *testing.T) {
	t.Skip("skipping test until I take out the polygon.io api key dependency")
	inputs := [][]int{
		{12, 26, 2022},
		{4, 16, 2023},
		{4, 17, 2023},
		{4, 18, 2023},
	}

	expectedOutputs := []bool{
		false,
		false,
		true,
		true,
	}

	for idx, input := range inputs {
		mktDate := state.NewMktDate(input[0], input[1], input[2])
		got := tl.isMarketOpen(mktDate)
		if got != expectedOutputs[idx] {
			t.Errorf("expected %v but got %v", expectedOutputs[idx], got)
		}
	}
}
