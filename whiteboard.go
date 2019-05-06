package main

import (
	"fmt"
	"time"
)

func durationToTimestamp_w(d time.Duration) {
	// Some possible ways to convert Golang duration to SRT timestamp

	// Way 1 : Directly act on the duration object
	//days := int64(d.Hours() / 24)
	//hours := int64(math.Mod(d.Hours(), 24))
	//minutes := int64(math.Mod(d.Minutes(), 60))
	//seconds := int64(math.Mod(d.Seconds(), 60))
	//millis := int64(math.Mod(float64(d.Nanoseconds()), 1000)) * 1000
	//millis := float64(d) / float64(time.Millisecond)
	//fmt.Println(days, hours, minutes, seconds, millis)

	// Way 2 : Use Sscanf to 'scan' the formatted string from duration's `.String()` method
	//var hour, minute, second, milli int
	//fmt.Sscanf(s.String(), "%dh%dm%d.%ds", &hour, &minute, &second, &milli)
	//fmt.Println(hour, minute, second, milli)

	// basically
	//var hour, minute, second, milli int
	//var s float64
	//fmt.Sscanf(d.String(), "%dh%dm%fs", &hour, &minute, &s)
	//second = int(s)
	//milli = int(math.Mod(s*1000., 1000))
	//fmt.Println(hour, minute, second, milli)

	// Way 3 : Regular expressions
	// Now you have one more problem and whatnot

	// Way 4 : The incorrect one
	//var hour, minute int
	//var second float64
	//fmt.Println(d)
	//fmt.Sscanf(d.String(), "%dh%dm%fs", &hour, &minute, &second)
	//res := fmt.Sprintf("%02d:%02d:%02.3f", hour, minute, second)

	// Initially I tried to make this work.
	// While it DOES work for very specific formatted strings, it absolutely fails if you provide a duration
	// Which does not contain hours, or minutes, as the Sscanf puts values wherever it wants to.
	// Which also led me to see why testing is important

	fmt.Println("Done")
}
