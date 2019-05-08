

# whiteboard.md

I like to have a 'whiteboard' for my projects.

It gets really messy, real quick, but provides a way to write down my thoughts, look back on previous design decisions, reason on future approaches, jot down various tests I've done and failed.

This is done best with pen-and-paper, but since I'm using Git, I'll try to go for digital.


## Design

Since I've now started with writing functionality, I'm not yet sure what will be the best way to package the whole thing, as a library, to enable it to be used as a cli application and whatnot. 

Because I'm lazy, I'll be throwing everything into `main.go`, and the structure will emerge later on.


I am not sure how to handle whole subtitle files.    
I am pretty confident with each subtitle being a struct like this will do, (maybe add another `metadata` field).
```go
Subtitle struct {
	Index    int
	Start    time.Duration
	End      time.Duration
	Content string
}
```

I am still on the fence about the whole file though. It can take a bunch of different forms.
```go
var subfile []Subtitle

var subfile map[int]subtitle

type subfile struct {
		subs []Subtitle
		metadata string
}
```

I think I might go with the KISS principle and just use a slice.
It's ordered, can be easily iterated, avoid nested structs, easy for others to use as API in the future.

## Parsing SRT File
I've been postponing this, but I really need to sit down and find a nice solution.

## SQL Queries
One of the reasons I'm excited about this project, and keep pushing it forwards, is this feature, the ability to run SQL queries on multiple subtitle files.

I'll probably be using SQLite, to create a singular `.db` file that I'll be querying. What's a good schema supposed to be like for this job?

```
filename  | subtitle_index | subtitle_start | subtitle_end | subtitle_content
```



## Testing


## Functions

### func ParseSRTFile()
I'm thinking about how to best parse SRT files. I see some options (not too many, for now). 

- I could use multiline regex matching, 
- Use a regex to match a subset and then expand forwards-backwards (like search for `-->` and/or `\n\n%d` to locate the subtitles
- Work directly on a line-by-line and detect two empty lines in succession.

I need to run some test with a video player and toy around with what's allowed and what not, eg. what errors will VLC ignore and what errors it will fail at displaying.

For now the last solution seems to be the simplest one, but I don't know how it will deal with errors. In case we go forward with regular expressions, we might need multiple alternative regexes to 'catch' error cases

### func durationToTimestamp()

There are a bunch of possible ways to conver a `time.Duration` to an SRT timestamp. If I were to write this in Python, I'd be using something like `datetime.timedelta`, so I think of `time.Duration` as something analogous

	// Some possible ways to convert Golang duration to SRT timestamp


Way 1 - Act directly on the duration object.   
Why it's bad : 
```go
days := int64(d.Hours() / 24)
hours := int64(math.Mod(d.Hours(), 24))
minutes := int64(math.Mod(d.Minutes(), 60))
seconds := int64(math.Mod(d.Seconds(), 60))
millis := int64(math.Mod(float64(d.Nanoseconds()), 1000)) * 1000
millis := float64(d) / float64(time.Millisecond)
fmt.Println(days, hours, minutes, seconds, millis)
```

Way 2 : Use Sscanf to 'scan' the formatted string from the object's `.String()` method    
Whu it's bad :
```go
var hour, minute, second, milli int
fmt.Sscanf(s.String(), "%dh%dm%d.%ds", &hour, &minute, &second, &milli)
fmt.Println(hour, minute, second, milli)

// basically
var hour, minute, second, milli int
var s float64
fmt.Sscanf(d.String(), "%dh%dm%fs", &hour, &minute, &s)
second = int(s)
milli = int(math.Mod(s*1000., 1000))
fmt.Println(hour, minute, second, milli)
```

Way 3 : Regular expressions    
Now you have one more problem and whatnot


Way 4 : The way I made it initially work. BUT   
it works for *very specific* formatted strings. It absolutely fails if e.g. you provide a duration that does not contain hours or minutes, and `Sscanf` just puts the first value it finds (s) in the 'hour' bucket.

This was one of the cases that was immediately obvious once I wrote some tests, and pointed me to keep on writing them and not get lazy.
```go
var hour, minute int
var second float64
fmt.Println(d)
fmt.Sscanf(d.String(), "%dh%dm%fs", &hour, &minute, &second)
res := fmt.Sprintf("%02d:%02d:%02.3f", hour, minute, second)
```

Way 5 : The way it currently works     
Using the native method of parsing the object and acting on its components is the easiest and most concise method I've found.
*It also allows for future improvements*, such as converting to another format of timestamps, for other subtitle files.
```go
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
```

### func StrToDuration()
I probably could get away without testing this function, but at least I found out more about error handling, and about writing tests for specific errors.

I'm using the wrapper instead of going full in with `time.ParseDuration`, to be able to extend it later when it's going to take cli arguments as well.

