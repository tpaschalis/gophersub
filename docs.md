Technical specs

- [SRT Subtitles](https://matroska.org/technical/specs/subtitles/srt.html)

SRT is perhaps the most basic of all subtitle formats.

It consists of four parts, all in text..

1. A number indicating which subtitle it is in the sequence.
2. The time that the subtitle should appear on the screen, and then disappear.
3. The subtitle itself.
4. A blank line indicating the start of a new subtitle.

When placing SRT in Matroska, part 3 is converted to UTF-8 (S_TEXT/UTF8) and placed in the data portion of the Block. Part 2 is used to set the timecode of the Block, and BlockDuration element. Nothing else is used.

Here is an example SRT file:
```
1
00:02:17,440 --> 00:02:20,375
Senator, we're making
our final approach into Coruscant.

2
00:02:20,476 --> 00:02:22,501
Very good, Lieutenant.
```
 
In this example, the text "Senator, we're making our final approach into Coruscant." would be converted into UTF-8 and placed in the Block. The timecode of the block would be set to "00:02:17,440". And the BlockDuration element would be set to "00:00:02,935".

The same is repeated for the next subtitle.

Because there are no general settings for SRT, the CodecPrivate is left blank.


