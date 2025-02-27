package br

import "github.com/SoggySaussages/syzygy/lib/when/rules"

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
	"domingo":       0,
	"dom":           0,
	"segunda-feira": 1,
	"segunda":       1,
	"seg":           1,
	"terça-feira":   2,
	"terça":         2,
	"ter":           2,
	"quarta-feira":  3,
	"quarta":        3,
	"qua":           3,
	"quinta-feira":  4,
	"quinta":        4,
	"qui":           4,
	"sexta-feira":   5,
	"sexta":         5,
	"sex":           5,
	"sábado":        6,
	"sab":           6,
}

var WEEKDAY_OFFSET_PATTERN = "(?:domingo|dom|segunda-feira|segunda|seg|terça-feira|terça|ter|quarta-feira|quarta|qua|quinta-feira|quinta|qui|sexta-feira|sexta|sex|sábado|sab)"

var MONTH_OFFSET = map[string]int{
	"janeiro":   1,
	"jan.":      1,
	"jan":       1,
	"fevereiro": 2,
	"fev.":      2,
	"fev":       2,
	"março":     3,
	"mar.":      3,
	"mar":       3,
	"abril":     4,
	"abr.":      4,
	"abr":       4,
	"maio":      5,
	"mai.":      5,
	"mai":       5,
	"junho":     6,
	"jun.":      6,
	"jun":       6,
	"julho":     7,
	"jul.":      7,
	"jul":       7,
	"agosto":    8,
	"ago.":      8,
	"ago":       8,
	"setembro":  9,
	"set.":      9,
	"set":       9,
	"outubro":   10,
	"out.":      10,
	"out":       10,
	"novembro":  11,
	"nov.":      11,
	"nov":       11,
	"dezembro":  12,
	"dez.":      12,
	"dez":       12,
}

var MONTH_OFFSET_PATTERN = `(?:janeiro|jan\.?|jan|fevereiro|fev\.?|fev|março|mar\.?|mar|abril|abr\.?|abr|maio|mai\.?|mai|junho|jun\.?|jun|julho|jul\.?|jul|agosto|ago\.?|ago|setembro|set\.?|set|outubro|out\.?|out|novembro|nov\.?|nov|dezembro|dez\.?|dez)`

var INTEGER_WORDS = map[string]int{
	"uma":    1,
	"um":     1,
	"duas":   2,
	"dois":   2,
	"três":   3,
	"quatro": 4,
	"cinco":  5,
	"seis":   6,
	"sete":   7,
	"oito":   8,
	"nove":   9,
	"dez":    10,
	"onze":   11,
	"doze":   12,
}

var INTEGER_WORDS_PATTERN = `(?:uma|um|duas|dois|três|quatro|cinco|seis|sete|oito|nove|dez|onze|doze)`

var ORDINAL_WORDS = map[string]int{
	"primeiro":          1,
	"1º":                1,
	"segunda":           2,
	"segundo":           2,
	"2ª":                2,
	"2º":                2,
	"terceira":          3,
	"terceiro":          3,
	"3ª":                3,
	"3º":                3,
	"quarta":            4,
	"quarto":            4,
	"4ª":                4,
	"4º":                4,
	"quinta":            5,
	"quinto":            5,
	"5ª":                5,
	"5º":                5,
	"sexta":             6,
	"sexto":             6,
	"6ª":                6,
	"6º":                6,
	"sétima":            7,
	"sétimo":            7,
	"7ª":                7,
	"7º":                7,
	"oitava":            8,
	"oitavo":            8,
	"8ª":                8,
	"8º":                8,
	"nona":              9,
	"nono":              9,
	"9ª":                9,
	"9º":                9,
	"décima":            10,
	"décimo":            10,
	"10ª":               10,
	"10º":               10,
	"décima-primeira":   11,
	"décima primeira":   11,
	"décimo-primeiro":   11,
	"décimo primeiro":   11,
	"11ª":               11,
	"11º":               11,
	"décima-segunda":    12,
	"décima segunda":    12,
	"décimo-segundo":    12,
	"décimo segundo":    12,
	"12ª":               12,
	"12º":               12,
	"décima-terceira":   13,
	"décima terceira":   13,
	"décimo-terceiro":   13,
	"décimo terceiro":   13,
	"13ª":               13,
	"13º":               13,
	"décima-quarta":     14,
	"décima quarta":     14,
	"décimo-quarto":     14,
	"décimo quarto":     14,
	"14ª":               14,
	"14º":               14,
	"décima-quinta":     15,
	"décima quinta":     15,
	"décimo-quinto":     15,
	"décimo quinto":     15,
	"15ª":               15,
	"15º":               15,
	"décima-sexta":      16,
	"décima sexta":      16,
	"décimo-sexto":      16,
	"décimo sexto":      16,
	"16ª":               16,
	"16º":               16,
	"décima-sétima":     17,
	"décima sétima":     17,
	"décimo-sétimo":     17,
	"décimo sétimo":     17,
	"17ª":               17,
	"17º":               17,
	"décima-oitava":     18,
	"décima oitava":     18,
	"décimo-oitavo":     18,
	"décimo oitavo":     18,
	"18ª":               18,
	"18º":               18,
	"décima-nona":       19,
	"décima nona":       19,
	"décimo-nono":       19,
	"décimo nono":       19,
	"19ª":               19,
	"19º":               19,
	"vigésima":          20,
	"vigésimo":          20,
	"20ª":               20,
	"20º":               20,
	"vigésima-primeira": 21,
	"vigésima primeira": 21,
	"vigésimo-primeiro": 21,
	"vigésimo primeiro": 21,
	"21ª":               21,
	"21º":               21,
	"vigésima-segunda":  22,
	"vigésima segunda":  22,
	"vigésimo-segundo":  22,
	"vigésimo segundo":  22,
	"22ª":               22,
	"22º":               22,
	"vigésima-terceira": 23,
	"vigésima terceira": 23,
	"vigésimo-terceiro": 23,
	"vigésimo terceiro": 23,
	"23ª":               23,
	"23º":               23,

	"vigésima-quarta":    24,
	"vigésima quarta":    24,
	"vigésimo-quarto":    24,
	"vigésimo quarto":    24,
	"24ª":                24,
	"24º":                24,
	"vigésima-quinta":    25,
	"vigésima quinta":    25,
	"vigésimo-quinto":    25,
	"vigésimo quinto":    25,
	"25ª":                25,
	"25º":                25,
	"vigésima-sexta":     26,
	"vigésima sexta":     26,
	"vigésimo-sexto":     26,
	"vigésimo sexto":     26,
	"26ª":                26,
	"26º":                26,
	"vigésima-sétima":    27,
	"vigésima sétima":    27,
	"vigésimo-sétimo":    27,
	"vigésimo sétimo":    27,
	"27ª":                27,
	"27º":                27,
	"vigésima-oitava":    28,
	"vigésima oitava":    28,
	"vigésimo-oitavo":    28,
	"vigésimo oitavo":    28,
	"28ª":                28,
	"28º":                28,
	"vigésima-nona":      29,
	"vigésima nona":      29,
	"vigésimo-nono":      29,
	"vigésimo nono":      29,
	"29ª":                29,
	"29º":                29,
	"trigésima":          30,
	"trigésimo":          30,
	"30ª":                30,
	"30º":                30,
	"trigésima-primeira": 31,
	"trigésima primeira": 31,
	"trigésimo-primeiro": 31,
	"trigésimo primeiro": 31,
	"31ª":                31,
	"31º":                31,
}

var ORDINAL_WORDS_PATTERN = `(?:primeir[ao]|1[ªº]|segund[ao]|2[ªº]|terceir[ao]|3[ªº]|quart[ao]|4[ªº]|quint[ao]|5[ªº]|sext[ao]|6[ªº]|sétim[ao]|7[ªº]|oitav[ao]|8[ªº]|non[ao]|9[ªº]|décim[ao]|10[ªº]|décim[ao][- ]primeir[ao]|11[ªº]|décima[ao][- ]segund[ao]|12[ªº]|décim[ao][- ]terceir[ao]|13[ªº]|décim[ao][- ]quart[ao]|14[ªº]|décim[ao][- ]quint[ao]|15[ªº]|décim[ao][- ]sext[ao]|16[ªº]|décim[ao][- ]sétim[ao]|17[ªº]|décim[ao][- ]oitav[ao]|18[ªº]|décim[ao][- ]non[ao]|19[ªº]|vigésim[ao]|20[ªº]|vigésim[ao][- ]primeir[ao]|21[ªº]|vigésim[ao][- ]segund[ao]|22[ªº]|vigésim[ao][- ]terceir[ao]|23[ªº]|vigésim[ao][- ]quart[ao]|24[ªº]|vigésim[ao][- ]quint[ao]|25[ªº]|vigésim[ao][- ]sext[ao]|26[ªº]|vigésim[ao][- ]sétim[ao]|27[ªº]|vigésim[ao][- ]oitav[ao]|28[ªº]|vigésim[ao][- ]non[ao]|29[ªº]|trigésim[ao]|30[ªº]|trigésim[ao][- ]primeir[ao]|31[ªº]|31º)`
