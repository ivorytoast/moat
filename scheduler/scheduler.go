package scheduler

import (
	prioq "github.com/jupp0r/go-priority-queue"
	"math/rand"
	"moat/state"
	"strconv"
	"time"
)

type Scheduler struct{}

func CreateSchedulerClient() *Scheduler {
	return &Scheduler{}
}

// TODO: All of these extra functions can be done by inputting the whole schedule struct
// Therefore, you can simply store the whole schedule struct and pass that around

func (scheduler *Scheduler) Schedule() *state.ScheduleState {
	inputTeams := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	daysInSchedule := map[int][]int{
		1: {1, 2, 3},
		2: {4, 5, 6},
		3: {7, 8, 9},
	}

	teamsToDivisions := map[int]string{}
	divisionOne := []int{1, 2, 7, 8, 9}
	divisionTwo := []int{3, 4, 5, 6}

	teamsToDivisions[1] = "one"
	teamsToDivisions[2] = "one"
	teamsToDivisions[7] = "one"
	teamsToDivisions[8] = "one"
	teamsToDivisions[9] = "one"

	teamsToDivisions[3] = "two"
	teamsToDivisions[4] = "two"
	teamsToDivisions[5] = "two"
	teamsToDivisions[6] = "two"

	divisions := map[string][]int{}
	divisions["one"] = divisionOne
	divisions["two"] = divisionTwo

	teamsToSpots := map[int][]int{}
	for _, team := range inputTeams {
		teamsToSpots[team] = []int{}
	}
	teamsToSpots[1] = []int{1, 2, 4, 5, 7, 8}

	games := make([]state.Game, 0)
	teamUsed := map[int]bool{}
	for dayKey, dayGameTimes := range daysInSchedule {
		shuffledTeams := shuffleTeams(inputTeams)
		for _, team := range shuffledTeams {
			teamUsed[team] = false
		}
		for _, gameTime := range dayGameTimes {
			gameResolved := false
			for _, teamOne := range shuffledTeams {
				if gameResolved {
					break
				}
				if teamUsed[teamOne] {
					continue
				}
				teamUsed[teamOne] = true
				for _, teamTwo := range shuffledTeams {
					if teamUsed[teamTwo] {
						continue
					}
					if teamsToDivisions[teamOne] != teamsToDivisions[teamTwo] {
						continue
					}
					games = append(games, state.Game{GameTime: gameTime, TeamOne: teamOne, TeamTwo: teamTwo})
					gameResolved = true
					teamUsed[teamTwo] = true
					break
				}
			}
		}

		for team, wasUsed := range teamUsed {
			if !wasUsed {
				println("The following team has a bye week: " + strconv.Itoa(team) + " on the following day: " + strconv.Itoa(dayKey))
			}
		}
	}

	scheduleState := &state.ScheduleState{
		Teams:            inputTeams,
		DaysInSchedule:   daysInSchedule,
		TeamsToDivisions: teamsToDivisions,
		Divisions:        divisions,
		Games:            games,
	}

	return scheduleState
}

func isTeamInValidSpot(team int, spot int, asks map[int][]int) bool {
	for _, badSpot := range asks[team] {
		if badSpot == spot {
			return false
		}
	}
	return true
}

func validPositions(teams []int, gameTimes map[int][]int, asks map[int][]int) map[int][]int {
	positions := map[int][]int{} // team to valid game times
	for _, team := range teams {
		positions[team] = []int{}
		for _, times := range gameTimes {
			for _, time := range times {
				if isTeamInValidSpot(team, time, asks) {
					positions[team] = append(positions[team], time)
				}
			}
		}
	}
	return positions
}

// If I know which one is valid
// And if I know where each team can play
// I can just go through the invalid teams and switch with a team that is valid at that time slot and division
func satisfyAsks(games []state.Game, asks map[int][]int, teamsToDivisions map[int]string, divisions map[string][]int) []state.Game {
	outputGames := make([]state.Game, 0)

	teamsNeededToSwitch := make([]int, 0)
	for _, game := range games {
		if !isTeamInValidSpot(game.TeamOne, game.GameTime, asks) {
			println("team " + strconv.Itoa(game.TeamOne) + " is not in a valid spot at " + strconv.Itoa(game.GameTime))
			teamsNeededToSwitch = append(teamsNeededToSwitch, game.TeamOne)
			division := teamsToDivisions[game.TeamOne]
			otherTeamsToPlayAgainst := divisions[division]

			for _, otherTeam := range otherTeamsToPlayAgainst {
				if otherTeam != game.TeamOne && isTeamInValidSpot(otherTeam, game.GameTime, asks) {
					println("we can switch " + strconv.Itoa(game.TeamOne) + " with " + strconv.Itoa(otherTeam) + " at " + strconv.Itoa(game.GameTime))
					//break
				}
			}
		}
		if !isTeamInValidSpot(game.TeamTwo, game.GameTime, asks) {
			println("team " + strconv.Itoa(game.TeamTwo) + " is not in a valid spot at " + strconv.Itoa(game.GameTime))
			teamsNeededToSwitch = append(teamsNeededToSwitch, game.TeamTwo)
			division := teamsToDivisions[game.TeamTwo]
			otherTeamsToPlayAgainst := divisions[division]

			for _, otherTeam := range otherTeamsToPlayAgainst {
				if otherTeam != game.TeamTwo && isTeamInValidSpot(otherTeam, game.GameTime, asks) {
					println("we can switch " + strconv.Itoa(game.TeamTwo) + " with " + strconv.Itoa(otherTeam) + " at " + strconv.Itoa(game.GameTime))
					//break
				}
			}
		}
	}

	return outputGames
}

func shuffleTeams(inputTeams []int) []int {
	if inputTeams == nil || len(inputTeams) == 0 {
		println("input teams is nil or empty")
		return nil
	}

	pq := prioq.New()

	for _, team := range inputTeams {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		weight := r1.Float64()
		pq.Insert(team, 1000.0*weight)
	}

	teamsInWeightedOrder := make([]int, 0)

	if pq.Len() != len(inputTeams) {
		println("the length of the priority queue does not equal the input teams")
		return nil
	}

	for pq.Len() > 0 {
		team, err := pq.Pop()
		if err != nil {
			println("error occurred while popping from priority queue" + err.Error())
			return nil
		}
		teamsInWeightedOrder = append(teamsInWeightedOrder, team.(int))
	}

	return teamsInWeightedOrder
}
