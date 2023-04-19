package state

import (
	"fmt"
	"strconv"
)

type ScheduleState struct {
	Teams            []int            `json:"teams"`
	DaysInSchedule   map[int][]int    `json:"days_in_schedule"`
	TeamsToDivisions map[int]string   `json:"teams_to_divisions"`
	Divisions        map[string][]int `json:"divisions"`
	Games            []Game           `json:"games"`
}

type Game struct {
	GameTime int `json:"game_time"`
	TeamOne  int `json:"team_one"`
	TeamTwo  int `json:"team_two"`
}

type TimestampInfo struct {
	Timestamp string `json:"timestamp"`
	Readable  string `json:"readable"`
}

type PriceInfo struct {
	Timestamp string  `json:"timestamp"`
	Readable  string  `json:"readable"`
	Price     float64 `json:"price"`
}

type CorrelationInfo struct {
	Timestamp              string  `json:"timestamp"`
	Readable               string  `json:"readable"`
	SymbolPrice            float64 `json:"symbol_price"`
	HedgePrice             float64 `json:"hedge_price"`
	SymbolPercentageChange float64 `json:"symbol_percentage_change"`
	HedgePercentageChange  float64 `json:"hedge_percentage_change"`
	SymbolPreviousClose    float64 `json:"symbol_previous_close"`
	HedgePreviousClose     float64 `json:"hedge_previous_close"`
	PreviousDayTimestamp   string  `json:"previous_day_timestamp"`
	PreviousDayReadable    string  `json:"previous_day_readable"`
	Correlation            float64 `json:"correlation"`
}

type MktDate struct {
	Month int
	Day   int
	Year  int
}

func NewMktDate(month int, day int, year int) *MktDate {
	if month < 1 || month > 12 {
		panic("Invalid month")
	}
	if day < 1 || day > 31 {
		panic("Invalid day")
	}
	if year < 0 {
		panic("Invalid year")
	}
	return &MktDate{Month: month, Day: day, Year: year}
}

func (day *MktDate) String() string {
	return fmt.Sprintf("%d/%d/%d", day.Month, day.Day, day.Year)
}

func (day *MktDate) IsBefore(other *MktDate) bool {
	if day.Year < other.Year {
		return true
	}
	if day.Year > other.Year {
		return false
	}
	if day.Month < other.Month {
		return true
	}
	if day.Month > other.Month {
		return false
	}
	return day.Day < other.Day
}

func (day *MktDate) IsAfter(other *MktDate) bool {
	return other.IsBefore(day)
}

func (day *MktDate) Equals(other *MktDate) bool {
	return day.Year == other.Year && day.Month == other.Month && day.Day == other.Day
}

func (day *MktDate) IsBetween(start *MktDate, end *MktDate) bool {
	return day.IsAfter(start) && day.IsBefore(end)
}

func (day *MktDate) IsOnOrBefore(other *MktDate) bool {
	return day.IsBefore(other) || day.Equals(other)
}

func (day *MktDate) IsOnOrAfter(other *MktDate) bool {
	return day.IsAfter(other) || day.Equals(other)
}

func (day *MktDate) ToKey() string {
	key := ""
	key = key + strconv.Itoa(day.Month) + "-"
	key = key + strconv.Itoa(day.Day) + "-"
	key = key + strconv.Itoa(day.Year)
	return key
}
