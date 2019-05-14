# go-subtitles

Development started on 2019-05-03 and was spurred as an idea a few days earlier, as I was struggling with the subtitles of the lastest episode in the "A Symphony of Frost and Flame" series.


## Core Ideas

- The initial version will work with `.srt` files, which is the most widespread subtitle format
- The main audience for subtitles are hearing impaired people and people not perfectly fluent in english. UTF-8 support and/or support for different symbols in various languages is very important

## First wishlist
During my first brainstorming sessions, here's what I imagine a program or library that deals with subtitles be able to achieve.    
I really need to talk with other people who actually *write subs* and check out *their* preferences, requirements and reality.

Feel free to open an issue, or create a pull request for anything you might want to add to the list.

- [x] Parse and Write to SRT Format
- [ ] Encode subtitles in different formats, change/preview their encoding
- [x] Add/Remove subtitles
- [ ] Modify subtitles
- [ ] Synchronize subtitles by adding-removing time from the whole file or a specific section (and then add audio-detection so it's done automatically)
- [x] Change subtitle duration in either *relative* or *absolute* time
- [ ] Search-and-replace subtitle text strings
- [x] Find overlapping subtitles
- [ ] Re-index (and re-sort) subtitles based on start times
- [ ] Auto report problems in subtitles (malformed files, non-sequential entries, and whatnot)
- [ ] Run SQL queries in one or more subtitles that exist in a directory
- [ ] Facilitate translating using side-to-side panes
- [ ] Simply view subtitle files (should be better than a text editor)
- [ ] Add or disable subtitle colors and other subtitle effects
- [ ] Live preview of changes (maybe tied-in with VLC or something) i.e. open video file with current sub and jump to specific time
- [ ] Diff two subtitle files
- [ ] Convert to/from other subtitle formats
- [ ] Hardcode subs to videos 
- [ ] Search for and download subtitles for your video automatically
- [ ] Help create subtitles for hearing impaired people (maybe by facilitating the addition of "tags" such as [LOUD MUSIC])
- [ ] Estimate subtitle "speed" (maybe characters-per-second?) and check for too-long or too-short ones.
- [ ]   
- [ ]    
- [ ] What are *your* ideas?
