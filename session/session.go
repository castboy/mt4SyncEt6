package session

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 0:0-12:0 -> (21:0-24:0, 0:0-9:0)
func sub3Hour(timeSpan string) (preDay string, thisDay string) {
	if timeSpan == "0:0-0:0" {
		return "", ""
	}

	t := strings.Split(timeSpan, "-")
	from := t[0]
	to := t[1]

	fromHM := strings.Split(from, ":")
	fromHourS := fromHM[0]

	toHM := strings.Split(to, ":")
	toHourS := toHM[0]

	fromHour, err := strconv.Atoi(fromHourS)
	if err != nil {
		panic(err)
	}

	toHour, err := strconv.Atoi(toHourS)
	if err != nil {
		panic(err)
	}

	hourLen := 24

	fromHourSub3 := (fromHour + hourLen -3) % hourLen
	toHourSub3 := (toHour + hourLen -3) % hourLen

	if fromHourSub3 < fromHour {
		return "", fmt.Sprintf("%d:%s-%d:%s", fromHourSub3, fromHM[1], toHourSub3, toHM[1])
	}

	if toHourSub3 > toHour {
		return fmt.Sprintf("%d:%s-%d:%s", fromHourSub3, fromHM[1], toHourSub3, toHM[1]), ""
	}

	return fmt.Sprintf("%d:%s-%d:%s", fromHourSub3, fromHM[1], 24, toHM[1]), fmt.Sprintf("%d:%d-%d:%s", 0, 0, toHourSub3, toHM[1])
}

// 20:0-21:0 21:0-24:0 -> 20:0-24:0
func mergeAdjacent(timeSpans []string, id int) []string {
	if len(timeSpans) <= 1 || id + 1 >= len(timeSpans) {
		return timeSpans
	}

	first := strings.Split(timeSpans[id], "-")
	second := strings.Split(timeSpans[id+1], "-")

	if first[1] == second[0] {
		var tail []string
		if id+2 < len(timeSpans) {
			tail = timeSpans[id+2:]
		}

		timeSpans = append(timeSpans[:id], first[0] + "-" + second[1])
		if len(tail) != 0 {
			timeSpans = append(timeSpans, tail...)
		}

		return mergeAdjacent(timeSpans, 0)
	}

	return mergeAdjacent(timeSpans, id+1)
}


// "0:0" -> "00:00"
func toStandardTime(mt4 string) (et6 string) {
	f := func(s string) string {
		if len(s) == 1 {
			return "0" + s
		}

		return s
	}

	t := strings.Split(mt4, ":")

	return f(t[0]) + ":" + f(t[1])
}

func toEt6String(timeSpans []string) {
	expr := `(\d+:\d+)`

	for i := range timeSpans {
		re, err := regexp.Compile(expr)
		if err != nil {
			panic(err)
		}

		timeSpans[i] = re.ReplaceAllStringFunc(timeSpans[i], toStandardTime)
	}
}

func mt4ToEt6(src map[string]string) map[string][]string {
	weekLen := 7
	pre := map[string][]string{}
	dst := map[string][]string{}

	for i := range src {
		timeSpans := strings.Split(src[i], ",")
		for j := range timeSpans {
			preDay, thisDay := sub3Hour(timeSpans[j])
			if preDay != "" && preDay != "0:0-0:0" {
				day, err := strconv.Atoi(i)
				if err != nil {
					panic(err)
				}

				pr := strconv.Itoa((day+weekLen-1)%weekLen)
				pre[pr] = append(pre[pr], preDay)
			}

			if thisDay != "" && thisDay != "0:0-0:0" {
				day, err := strconv.Atoi(i)
				if err != nil {
					panic(err)
				}

				da := strconv.Itoa(day)
				dst[da] = append(dst[da], thisDay)
			}
		}
	}

	for i := 0; i < weekLen; i++ {
		day := strconv.Itoa(i)
		dst[day] = append(dst[day], pre[day]...)
	}

	return dst
}

func Mt4ToEt6(mt4 map[string]string) map[string][]string {
	et6 := mt4ToEt6(mt4)

	dst := map[string][]string{}

	dst["0"] = mergeAdjacent(et6["0"], 0)
	dst["1"] = mergeAdjacent(et6["1"], 0)
	dst["2"] = mergeAdjacent(et6["2"], 0)
	dst["3"] = mergeAdjacent(et6["3"], 0)
	dst["4"] = mergeAdjacent(et6["4"], 0)
	dst["5"] = mergeAdjacent(et6["5"], 0)
	dst["6"] = mergeAdjacent(et6["6"], 0)

	for i := range dst {
		toEt6String(dst[i])
	}

	return dst
}

//func main() {
//
//	// et6 = mt4 - 3
//	et6 := mt4ToEt6(mt4)
//
//	dst := map[string][]string{}
//
//	dst["0"] = mergeAdjacent(et6["0"], 0)
//	dst["1"] = mergeAdjacent(et6["1"], 0)
//	dst["2"] = mergeAdjacent(et6["2"], 0)
//	dst["3"] = mergeAdjacent(et6["3"], 0)
//	dst["4"] = mergeAdjacent(et6["4"], 0)
//	dst["5"] = mergeAdjacent(et6["5"], 0)
//	dst["6"] = mergeAdjacent(et6["6"], 0)
//
//	fmt.Println(dst)
//
//	for i := range dst {
//		toEt6String(dst[i])
//	}
//
//
//	// et6 need currently.
//	fmt.Println(dst)
//}





