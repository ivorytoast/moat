package timeline

import (
	"context"
	"errors"
	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
	"moat/state"
	"os"
	"time"
)

type Timeline struct {
	Context context.Context
	Client  *polygon.Client
}

func CreateTimelineClient() *Timeline {
	polygonApiKey := os.Getenv("POLYGON_API_KEY")
	return &Timeline{
		Context: context.Background(),
		Client:  polygon.New(polygonApiKey),
	}
}

func (timeline *Timeline) isMarketOpen(day *state.MktDate) bool {
	loc, _ := time.LoadLocation("America/New_York")
	_, err := timeline.Client.GetDailyOpenCloseAgg(context.Background(), models.GetDailyOpenCloseAggParams{
		Ticker: "SPY",
		Date:   models.Date(time.Date(day.Year, time.Month(day.Month), day.Day, 0, 0, 0, 0, loc)),
	}.WithAdjusted(true))
	if err != nil {
		return false
	}
	return true
}

func (timeline *Timeline) Yesterday(mktDate *state.MktDate) (*state.MktDate, error) {
	loc, _ := time.LoadLocation("America/New_York")
	start := time.Date(mktDate.Year, time.Month(mktDate.Month), mktDate.Day, 0, 0, 0, 0, loc)
	counter := -1
	for {
		currentDate := start.AddDate(0, 0, counter)
		previousMktDate := state.NewMktDate(int(currentDate.Month()), currentDate.Day(), currentDate.Year())
		if timeline.isMarketOpen(previousMktDate) {
			return previousMktDate, nil
		}
		counter--
		if counter < -15 {
			break
		}
	}
	return nil, errors.New("there is a serious issue here...for 15 days (2+ weeks) polygon.io said the market was not open")
}
