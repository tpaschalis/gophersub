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
	Index    int
	Start    time.Duration
	End      time.Duration
	Content  string
	Metadata string
	Header   string
}

type SubtitleFile struct {
	Subtitles []Subtitle
	Headers   string
}

func main() {
}

func DurationToTimestampSRT(d time.Duration) string {
	var hour, minute int
	var second float64
	stringDuration, err := time.ParseDuration(d.String())
	if err != nil {
		//fmt.Println("Could not parse provided time.Duration")
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
	for _, sub := range in.Subtitles {
		sub.Start = sub.Start + shift
		sub.End = sub.End + shift
		res.Subtitles = append(res.Subtitles, sub)
	}
	return res
}

func PaceSubtitleFile(in SubtitleFile, rate float64) (SubtitleFile, error) {
	var res SubtitleFile
	if rate <= 0 {
		return res, errors.New("Input rate should be a positive, floating-point number")
	}

	whole, frac := math.Modf(1. / rate)
	for _, sub := range in.Subtitles {
		sub.Start = sub.Start*time.Duration(whole) + sub.Start*time.Duration(int(frac*1000))/1000
		sub.End = sub.End*time.Duration(whole) + sub.End*time.Duration(int(frac*1000))/1000
		res.Subtitles = append(res.Subtitles, sub)
	}
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
		return res, []error{errors.New("Something went wrong while trying to parse the provided file!")}
	}
	srt := string(content)
	_ = srt

	reader := new(SRTReader)
	reader.p, reader.s = 0, bufio.NewScanner(file)
	reader.s.Split(SRTScanner)

	for reader.s.Scan() {
		var current Subtitle
		cur := strings.FieldsFunc(reader.s.Text(), EOLSplit)

		idx, err1 := strconv.Atoi(cur[0])
		if err1 != nil {
			errCollection = append(errCollection, err1)
		} else {
			current.Index = idx
		}

		r1, _ := regexp.Compile(`\d+[,.:]\d+[,.:]\d+[,.:]\d+ -[ -]>`)
		r2, _ := regexp.Compile(`-[ -]> \d+[,.:]\d+[,.:]\d+[,.:]\d+`)

		start, err2 := TimestampToDurationSRT(r1.FindString(cur[1]))
		if err2 != nil {
			errCollection = append(errCollection, err2)
		} else {
			current.Start = start
		}

		end, err3 := TimestampToDurationSRT(r2.FindString(cur[1]))
		if err3 != nil {
			errCollection = append(errCollection, err3)
		} else {
			current.End = end
		}

		current.Content = strings.Join(cur[2:len(cur)], "\n")
		res.Subtitles = append(res.Subtitles, current)
	}

	return res, errCollection
}

func SRTScanner(data []byte, atEOF bool) (adv int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if i < len(data)-1 && string(data[i:i+2]) == "\n\n" {
			return i + 2, data[:i+2], nil
		}
		if i < len(data)-3 && string(data[i:i+4]) == "\r\n\r\n" {
			return i + 4, data[:i+4], nil
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

// SearchSubtitleFile scans the contents of all subtitle entries
// in a subtitle file for matches with the provided string,
// which can be a valid Regular Expression.
// It returns a slice of all entries that matched this input.
func SearchSubtitleFile(subfile SubtitleFile, in string) ([]Subtitle, error) {
	var res []Subtitle
	r, err := regexp.Compile(in)
	if err != nil {
		return res, errors.New("The provided search term is invalid :`" + in + "`")
	}
	for _, sub := range subfile.Subtitles {
		if r.MatchString(sub.Content) {
			res = append(res, sub)
		}
	}
	return res, nil
}

func SerializeSubtitles(subfile SubtitleFile) SubtitleFile {
	// This is a simple operation, that normally would take place
	// after parsing a subtitle files. Should we return some errors
	// to help with debugging and future expansion?
	var res SubtitleFile = subfile
	for i, _ := range subfile.Subtitles {
		res.Subtitles[i].Index = i + 1
	}
	return res
}

// DetectOverlaps scans the provided subtitle file serially,
// and checks consecutive subtitles for invalid start/end times
// Returns the pair or pairs of subtitles where overlaps where detected,
// one of which is the culprit
func DetectOverlaps(subfile SubtitleFile) []Subtitle {
	// Maybe in case of "longer" overlaps eg. sub 5 w/ sub 20, we should return the offending pair instead of the consecutive ones
	var overlaps []Subtitle
	for i := 0; i < len(subfile.Subtitles)-1; i++ {
		if subfile.Subtitles[i].End > subfile.Subtitles[i+1].Start {
			overlaps = append(overlaps, subfile.Subtitles[i], subfile.Subtitles[i+1])
		}
	}
	return overlaps
}

func ConvertToUTF8(filename string) {
}

func RemoveSubtitle(subfile SubtitleFile, idx int) (SubtitleFile, error) {
	// For now, we assume that the provided SubtitleFile is okay
	// and that the parser has taken care of any glaring issues

	// We need to use a copy, otherwise the slice will retain
	// references to the input slice, and will have side-effects
	// either we want to, or not
	res := make([]Subtitle, len(subfile.Subtitles))
	copy(res, subfile.Subtitles)
	if idx <= 0 || idx > len(subfile.Subtitles) {
		idxerr := strconv.Itoa(idx)
		return SubtitleFile{res, subfile.Headers}, errors.New("The index marked for removal is invalid :" + idxerr)
	}
	// Turn human input to zero-based index
	idx -= 1
	res = append(res[:idx], res[idx+1:]...)
	ret := SubtitleFile{res, subfile.Headers}
	ret = SerializeSubtitles(ret)
	return ret, nil
}

func AddSubtitle(subfile SubtitleFile, start, end, content, metadata, header string) (SubtitleFile, error) {
	var res SubtitleFile
	startTime, _ := StrToDuration(start)
	endTime, _ := StrToDuration(end)
	if startTime < 0 || endTime < 0 || endTime < startTime {
		return subfile, errors.New("Start and End times should be positive and ordered, ignoring input... " + start + " - " + end)
	}
	placed := false

	// added conditional if file is dead last or dead front
	// this is a bad practice, I think I can come up with something more elegant.
	// TODO TODO TODO TODO
	if startTime > subfile.Subtitles[len(subfile.Subtitles)-1].End {
		res.Subtitles = append(res.Subtitles, subfile.Subtitles...)
		res.Subtitles = append(res.Subtitles, Subtitle{len(subfile.Subtitles) + 1, startTime, endTime, content, metadata, header})
		res = SerializeSubtitles(res)
		return res, nil
	}

	if endTime < subfile.Subtitles[0].Start {
		res.Subtitles = append(res.Subtitles, Subtitle{1, startTime, endTime, content, metadata, header})
		res.Subtitles = append(res.Subtitles, subfile.Subtitles...)
		res = SerializeSubtitles(res)
		return res, nil
	}
	for i, sub := range subfile.Subtitles {
		if placed == false {
			res.Subtitles = append(res.Subtitles, sub)
		}

		if startTime > subfile.Subtitles[i].End && endTime < subfile.Subtitles[i+1].Start {
			placed = true
			// Bumped once for skipping current entry in loop, once for zero-based indexing
			res.Subtitles = append(res.Subtitles, Subtitle{i + 2, startTime, endTime, content, metadata, header})
			continue
		}

		if placed == true {
			// New index is n+2, one for the new entry, one for the zero-based indexing
			res.Subtitles = append(res.Subtitles, Subtitle{i + 2, sub.Start, sub.End, sub.Content, metadata, header})
		}
	}
	if placed == false {
		return subfile, errors.New("New subtitle would overlap with existing ones, ignoring it..." + start + " - " + end)
	}
	return res, nil
}

func PrintSubfileInfo(subfile SubtitleFile) {

	cpmLo, cpmHi, cpmAvg, runtime := 10000., 0., 0., 0.
	cpmLoIdx, cpmHiIdx, runtime := 0, 0, 0.
	for _, sub := range subfile.Subtitles {
		currentDur := (sub.End - sub.Start).Seconds()
		currentCpm := float64(len(sub.Content)) / currentDur
		if currentCpm > cpmHi {
			cpmHi = currentCpm
			cpmHiIdx = sub.Index
		}
		if currentCpm < cpmLo {
			cpmLo = currentCpm
			cpmLoIdx = sub.Index
		}
		cpmAvg += currentCpm
		runtime += currentDur
	}
	cpmAvg = cpmAvg / runtime

	fmt.Printf("Headers : %v\n", subfile.Headers)
	fmt.Printf("Number of subtitles : %d\n", len(subfile.Subtitles))
	fmt.Printf("Start Time : %v\n", subfile.Subtitles[0].Start)
	fmt.Printf("End Time : %v\n", subfile.Subtitles[len(subfile.Subtitles)-1].End)
	fmt.Printf("First-to-last Runtime : %v\n", (subfile.Subtitles[len(subfile.Subtitles)-1].End - subfile.Subtitles[0].Start))
	fmt.Printf("Subtitle Runtime : %v\n", runtime)
	fmt.Printf("Highest Character-Per-Minute : %.2f on subtitle index : %d\n", cpmHi, cpmHiIdx)
	fmt.Printf("Lowest Character-Per-Minute : %.2f on subtitle index : %d\n", cpmLo, cpmLoIdx)
	fmt.Printf("Average Character-Per-Minute : %.2f\n ", cpmAvg)
}
