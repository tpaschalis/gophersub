package main

import (
	"errors"
	"fmt"
	"math"
	"time"
)

const (
	day  = time.Minute * 60 * 24
	year = 365 * day
)

type Subtitle struct {
	Index   int
	Start   time.Duration
	End     time.Duration
	Content string
}

type SubtitleFile []Subtitle

func main() {
	parsedSRTFile := SubtitleFile{
		{1, time.Duration(time.Second*1 + time.Millisecond*602), time.Duration(time.Second*3 + time.Millisecond*314), `Έχουμε όλοι υποφέρει.`},
		{2, time.Duration(time.Second*4 + time.Millisecond*536), time.Duration(time.Second*7 + time.Millisecond*379), `Έχουμε χάσει αγαπημένους μας.`},
		{3, time.Duration(time.Second*10 + time.Millisecond*88), time.Duration(time.Second*14 + time.Millisecond*500), `Αυτό δεν αφορά τους Οίκους των ευγενών,
αλλά τους ζωντανούς και τους νεκρούς.`},
		{4, time.Duration(time.Second*14 + time.Millisecond*611), time.Duration(time.Second*16 + time.Millisecond*568), `Κι εγώ σκοπεύω να ζήσω.`},
		{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή..`},
	}

	PaceSRTFile(parsedSRTFile, 2.)

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

func TimeshiftSRTFile(in SubtitleFile, shift time.Duration) SubtitleFile {
	var res SubtitleFile
	for _, sub := range in {
		sub.Start = sub.Start + shift
		sub.End = sub.End + shift
		res = append(res, sub)
	}
	return res
}

func PaceSRTFile(in SubtitleFile, rate float64) (SubtitleFile, error) {
	var res SubtitleFile
	if rate <= 0 {
		return res, errors.New("Input rate should be a positive, floating-point number")
	}

	whole, frac := math.Modf(1. / rate)
	for _, sub := range in {
		sub.Start = sub.Start*time.Duration(whole) + sub.Start*time.Duration(int(frac*1000))/1000
		sub.End = sub.End*time.Duration(whole) + sub.End*time.Duration(int(frac*1000))/1000
		res = append(res, sub)
	}
	//fmt.Println(in)
	//fmt.Println("--------------------------------------=============-----------------------------------------------")
	//fmt.Println(res)
	return res, nil
}
