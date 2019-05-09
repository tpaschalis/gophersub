package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
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
		{5, time.Duration(time.Second*17 + time.Millisecond*929), time.Duration(time.Second*19 + time.Millisecond*751), `Σας προσφέρω την επιλογή...`},
	}
	_ = parsedSRTFile
	//PaceSubtitleFile(parsedSRTFile, 2.)

	//ParseSRTFile("samples/sample.srt")
	//TimestampToDurationSRT("00:00:17,929 -->")
	//TimestampToDurationSRT("--> 00:00:19,751")

	//a, _ := ParseSRTFile("samples/sample.srt")

	//fmt.Println("parsedSRTFile :\n\n", parsedSRTFile)
	//fmt.Println("funcResult :\n\n", a)
	//_ = a
	//fmt.Println(cmp.Equal(a, parsedSRTFile))

	//b, c := ParseSRTFile("samples/sample_wrong_timestamps.srt")
	//fmt.Println("-----------")
	//fmt.Println(b)
	//fmt.Println("-----------")
	//fmt.Println(c)
	//fmt.Println("-----------")
	//_, _ = b, c
}

func DurationToTimestampSRT(d time.Duration) string {
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
	//fmt.Println("----------------------------------------")
	//fmt.Println(len(splitInput))
	//fmt.Println(splitInput)
	//fmt.Println(splitInput[0], hour)
	//fmt.Println(splitInput[1], minute)
	//fmt.Println(splitInput[2], second)
	//fmt.Println(splitInput[3], millisecond)

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

func TimeshiftSubtitleFile(in SubtitleFile, shift time.Duration) SubtitleFile {
	var res SubtitleFile
	for _, sub := range in {
		sub.Start = sub.Start + shift
		sub.End = sub.End + shift
		res = append(res, sub)
	}
	return res
}

func PaceSubtitleFile(in SubtitleFile, rate float64) (SubtitleFile, error) {
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

type SRTReader struct {
	s *bufio.Scanner
	b bytes.Buffer
	p int
}

func ParseSRTFile(filename string) (SubtitleFile, []error) {
	var res SubtitleFile
	var errCollection []error

	content, err := ioutil.ReadFile(filename)
	file, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return res, []error{errors.New("Something went wrong while trying to parse the provided file!")}
	}
	srt := string(content)
	_ = srt

	reader := new(SRTReader)
	reader.p, reader.s = 0, bufio.NewScanner(file)
	reader.s.Split(SRTScanner)

	for reader.s.Scan() {
		var current Subtitle
		//cur := strings.Split(reader.s.Text(), "\n")
		cur := strings.FieldsFunc(reader.s.Text(), EOLSplit)

		idx, err1 := strconv.Atoi(cur[0])
		if err1 != nil {
			fmt.Println("Error in block :", cur, "\nCould not parse an index correctly")
			errCollection = append(errCollection, err1)
		} else {
			current.Index = idx
		}

		r1, _ := regexp.Compile(`\d+[,.:]\d+[,.:]\d+[,.:]\d+ -[ -]>`)
		r2, _ := regexp.Compile(`-[ -]> \d+[,.:]\d+[,.:]\d+[,.:]\d+`)

		start, err2 := TimestampToDurationSRT(r1.FindString(cur[1]))
		if err2 != nil {
			fmt.Println("Error in block :", cur, "\nCould not parse Start Timestamp correctly")
			errCollection = append(errCollection, err2)
		} else {
			current.Start = start
		}

		end, err3 := TimestampToDurationSRT(r2.FindString(cur[1]))
		if err3 != nil {
			fmt.Println("Error in block :", cur, "\nCould not parse End Timestamp correctly")
			errCollection = append(errCollection, err3)
		} else {
			current.End = end
		}

		current.Content = strings.Join(cur[2:len(cur)], "\n")
		res = append(res, current)
	}

	fmt.Println("Completed parsing the following SRT file :", filename)
	fmt.Println("Parsed a total of", len(res), "subtitle lines")
	if len(errCollection) != 0 {
		fmt.Println("Encountered a total of", len(errCollection), "issues, as presented below :")
		fmt.Println(errCollection)
	}

	return res, errCollection
}

func SRTScanner(data []byte, atEOF bool) (adv int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if i < len(data)-1 && string(data[i:i+2]) == "\n\n" {
			// Do not return the empty lines at the end of the detected block
			return i + 2, data[:i+2], nil
			//return i + 2, data[:i], nil
		}
		if i < len(data)-3 && string(data[i:i+4]) == "\r\n\r\n" {
			// Do not return the empty lines at the end of the detected block
			return i + 4, data[:i+4], nil
			//return i + 4, data[:i], nil
		}
	}
	if atEOF && len(data) != 0 {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func TimestampSplitSRT(r rune) bool {
	return r == ',' || r == '.' || r == ':'
}

func EOLSplit(r rune) bool {
	return r == '\n' || r == '\r'
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
