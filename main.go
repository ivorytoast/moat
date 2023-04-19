package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"log"
	"math"
	"moat/scheduler"
	"moat/state"
	"moat/timeline"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	ConfigRuntime()
	StartGin()
}

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

func StartGin() {
	moatClient := CreateMoatClient()
	timelineClient := timeline.CreateTimelineClient()
	schedulerClient := scheduler.CreateSchedulerClient()

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "Hello World From Golang!",
		})
	})

	router.GET("/api/v0/timestamps", func(context *gin.Context) {
		println("GET /api/v0/timestamps called")
		startTime := time.Now()

		type response struct {
			Timestamps []state.TimestampInfo `json:"timestamps"`
			Error      string                `json:"error"`
			TimeTaken  string                `json:"time_taken"`
		}

		mktDate := state.NewMktDate(4, 14, 2023)
		timestamps, err := moatClient.GetTimestampsForTradingDay(mktDate)

		res := response{
			Timestamps: timestamps,
			Error:      "",
			TimeTaken:  time.Now().Sub(startTime).String(),
		}

		if err != nil {
			res.Error = err.Error()
			context.JSON(http.StatusOK, res)
		}
		fmt.Print(res)
		context.JSON(http.StatusOK, res)
	})

	router.GET("/api/v0/prices", func(context *gin.Context) {
		symbol := context.Query("symbol")
		println("GET /api/v0/timestamps called with query param 'symbol' as: " + symbol)

		startTime := time.Now()

		type response struct {
			Prices    []state.PriceInfo `json:"prices"`
			Error     string            `json:"error"`
			TimeTaken string            `json:"time_taken"`
		}

		mktDate := state.NewMktDate(4, 14, 2023)
		prices, pricesErr := moatClient.GetPricesForSymbolOnTradingDay(mktDate, symbol)

		res := response{
			Prices:    nil,
			Error:     "",
			TimeTaken: time.Now().Sub(startTime).String(),
		}

		if pricesErr != nil {
			res.Error = pricesErr.Error()
			context.JSON(http.StatusOK, res)
		}

		res.Prices = prices
		context.JSON(http.StatusOK, res)
	})

	router.GET("/api/v0/correlations", func(context *gin.Context) {
		symbol := context.Query("symbol")
		hedge := context.Query("hedge")
		println("GET /api/v0/correlations called with query param 'symbol' as: " + symbol + " and 'hedge' as: " + hedge)

		startTime := time.Now()

		type response struct {
			Correlations     []state.CorrelationInfo `json:"correlations"`
			TotalCorrelation float64                 `json:"total_correlation"`
			Error            string                  `json:"error"`
			TimeTaken        string                  `json:"time_taken"`
		}

		mktDate := state.NewMktDate(4, 14, 2023)

		timestamps, timestampsErr := moatClient.GetTimestampsForTradingDay(mktDate)
		if timestampsErr != nil {
			context.JSON(http.StatusOK, response{
				Correlations: nil,
				Error:        timestampsErr.Error(),
				TimeTaken:    time.Now().Sub(startTime).String(),
			})
			return
		}

		symbolPrices, symbolErr := moatClient.GetPricesForSymbolOnTradingDay(mktDate, symbol)
		if symbolErr != nil {
			context.JSON(http.StatusOK, response{
				Correlations: nil,
				Error:        symbolErr.Error(),
				TimeTaken:    time.Now().Sub(startTime).String(),
			})
			return
		}

		hedgePrices, hedgeErr := moatClient.GetPricesForSymbolOnTradingDay(mktDate, hedge)
		if hedgeErr != nil {
			context.JSON(http.StatusOK, response{
				Correlations: nil,
				Error:        hedgeErr.Error(),
				TimeTaken:    time.Now().Sub(startTime).String(),
			})
			return
		}

		if len(symbolPrices) != len(hedgePrices) || len(timestamps) != len(symbolPrices) {
			context.JSON(http.StatusOK, response{
				Correlations: nil,
				Error:        "the timestamps, symbol, and hedge do not have the same amount of prices",
				TimeTaken:    time.Now().Sub(startTime).String(),
			})
			return
		}

		res := response{
			Correlations:     nil,
			Error:            "",
			TotalCorrelation: 0.0,
			TimeTaken:        time.Now().Sub(startTime).String(),
		}

		previousMktDate, yesterdayErr := timelineClient.Yesterday(mktDate)
		if yesterdayErr != nil {
			res.Error = yesterdayErr.Error()
			context.JSON(http.StatusOK, res)
		}

		previousSymbolPriceInfo, previousSymErr := moatClient.GetClosePriceForSymbolOnTradingDay(previousMktDate, symbol)
		if previousSymErr != nil {
			res.Error = previousSymErr.Error()
			context.JSON(http.StatusOK, res)
		}

		previousHedgePriceInfo, previousHedgeErr := moatClient.GetClosePriceForSymbolOnTradingDay(previousMktDate, hedge)
		if previousHedgeErr != nil {
			res.Error = previousHedgeErr.Error()
			context.JSON(http.StatusOK, res)
		}

		if previousSymbolPriceInfo.Timestamp != previousHedgePriceInfo.Timestamp {
			res.Error = "the previous day for the symbol and hedge do not have the same timestamps. The timestamps are: (symbol) " + previousSymbolPriceInfo.Timestamp + ", (hedge) " + previousHedgePriceInfo.Timestamp
			context.JSON(http.StatusOK, res)
		}

		correlationSum := 0.0
		correlationCounter := 0.0
		latestCorrelation := 0.0
		correlations := make([]state.CorrelationInfo, 0)
		for i := 0; i < len(timestamps); i++ {

			symbolPercentageChange := ((symbolPrices[i].Price - previousSymbolPriceInfo.Price) / previousSymbolPriceInfo.Price) * 100
			hedgePercentageChange := ((hedgePrices[i].Price - previousHedgePriceInfo.Price) / previousHedgePriceInfo.Price) * 100

			if symbolPrices[i].Timestamp != timestamps[i].Timestamp || hedgePrices[i].Timestamp != timestamps[i].Timestamp {
				res.Error = "the symbol, hedge, and timestamps do not have the same timestamps. The timestamps are: (timestamp) " + timestamps[i].Timestamp + ", (symbol) " + symbolPrices[i].Timestamp + ", (hedge) " + hedgePrices[i].Timestamp
				context.JSON(http.StatusOK, res)
			}

			corr := 0.0
			if hedgePercentageChange == 0 {
				corr = latestCorrelation
			} else {
				corr = symbolPercentageChange / hedgePercentageChange
				latestCorrelation = corr
			}

			if math.Abs(hedgePercentageChange) <= 0.50 || math.Abs(symbolPercentageChange) <= 0.50 {
				corr = 0.0
			}

			if corr != 0.0 {
				correlationSum += corr
				correlationCounter += 1.0
			}

			correlation := state.CorrelationInfo{
				Timestamp:              timestamps[i].Timestamp,
				Readable:               timestamps[i].Readable,
				SymbolPrice:            symbolPrices[i].Price,
				HedgePrice:             hedgePrices[i].Price,
				SymbolPercentageChange: symbolPercentageChange,
				HedgePercentageChange:  hedgePercentageChange,
				SymbolPreviousClose:    previousSymbolPriceInfo.Price,
				HedgePreviousClose:     previousHedgePriceInfo.Price,
				PreviousDayTimestamp:   previousSymbolPriceInfo.Timestamp,
				PreviousDayReadable:    previousSymbolPriceInfo.Readable,
				Correlation:            corr,
			}

			correlations = append(correlations, correlation)
		}

		if correlationSum == 0 {
			res.Error = "no correlations found"
			context.JSON(http.StatusOK, res)
			return
		}

		if correlationCounter == 0 {
			res.Error = "not enough data to calculate correlation"
			context.JSON(http.StatusOK, res)
			return
		}

		res.Correlations = correlations
		res.TotalCorrelation = correlationSum / correlationCounter
		context.JSON(http.StatusOK, res)
	})

	router.GET("/api/v0/scheduler", func(context *gin.Context) {
		startTime := time.Now()

		schedule := schedulerClient.Schedule()

		type response struct {
			Schedule  *state.ScheduleState `json:"schedule"`
			Error     string               `json:"error"`
			TimeTaken string               `json:"time_taken"`
		}

		res := response{
			Schedule:  schedule,
			Error:     "",
			TimeTaken: time.Now().Sub(startTime).String(),
		}

		context.JSON(http.StatusOK, res)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
