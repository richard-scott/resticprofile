//+build darwin

package schedule

import (
	"bytes"
	"testing"

	"github.com/creativeprojects/resticprofile/calendar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"howett.net/plist"
)

func TestPListEncoderWithCalendarInterval(t *testing.T) {
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict><key>Day</key><integer>1</integer><key>Hour</key><integer>0</integer></dict></plist>`
	entry := newCalendarInterval()
	setCalendarIntervalValueFromType(entry, 1, calendar.TypeDay)
	setCalendarIntervalValueFromType(entry, 0, calendar.TypeHour)
	buffer := &bytes.Buffer{}
	encoder := plist.NewEncoder(buffer)
	err := encoder.Encode(entry)
	require.NoError(t, err)
	assert.Equal(t, expected, buffer.String())
}

func TestGetCalendarIntervalsFromScheduleTree(t *testing.T) {
	testData := []struct {
		input    string
		expected []CalendarInterval
	}{
		{"*-*-*", []CalendarInterval{
			{"Hour": 0, "Minute": 0},
		}},
		{"*:0,30", []CalendarInterval{
			{"Minute": 0},
			{"Minute": 30},
		}},
		{"0,12:20", []CalendarInterval{
			{"Hour": 0, "Minute": 20},
			{"Hour": 12, "Minute": 20},
		}},
		{"0,12:20,40", []CalendarInterval{
			{"Hour": 0, "Minute": 20},
			{"Hour": 0, "Minute": 40},
			{"Hour": 12, "Minute": 20},
			{"Hour": 12, "Minute": 40},
		}},
		{"Mon..Fri *-*-* *:0,30:00", []CalendarInterval{
			{"Weekday": 1, "Minute": 0},
			{"Weekday": 1, "Minute": 30},
			{"Weekday": 2, "Minute": 0},
			{"Weekday": 2, "Minute": 30},
			{"Weekday": 3, "Minute": 0},
			{"Weekday": 3, "Minute": 30},
			{"Weekday": 4, "Minute": 0},
			{"Weekday": 4, "Minute": 30},
			{"Weekday": 5, "Minute": 0},
			{"Weekday": 5, "Minute": 30},
		}},
	}

	for _, testItem := range testData {
		t.Run(testItem.input, func(t *testing.T) {
			event := calendar.NewEvent()
			err := event.Parse(testItem.input)
			assert.NoError(t, err)
			assert.ElementsMatch(t, testItem.expected, getCalendarIntervalsFromScheduleTree(generateTreeOfSchedules(event)))
		})
	}
}

func TestParseStatus(t *testing.T) {
	status := `{
	"StandardOutPath" = "local.resticprofile.self.check.log";
	"LimitLoadToSessionType" = "Aqua";
	"StandardErrorPath" = "local.resticprofile.self.check.log";
	"Label" = "local.resticprofile.self.check";
	"OnDemand" = true;
	"LastExitStatus" = 0;
	"Program" = "/Users/go/src/github.com/creativeprojects/resticprofile/resticprofile";
	"ProgramArguments" = (
		"/Users/go/src/github.com/creativeprojects/resticprofile/resticprofile";
		"--no-ansi";
		"--config";
		"examples/dev.yaml";
		"--name";
		"self";
		"check";
	);
};`
	expected := map[string]string{
		"StandardOutPath":        "local.resticprofile.self.check.log",
		"LimitLoadToSessionType": "Aqua",
		"StandardErrorPath":      "local.resticprofile.self.check.log",
		"Label":                  "local.resticprofile.self.check",
		"OnDemand":               "true",
		"LastExitStatus":         "0",
		"Program":                "/Users/go/src/github.com/creativeprojects/resticprofile/resticprofile",
	}

	output := parseStatus(status)
	assert.Equal(t, expected, output)
}
