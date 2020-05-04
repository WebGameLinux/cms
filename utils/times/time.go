package times

import "time"

var days = [][]string{
		{"Sunday", "Sun", "sunday", "sun", "周日", "7"},
		{"Monday", "Mon", "monday", "mon", "周一", "1"},
		{"Tuesday", "Tues", "tuesday", "tues", "周二", "2"},
		{"Wednesday", "Wed", "wednesday", "wed", "周三", "3"},
		{"Thursday", "Thurs", "thursday", "thurs", "周四", "4"},
		{"Friday", "Fri", "friday", "fri", "周五", "5"},
		{"Saturday", "Satur", "saturday", "satur", "周六", "6"},
}

var layouts = []string{
		"2006-01-02 ",
		"2006 01 02",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		// Handy time stamps.
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
}

func ParseWeekday(str string) (time.Weekday, bool) {
		for i, arr := range days {
				if len(arr) == 0 {
						continue
				}
				for _, v := range arr {
						if v == str {
								return time.Weekday(i), true
						}
				}
		}
		if t, ok := ParseTime(str); ok {
				return t.Weekday(), true
		}
		return time.Monday, false
}

func ParseTime(str string) (time.Time, bool) {
		for _, layout := range layouts {
				if t, err := time.Parse(layout, str); err == nil {
						return t, true
				}
		}
		return time.Time{}, false
}
