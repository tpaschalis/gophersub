package main

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func DurationToTimestampSRT(d time.Duration) string {
	var hour, minute, second, millisec float64
	stringDuration, err := time.ParseDuration(d.String())
	if err != nil {
		//fmt.Println("Could not parse provided time.Duration")
		panic(err)
	}
	hour = stringDuration.Hours()
	minute = math.Mod(stringDuration.Minutes(), 60)
	second = math.Mod(stringDuration.Seconds(), 60)
	_, millisec = math.Modf(second)
	millisec = math.Round(millisec * 1000)

	res := fmt.Sprintf("%02d:%02d:%02d,%03d", int(hour), int(minute), int(second), int(millisec))
	return res
}

func TimestampToDurationSRT(in string) (time.Duration, error) {
	var res time.Duration

	r, _ := regexp.Compile("[^0-9,.:]")
	in = r.ReplaceAllString(in, "")
	splitInput := strings.FieldsFunc(in, TimestampSplitSRT)
	if len(splitInput) != 4 {
		return res, errors.New("Wrong Number of fields resulting from input timestamp")
	}

	// TODO : Implement concise error handling
	hour, _ := strconv.Atoi(splitInput[0])
	minute, _ := strconv.Atoi(splitInput[1])
	second, _ := strconv.Atoi(splitInput[2])
	millisecond, _ := strconv.Atoi(splitInput[3])
	if minute > 59 {
		return res, errors.New("Unexpected parsed minute value, should be between 0 and 60")
	}
	if second > 59 {
		return res, errors.New("Unexpected parsed seconds value, should be between 0 and 60")
	}
	if millisecond > 999 {
		return res, errors.New("Unexpected parsed millisecond value, should be between 0 and 999")
	}

	res = time.Duration(time.Hour*time.Duration(hour) + time.Minute*time.Duration(minute) + time.Second*time.Duration(second) + time.Millisecond*time.Duration(millisecond))

	return res, nil
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

func ErrorSlicesEqual(a, b []error) bool {
	if len(a) != len(b) {
		return false
	}
	for i, _ := range a {
		if a[i].Error() != b[i].Error() {
			return false
		}
	}
	return true
}

func ConvertToUTF8(filename string) {
}
