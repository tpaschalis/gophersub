package main

import (
	"fmt"
	"math"
	"time"
)

const (
	day  = time.Minute * 60 * 24
	year = 365 * day
)

type Subtitle struct {
	Index    int
	Start    time.Duration
	End      time.Duration
	Metadata string
}

func main() {

}

func DurationToTimestamp(d time.Duration) string {
	var hour, minute int
	var second float64
	stringDuration, err := time.ParseDuration(d.String())
	if err != nil {
		fmt.Println("Could not parse provided time.Duration")
		panic(err)
	}
	hour = int(stringDuration.Hours())
	minute = int(math.Mod(stringDuration.Minutes(), 60))
	second = math.Mod(stringDuration.Seconds(), 60)

	res := fmt.Sprintf("%02d:%02d:%02.3f", hour, minute, second)
	return res
}

func StrToDuration(in string) (time.Duration, error) {
	var res time.Duration

	res, err := time.ParseDuration(in)
	if err != nil {
		//fmt.Println("Could not parse provided time.Duration")
		return res, err
	}
	return res, nil
}
