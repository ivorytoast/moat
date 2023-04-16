package main

import (
	"context"
	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
	"moat/state"
	"os"
	"sort"
	"strconv"
	"time"
)

type Moat struct {
	Context context.Context
	Client  *polygon.Client
}

func CreateMoatClient() *Moat {
	polygonApiKey := os.Getenv("POLYGON_API_KEY")
	return &Moat{
		Context: context.Background(),
		Client:  polygon.New(polygonApiKey),
	}
}

func (moat Moat) GetPricesForSymbolOnTradingDay(day state.Day, symbol string) ([]state.PriceInfo, error) {
	polygonResponse, err := moat.Client.GetAggs(moat.Context, models.GetAggsParams{
		Ticker:     symbol,
		Multiplier: 1,
		Timespan:   models.Minute,
		From:       models.Millis(time.Date(day.Year, time.Month(day.Month), day.Day, 13, 30, 0, 0, time.UTC)),
		To:         models.Millis(time.Date(day.Year, time.Month(day.Month), day.Day, 20, 0, 0, 0, time.UTC)),
	}.WithOrder(models.Asc).WithAdjusted(true))
	if err != nil {
		return nil, err
	}
	output := getPricesForSymbolOnTradingDayImpl(polygonResponse)
	return output, nil
}

func getPricesForSymbolOnTradingDayImpl(polygonResponse *models.GetAggsResponse) []state.PriceInfo {
	prices := map[string]float64{}
	for _, agg := range polygonResponse.Results {
		timestamp, err := agg.Timestamp.MarshalJSON()
		if err != nil {
			println(err.Error())
			return nil
		}
		prices[string(timestamp)] = agg.Close
	}

	priceObjects := make([]state.PriceInfo, 0)
	for timestamp, price := range prices {
		_, readableTimestamp := convertToEst(timestamp)
		priceObjects = append(priceObjects, state.PriceInfo{
			Timestamp: timestamp,
			Readable:  readableTimestamp,
			Price:     price,
		})
	}

	sort.Slice(priceObjects, func(i, j int) bool {
		one := priceObjects[i].Timestamp
		two := priceObjects[j].Timestamp
		return one < two
	})

	return priceObjects
}

func (moat Moat) GetTimestampsForTradingDay(day state.Day) ([]state.TimestampInfo, error) {
	polygonResponse, err := moat.Client.GetAggs(moat.Context, models.GetAggsParams{
		Ticker:     "SPY",
		Multiplier: 1,
		Timespan:   models.Minute,
		From:       models.Millis(time.Date(day.Year, time.Month(day.Month), day.Day, 13, 30, 0, 0, time.UTC)),
		To:         models.Millis(time.Date(day.Year, time.Month(day.Month), day.Day, 20, 0, 0, 0, time.UTC)),
	}.WithOrder(models.Asc).WithAdjusted(true))
	if err != nil {
		return nil, err
	}
	output := getTimestampsForTradingDayImpl(polygonResponse)
	return output, nil
}

func getTimestampsForTradingDayImpl(polygonResponse *models.GetAggsResponse) []state.TimestampInfo {
	timestamps := map[string]string{}
	for _, agg := range polygonResponse.Results {
		timestamp, err := agg.Timestamp.MarshalJSON()
		if err != nil {
			println(err.Error())
			return nil
		}
		_, readableTimestamp := convertToEst(string(timestamp))
		timestamps[string(timestamp)] = readableTimestamp
	}

	timestampObjects := make([]state.TimestampInfo, 0)
	for key, value := range timestamps {
		timestampObjects = append(timestampObjects, state.TimestampInfo{
			Timestamp: key,
			Readable:  value,
		})
	}

	sort.Slice(timestampObjects, func(i, j int) bool {
		one := timestampObjects[i].Timestamp
		two := timestampObjects[j].Timestamp
		return one < two
	})

	return timestampObjects
}

func convertToEst(timestamp string) (time.Time, string) {
	unix, _ := strconv.ParseInt(timestamp, 10, 64)
	t := time.UnixMilli(unix)

	loc, _ := time.LoadLocation("America/New_York")
	newYorkTime := t.In(loc)

	newYorkTimeHours := (newYorkTime.Hour()+11)%12 + 1

	weekday := newYorkTime.Weekday().String()
	month := newYorkTime.Month().String()
	day := strconv.Itoa(newYorkTime.Day())
	year := strconv.Itoa(newYorkTime.Year())
	hour := strconv.Itoa(newYorkTimeHours)
	amOrPM := "AM"
	if newYorkTime.Hour() >= 12 {
		amOrPM = "PM"
	}
	minute := strconv.Itoa(newYorkTime.Minute())
	if len(minute) < 2 {
		minute = "0" + minute
	}
	second := strconv.Itoa(newYorkTime.Second())
	if len(second) < 2 {
		second = "0" + second
	}
	asString := weekday + " " + month + " " + day + " " + year + " " + hour + ":" + minute + ":" + second + " " + amOrPM
	return newYorkTime, asString
}
