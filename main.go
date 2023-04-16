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

// StartGin starts gin web server with setting router.
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
		prices, err := moat.GetPricesForSymbolOnTradingDay(day, symbol)

		res := response{
			Prices:    prices,
			Error:     "",
			TimeTaken: time.Now().Sub(startTime).String(),
		}

		if err != nil {
			res.Error = err.Error()
			context.JSON(http.StatusOK, res)
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
