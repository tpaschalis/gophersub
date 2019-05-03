package main

import (
	"fmt"
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
	s := time.Duration(time.Hour*2 + time.Minute*10 + time.Second*20 + time.Millisecond*183)
	e := time.Duration(time.Hour*2 + time.Minute*80 + time.Second*90 + time.Millisecond*1800)
	foo := time.Duration(time.Hour*1 + time.Minute*0 + time.Second*10 + time.Millisecond*0)
	bar := time.Duration(time.Hour*20 + time.Minute*5 + time.Second*0 + time.Millisecond*70)
	//fmt.Println(s.String())
	//fmt.Println(e.String())
	//2h10m20.183s
	//00:51:32,817
	DurationToTimestamp(s)
	DurationToTimestamp(e)
	DurationToTimestamp(foo)
	DurationToTimestamp(bar)
}

func DurationToTimestamp(d time.Duration) string {
	var hour, minute int
	var second float64
	fmt.Sscanf(d.String(), "%dh%dm%fs", &hour, &minute, &second)
	res := fmt.Sprintf("%02d:%02d:%02.3f", hour, minute, second)
	//fmt.Println(d)
	//fmt.Println(res)
	//fmt.Println("")
	return res
}
