package timeline

import (
	"context"
	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
	"moat/state"
	"os"
	"time"
)

type timeline struct {
	Context context.Context
	Client  *polygon.Client
}

func CreateTimelineClient() *timeline {
	polygonApiKey := os.Getenv("POLYGON_API_KEY")
	return &timeline{
		Context: context.Background(),
		Client:  polygon.New(polygonApiKey),
	}
}

func (timeline *timeline) MarketDaysInRange(startDay *state.MktDate, endDay *state.MktDate) map[string]bool {
	if startDay.IsOnOrAfter(endDay) {
		panic("Start day must be before end day")
	}
	days := map[string]bool{}
	counter := 0
	loc, _ := time.LoadLocation("GMT")
	start := time.Date(startDay.Year, time.Month(startDay.Month), startDay.Day, 0, 0, 0, 0, loc)
	for {
		currentDate := start.AddDate(0, 0, counter)
		plumDate := state.NewMktDate(int(currentDate.Month()), currentDate.Day(), currentDate.Year())
		if plumDate.Equals(endDay) {
			break
		}
		if timeline.isMarketOpen(plumDate) {
			days[plumDate.ToKey()] = true
		}
		counter++
	}
	return days
}

func (timeline *timeline) Yesterday(mktDate *state.MktDate) *state.MktDate {
	counter := -3
	loc, _ := time.LoadLocation("GMT")
	start := time.Date(mktDate.Year, time.Month(mktDate.Month), mktDate.Day, 0, 0, 0, 0, loc)
	lastOpenDay := state.NewMktDate(1, 1, 1000)
	for {
		currentDate := start.AddDate(0, 0, counter)
		plumDate := state.NewMktDate(int(currentDate.Month()), currentDate.Day(), currentDate.Year())
		if plumDate.Equals(mktDate) {
			break
		}
		if timeline.isMarketOpen(plumDate) {
			lastOpenDay = plumDate
		}
		counter++
	}
	return lastOpenDay
}

func (timeline *timeline) YesterdaysInRange(startDay *state.MktDate, endDay *state.MktDate) map[string]string {
	if startDay.IsOnOrAfter(endDay) {
		panic("Start day must be before end day")
	}
	days := map[string]string{}
	counter := -10
	loc, _ := time.LoadLocation("GMT")
	start := time.Date(startDay.Year, time.Month(startDay.Month), startDay.Day, 0, 0, 0, 0, loc)
	lastOpenDay := state.NewMktDate(1, 1, 1000)
	for {
		currentDate := start.AddDate(0, 0, counter)
		plumDate := state.NewMktDate(int(currentDate.Month()), currentDate.Day(), currentDate.Year())
		if plumDate.Equals(endDay) {
			break
		}
		if timeline.isMarketOpen(plumDate) {
			days[plumDate.ToKey()] = lastOpenDay.ToKey()
			lastOpenDay = plumDate
		}
		counter++
	}
	return days
}

func (timeline *timeline) isMarketOpen(day *state.MktDate) bool {
	loc, _ := time.LoadLocation("GMT")
	_, err := timeline.Client.GetDailyOpenCloseAgg(context.Background(), models.GetDailyOpenCloseAggParams{
		Ticker: "SPY",
		Date:   models.Date(time.Date(day.Year, time.Month(day.Month), day.Day, 0, 0, 0, 0, loc)),
	}.WithAdjusted(true))
	if err != nil {
		//println(day.String() + " is not a market day")
		return false
	}
	return true
}
