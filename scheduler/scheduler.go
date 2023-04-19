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
