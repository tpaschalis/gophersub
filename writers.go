package main

import (
	"bufio"
	"errors"
	"os"
	"strconv"
)

// Exports a SubtitleFile object to an SRT file format.
// If the file exists, no error will be raised, but will
// append the data towards its end.
func ToSRTFile(subfile SubtitleFile, outfile string) error {

	f, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return errors.New("Could not open file " + outfile + " for writing")
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	// SRT Files do not feature header or metadata information,
	// so this information will not be written to the file
	// https://matroska.org/technical/specs/subtitles/srt.html

	var idxStr, startStr, endStr string
	for _, sub := range subfile.Subtitles {
		startStr = DurationToTimestampSRT(sub.Start)
		endStr = DurationToTimestampSRT(sub.End)
		idxStr = strconv.Itoa(sub.Index)
		_, err = w.WriteString(idxStr + "\n")
		_, err = w.WriteString(startStr + " --> " + endStr + "\n")
		_, err = w.WriteString(sub.Content)
		_, err = w.WriteString("\n\n")
		w.Flush()
	}
	//w.WriteString("\n")
	w.Flush()

	return nil
}
