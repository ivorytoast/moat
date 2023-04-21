package scheduler

import (
	"moat/state"
	"testing"
)

func Test_SatisfyAsks(t *testing.T) {
	teamsToDivisions := map[int]string{}

	teamsToDivisions[1] = "one"
	teamsToDivisions[2] = "one"
	teamsToDivisions[3] = "one"
	teamsToDivisions[7] = "one"
	teamsToDivisions[8] = "one"

	teamsToDivisions[4] = "two"
	teamsToDivisions[5] = "two"
	teamsToDivisions[6] = "two"
	teamsToDivisions[9] = "two"

	divisionOne := []int{1, 2, 3, 7, 8}
	divisionTwo := []int{4, 5, 6, 9}

	divisions := map[string][]int{}
	divisions["one"] = divisionOne
	divisions["two"] = divisionTwo

	inputGames := make([]state.Game, 0)
	inputGames = append(inputGames, state.Game{TeamOne: 1, TeamTwo: 2, GameTime: 1})
	inputGames = append(inputGames, state.Game{TeamOne: 3, TeamTwo: 4, GameTime: 2})
	inputGames = append(inputGames, state.Game{TeamOne: 5, TeamTwo: 6, GameTime: 3})

	inputAsks := map[int][]int{}
	inputAsks[1] = []int{3}
	inputAsks[2] = []int{3}
	inputAsks[3] = []int{2}
	inputAsks[6] = []int{3}

	_ = satisfyAsks(inputGames, inputAsks, teamsToDivisions, divisions)
}

func Test_ValidTimes(t *testing.T) {
	inputTeams := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	daysInSchedule := map[int][]int{
		1: {1, 2, 3},
		2: {4, 5, 6},
		3: {7, 8, 9},
	}
	inputAsks := map[int][]int{}
	inputAsks[1] = []int{3}
	inputAsks[2] = []int{3}
	inputAsks[3] = []int{2}
	inputAsks[6] = []int{3}
	positions := validPositions(inputTeams, daysInSchedule, inputAsks)
	if len(positions) != 9 {
		println("Expected 9 positions, got", len(positions))
	}
}
