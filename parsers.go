package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type SRTReader struct {
	s *bufio.Scanner
	b bytes.Buffer
	p int
}

type WebVTTReader struct {
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

		current.Content = strings.Join(cur[2:], "\n")
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

func WebVTTScanner(data []byte, atEOF bool) (adv int, token []byte, err error) {
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

func ParseWebVTTFile(filename string) (SubtitleFile, []error) {
	// One source of inspiration, ehm, theft
	// https://github.com/glut23/webvtt-py/blob/master/webvtt/parsers.py
	//CUE_TIMING_PATTERN, _ := regexp.Compile(`\s*((?:\d+:)?\d{2}:\d{2}.\d{3})\s*-->\s*((?:\d+:)?\d{2}:\d{2}.\d{3})`)
	//COMMENT_PATTERN, _ := regexp.Compile(`NOTE(?:\s.+|$)`)
	//STYLE_PATTERN, _ := regexp.Compile(`STYLE[ \t]*$`)

	var res SubtitleFile
	var errCollection []error

	content, err := ioutil.ReadFile(filename)
	file, err := os.Open(filename)

	if err != nil {
		return res, []error{errors.New("Something went wrong while trying to parse the provided WebVTT file!")}
	}
	srt := string(content)
	_ = srt

	reader := new(WebVTTReader)
	reader.p, reader.s = 0, bufio.NewScanner(file)
	reader.s.Split(WebVTTScanner)

	currentBlock := 0
	for reader.s.Scan() {
		//var current Subtitle
		cur := strings.FieldsFunc(reader.s.Text(), EOLSplit)
		fmt.Println(cur)
		if blockno == 0 && cur[0] != "WEBVTT" {
			return res, []error{errors.New("Provided file does not start with `WEBVTT`")}
		}

		//idx, err1 := strconv.Atoi(cur[0])
		//if err1 != nil {
		//	errCollection = append(errCollection, err1)
		//} else {
		//	current.Index = idx
		//}

		//r1, _ := regexp.Compile(`\d+[,.:]\d+[,.:]\d+[,.:]\d+ -[ -]>`)
		//r2, _ := regexp.Compile(`-[ -]> \d+[,.:]\d+[,.:]\d+[,.:]\d+`)

		//start, err2 := TimestampToDurationSRT(r1.FindString(cur[1]))
		//if err2 != nil {
		//	errCollection = append(errCollection, err2)
		//} else {
		//	current.Start = start
		//}

		//end, err3 := TimestampToDurationSRT(r2.FindString(cur[1]))
		//if err3 != nil {
		//	errCollection = append(errCollection, err3)
		//} else {
		//	current.End = end
		//}

		//current.Content = strings.Join(cur[2:], "\n")
		//res.Subtitles = append(res.Subtitles, current)
		currentBlock += 1
	}

	return res, errCollection
}

// Returns true if input slice of lines
// matches a 'cue' block, i.e. one of the
// first two lines contains a cue timing
func isCueBlockVTT(block []string) bool {
	CUE_TIMING_PATTERN, _ := regexp.Compile(`\s*((?:\d+:)?\d{2}:\d{2}.\d{3})\s*-->\s*((?:\d+:)?\d{2}:\d{2}.\d{3})`)
	return CUE_TIMING_PATTERN.Match([]byte(block[0])) || CUE_TIMING_PATTERN.Match([]byte(block[1]))
}

// Returns true if input slice of lines
// matches a 'comment' block, i.e.
// contains the NOTE directive
func isCommentBlockWebVTT(block []string) bool {
	COMMENT_PATTERN, _ := regexp.Compile(`NOTE(?:\s.+|$)`)
	return COMMENT_PATTERN.Match([]byte(block[0]))
}

// Returns true if input slice of lines
// matches a 'style' block i.e.
// contains the STYLE directive
func isStyleBlockWebVTT(block []string) bool {
	STYLE_PATTERN, _ := regexp.Compile(`STYLE[ \t]*$`)
	return STYLE_PATTERN.Match([]byte(block[0]))
}
