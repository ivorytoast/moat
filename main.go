package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"log"
	"moat/state"
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
	moat := CreateMoatClient()

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

		day := state.Day{
			Month: 4,
			Day:   14,
			Year:  2023,
		}
		timestamps, err := moat.GetTimestampsForTradingDay(day)

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

		day := state.Day{
			Month: 4,
			Day:   14,
			Year:  2023,
		}
		prices, pricesErr := moat.GetPricesForSymbolOnTradingDay(day, symbol)

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
			Correlations []state.CorrelationInfo `json:"correlations"`
			Error        string                  `json:"error"`
			TimeTaken    string                  `json:"time_taken"`
		}

		day := state.Day{
			Month: 4,
			Day:   14,
			Year:  2023,
		}

		timestamps, timestampsErr := moat.GetTimestampsForTradingDay(day)
		if timestampsErr != nil {
			context.JSON(http.StatusOK, response{
				Correlations: nil,
				Error:        timestampsErr.Error(),
				TimeTaken:    time.Now().Sub(startTime).String(),
			})
			return
		}

		symbolPrices, symbolErr := moat.GetPricesForSymbolOnTradingDay(day, symbol)
		if symbolErr != nil {
			context.JSON(http.StatusOK, response{
				Correlations: nil,
				Error:        symbolErr.Error(),
				TimeTaken:    time.Now().Sub(startTime).String(),
			})
			return
		}

		hedgePrices, hedgeErr := moat.GetPricesForSymbolOnTradingDay(day, hedge)
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
			Correlations: nil,
			Error:        "",
			TimeTaken:    time.Now().Sub(startTime).String(),
		}

		for i := 0; i < len(timestamps); i++ {

			if symbolPrices[i].Timestamp != timestamps[i].Timestamp || hedgePrices[i].Timestamp != timestamps[i].Timestamp {
				res.Error = "the symbol, hedge, and timestamps do not have the same timestamps. The timestamps are: (timestamp) " + timestamps[i].Timestamp + ", (symbol) " + symbolPrices[i].Timestamp + ", (hedge) " + hedgePrices[i].Timestamp
				context.JSON(http.StatusOK, res)
			}

			correlation := state.CorrelationInfo{
				Timestamp:   timestamps[i].Timestamp,
				Readable:    timestamps[i].Readable,
				SymbolPrice: symbolPrices[i].Price,
				HedgePrice:  hedgePrices[i].Price,
				Correlation: symbolPrices[i].Price / hedgePrices[i].Price,
			}

			res.Correlations = append(res.Correlations, correlation)
		}

		fmt.Print(res)
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
