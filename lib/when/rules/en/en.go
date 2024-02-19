package en

import "github.com/botlabs-gg/quackpdb/v2/lib/when/rules"

var All = []rules.Rule{
	Weekday(rules.Override),
	CasualDate(rules.Override),
	CasualTime(rules.Override),
	Hour(rules.Override),
	HourMinute(rules.Override),
	Deadline(rules.Override),
	PastTime(rules.Override),
	ExactMonthDate(rules.Override),
}

var WEEKDAY_OFFSET = map[string]int{
	"sunday":    0,
	"sun":       0,
	"monday":    1,
	"mon":       1,
	"tuesday":   2,
	"tue":       2,
	"wednesday": 3,
	"wed":       3,
	"thursday":  4,
	"thur":      4,
	"thu":       4,
	"friday":    5,
	"fri":       5,
	"saturday":  6,
	"sat":       6,
}

var WEEKDAY_OFFSET_PATTERN = "(?:sunday|sun|monday|mon|tuesday|tue|wednesday|wed|thursday|thur|thu|friday|fri|saturday|sat)"

var MONTH_OFFSET = map[string]int{
	"january":   1,
	"jan":       1,
	"jan.":      1,
	"february":  2,
	"feb":       2,
	"feb.":      2,
	"march":     3,
	"mar":       3,
	"mar.":      3,
	"april":     4,
	"apr":       4,
	"apr.":      4,
	"may":       5,
	"june":      6,
	"jun":       6,
	"jun.":      6,
	"july":      7,
	"jul":       7,
	"jul.":      7,
	"august":    8,
	"aug":       8,
	"aug.":      8,
	"september": 9,
	"sep":       9,
	"sep.":      9,
	"sept":      9,
	"sept.":     9,
	"october":   10,
	"oct":       10,
	"oct.":      10,
	"november":  11,
	"nov":       11,
	"nov.":      11,
	"december":  12,
	"dec":       12,
	"dec.":      12,
}

var MONTH_OFFSET_PATTERN = `(?:january|jan\.?|february|feb\.?|march|mar\.?|april|apr\.?|may|june|jun\.?|july|jul\.?|august|aug\.?|september|sept?\.?|october|oct\.?|november|nov\.?|december|dec\.?)`

var INTEGER_WORDS = map[string]int{
	"one":    1,
	"two":    2,
	"three":  3,
	"four":   4,
	"five":   5,
	"six":    6,
	"seven":  7,
	"eight":  8,
	"nine":   9,
	"ten":    10,
	"eleven": 11,
	"twelve": 12,
}

var INTEGER_WORDS_PATTERN = `(?:one|two|three|four|five|six|seven|eight|nine|ten|eleven|twelve)`

var ORDINAL_WORDS = map[string]int{
	"first":          1,
	"1st":            1,
	"second":         2,
	"2nd":            2,
	"third":          3,
	"3rd":            3,
	"fourth":         4,
	"4th":            4,
	"fifth":          5,
	"5th":            5,
	"sixth":          6,
	"6th":            6,
	"seventh":        7,
	"7th":            7,
	"eighth":         8,
	"8th":            8,
	"ninth":          9,
	"9th":            9,
	"tenth":          10,
	"10th":           10,
	"eleventh":       11,
	"11th":           11,
	"twelfth":        12,
	"12th":           12,
	"thirteenth":     13,
	"13th":           13,
	"fourteenth":     14,
	"14th":           14,
	"fifteenth":      15,
	"15th":           15,
	"sixteenth":      16,
	"16th":           16,
	"seventeenth":    17,
	"17th":           17,
	"eighteenth":     18,
	"18th":           18,
	"nineteenth":     19,
	"19th":           19,
	"twentieth":      20,
	"20th":           20,
	"twenty first":   21,
	"twenty-first":   21,
	"21st":           21,
	"twenty second":  22,
	"twenty-second":  22,
	"22nd":           22,
	"twenty third":   23,
	"twenty-third":   23,
	"23rd":           23,
	"twenty fourth":  24,
	"twenty-fourth":  24,
	"24th":           24,
	"twenty fifth":   25,
	"twenty-fifth":   25,
	"25th":           25,
	"twenty sixth":   26,
	"twenty-sixth":   26,
	"26th":           26,
	"twenty seventh": 27,
	"twenty-seventh": 27,
	"27th":           27,
	"twenty eighth":  28,
	"twenty-eighth":  28,
	"28th":           28,
	"twenty ninth":   29,
	"twenty-ninth":   29,
	"29th":           29,
	"thirtieth":      30,
	"30th":           30,
	"thirty first":   31,
	"thirty-first":   31,
	"31st":           31,
}

var ORDINAL_WORDS_PATTERN = `(?:1st|first|2nd|second|3rd|third|4th|fourth|5th|fifth|6th|sixth|7th|seventh|8th|eighth|9th|ninth|10th|tenth|11th|eleventh|12th|twelfth|13th|thirteenth|14th|fourteenth|15th|fifteenth|16th|sixteenth|17th|seventeenth|18th|eighteenth|19th|nineteenth|20th|twentieth|21st|twenty[ -]first|22nd|twenty[ -]second|23rd|twenty[ -]third|24th|twenty[ -]fourth|25th|twenty[ -]fifth|26th|twenty[ -]sixth|27th|twenty[ -]seventh|28th|twenty[ -]eighth|29th|twenty[ -]ninth|30th|thirtieth|31st|thirty[ -]first)`
