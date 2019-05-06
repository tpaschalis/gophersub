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

	//SecArgToDuration("1.4s")
	//SecArgToDuration("5s")
	//	SecArgToDuration("-1.214211232321")
	//	SecArgToDuration("1.1s")
	//	SecArgToDuration("s")
	//	SecArgToDuration("0s")
	//	SecArgToDuration("0.181s")
	//	SecArgToDuration("0.00001s")
	//	SecArgToDuration("0.0000s")

	//e := time.Duration(time.Hour*0 + time.Minute*0 + time.Second*0 + time.Millisecond*181)
	e2 := time.Duration(time.Millisecond * 181)

	s := time.Duration(time.Hour*2 + time.Minute*10 + time.Second*20 + time.Millisecond*183)
	e := time.Duration(time.Hour*2 + time.Minute*80 + time.Second*90 + time.Millisecond*1800)
	foo := time.Duration(time.Hour*1 + time.Minute*0 + time.Second*10 + time.Millisecond*0)
	bar := time.Duration(time.Hour*20 + time.Minute*5 + time.Second*0 + time.Millisecond*70)
	fmt.Println(DurationToTimestamp(e))
	DurationToTimestamp(s)
	DurationToTimestamp(e)
	DurationToTimestamp(foo)
	DurationToTimestamp(bar)
	DurationToTimestamp(e2)
}

func DurationToTimestamp(d time.Duration) string {
	var hour, minute int
	var second float64
	stringDuration, _ := time.ParseDuration(d.String())
	hour = int(stringDuration.Hours())
	minute = int(math.Mod(stringDuration.Minutes(), 60))
	second = math.Mod(stringDuration.Seconds(), 60)

	res := fmt.Sprintf("%02d:%02d:%02.3f", hour, minute, second)
	return res
}

func SecArgToDuration(arg string) {
	var in float64
	fmt.Sscanf(arg, "%fs", &in)
	second := int(in)
	milli := math.Round((in - float64(second)) * 1000)
	fmt.Println(arg)
	fmt.Println(in, second, milli)
	res := time.Duration(time.Second*time.Duration(second) + time.Millisecond*time.Duration(milli))
	fmt.Println(res)
	fmt.Println(DurationToTimestamp(res))
	fmt.Println("")
}
