package utils

import "testing"

func TestDate(t *testing.T) {
	t.Logf("Get current time: %v.", GetCurrentTime())
	t.Logf("Get current timestamp: %d.", GetCurrentTimestamp())
	t.Logf("Get current milli timestamp: %d.", GetCurrentMilliTimestamp())
	t.Logf("Get current naco timestamp: %d.", GetCurrentNanoTimestamp())
	t.Logf("Get current naco: %d.", GetCurrentNano())
	t.Logf("Get current year: %d.", GetCurrentYear())
	t.Logf("Get current month: %s.", GetCurrentMonth())
	t.Logf("Get current weekday: %d.", GetCurrentWeekday())
	t.Logf("Get current day: %d.", GetCurrentDay())

	year, month, day := GetCurrentYearMonthDay()
	t.Logf("Get current year: %d, month: %d, day: %d.", year, month, day)

	t.Logf("Get current hour: %d.", GetCurrentHour())
	t.Logf("Get current minute: %d.", GetCurrentMinute())
	t.Logf("Get current second: %d.", GetCurrentSecond())
	t.Logf("Get current day of year: %d.", GetCurrentDayOfYear())
	t.Logf("Get current day of date: %s.", GetCurrentDate())
	t.Logf("Get current day of date until second: %s.", GetCurrentDateUntilSecond())
	t.Logf("Get current day of date until mill second: %s.", GetCurrentDateUntilMilli())
	t.Logf("Get current day of date until nano second: %s.", GetCurrentDateUntilNano())

	defaultStamp, err := DateUntilSecondToTimeStamp("2018-07-11 15:07:51", "")
	if err != nil {
		t.Errorf("DateUntilSecondToTimeStamp default err: %v.", err)
	}
	t.Logf("DateUntilSecondToTimeStamp default: %d.", defaultStamp)

	stamp, err := DateUntilSecondToTimeStamp("2018-07-11 15:07:51", "Asia/Shanghai")
	if err != nil {
		t.Errorf("DateUntilSecondToTimeStamp err: %v.", err)
	}
	t.Logf("DateUntilSecondToTimeStamp: %d.", stamp)

	defaultMillStamp, err := DateUntilMilliToTimeStamp("2018-07-11 15:07:51.100", "")
	if err != nil {
		t.Errorf("DateUntilMillToTimeStamp default err: %v.", err)
	}
	t.Logf("DateUntilMillToTimeStamp default: %d.", defaultMillStamp)

	millStamp, err := DateUntilMilliToTimeStamp("2018-07-11 15:07:51.100", "Asia/Shanghai")
	if err != nil {
		t.Errorf("DateUntilMillToTimeStamp err: %v.", err)
	}
	t.Logf("DateUntilMillToTimeStamp: %d.", millStamp)

	defaultNanoStamp, err := DateUntilNanoToTimeStamp("2018-07-11 15:07:51.100200000", "")
	if err != nil {
		t.Errorf("DateUntilNanoToTimeStamp default err: %v.", err)
	}
	t.Logf("DateUntilNanoToTimeStamp default: %d.", defaultNanoStamp)

	nanoStamp, err := DateUntilNanoToTimeStamp("2018-07-11 15:07:51.100200000", "Asia/Shanghai")
	if err != nil {
		t.Errorf("DateUntilMillToTimeStamp err: %v.", err)
	}
	t.Logf("DateUntilMillToTimeStamp: %d.", nanoStamp)

	t.Logf("Timestamp to date: %s.", TimeStampToDate(1531293019))
	t.Logf("Timestamp to date until second: %s.", TimeStampToDateUntilSecond(1531293019))
	t.Logf("Timestamp to date until mill: %s.", TimeStampToDateUntilMilli(1531293019))
	t.Logf("Timestamp to date until nano: %s.", TimeStampToDateUntilNano(1531293019))

}
