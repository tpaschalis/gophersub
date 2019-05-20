package main

import (
	"bufio"
	"bytes"
	"errors"
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

func TimestampSplitSRT(r rune) bool {
	return r == ',' || r == '.' || r == ':'
}

func EOLSplit(r rune) bool {
	return r == '\n' || r == '\r'
}
