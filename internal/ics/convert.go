package ics

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/HubertBel/lazyorg/internal/utils"
)

func findEventsRange(s []string) ([]int, []int, bool) {
	var rb = regexp.MustCompile(`(?m)\b(BEGIN:VEVENT)\b`)
	var re = regexp.MustCompile(`(?m)\b(END:VEVENT)\b`)
	b := []int{}
	e := []int{}
	for i, n := range s {
		if len(rb.FindString(n)) != 0 {
			b = append(b, i)
		}
		if len(re.FindString(n)) != 0 {
			e = append(e, i)
		}
	}
	return b, e, len(e) == len(b)
}

func getEvents(s string) []string {
	splited := strings.Split(s, "\n")
	b, e, same := findEventsRange(splited)
	if !same {
		panic("[ERROR] (from getEvents) The file is not formated properly")
	}
	events := []string{""}
	j := 0
	for i, v := range splited {
		if i > e[j] && j < len(e)-1 {
			events = append(events, "")
			j++
		} else if i > b[j] && i < e[j] {
			events[j] = fmt.Sprintf("%s\n%s", events[j], v)
		}
	}

	return events
}

func getEventTime(s []string, timezone *time.Location) (float64, string) {
	if strings.Contains(s[1], "VALUE=DATE") {
		date := strings.Split(s[1], ":")[1]
		t, err := time.Parse("20060102", strings.Trim(date, "\x0d"))
		if err != nil {
			panic(err)
		}
		return 1.0, fmt.Sprintf("%s@%s", utils.FormatDate(t), utils.FormatHourFromTime(t))
	} else {
		start := strings.Split(s[1], ":")[1]
		end := strings.Split(s[2], ":")[1]
		t0, err := time.Parse("20060102T150405Z", strings.Trim(start, "\x0d"))
		if err != nil {
			panic(err)
		}
		t1, err := time.Parse("20060102T150405Z", strings.Trim(end, "\x0d"))
		if err != nil {
			panic(err)
		}
		t0 = t0.In(timezone)
		t1 = t1.In(timezone)
		diff := t1.Sub(t0)
		r := fmt.Sprintf("%s@%s", utils.FormatDate(t0), utils.FormatHourFromTime(t0))
		return diff.Hours(), r
	}
}

func getEventNameDescription(s []string) (string, string) {
	n := ""
	d := ""
	for _, v := range s {
		if strings.Contains(v, "SUMMARY") {
			n = strings.Split(v, ":")[1]
			n = strings.Trim(n, "\x0d")
		} else if strings.Contains(v, "DESCRIPTION") {
			d = strings.Split(strings.SplitN(v, ":", 2)[1], "-::")[0]
			d = strings.Trim(d, "\x0d")
		}
	}

	return n, d
}

func getLocation(s []string) string {
	for _, v := range s {
		if strings.Contains(v, "LOCATION") {
			u := strings.Split(v, ":")[1]
			u = strings.Trim(u, "\x0d")
			u = strings.ReplaceAll(u, "\\", "")

			return u
		}
	}
	return ""
}

func ConvertIcs2LO(c []byte, tz int, isPCTimeZone bool) {
	s := string(c)
	v := getEvents(s)
	t := time.Now().Location()
	if isPCTimeZone {
		t = time.FixedZone("UTC-0", tz*60*60)
	}
	for _, e := range v {
		split := strings.Split(e, "\n")
		f, d := getEventTime(split, t)
		n, desc := getEventNameDescription(split)
		l := getLocation(split)

		o := fmt.Sprintf("%f|%s|%s|%s|%s", f, d, n, l, desc)
		fmt.Println(o)
	}
}
