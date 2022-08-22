package utils

import (
	"time"
)

// DATE 日期格式常量
const (
	DATE              = "2006-01-02"                    // 日期
	DATA_UNTIL_SECOND = "2006-01-02 15:04:05"           // 日期具体到秒
	DATE_UNTIL_MILLI  = "2006-01-02 15:04:05.000"       // 日期具体到毫秒
	DATE_UNTIL_NANO   = "2006-01-02 15:04:05.000000000" // 日期具体到纳秒
	UTC_DATE          = "2006-01-02T15:04:05.000Z"      // 日期（附带UTC 时区）
	BEIJING_DATE      = "2006-01-02T15:04:05+08:00"     // 日期（北京时区）
)

// GetCurrentTime 获取当前时间
func GetCurrentTime() time.Time {
	return time.Now()
}

// GetCurrentTimestamp 获取当前Unix 时间戳
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// GetCurrentMilliTimestamp 获取当前Unix 时间戳（精确到毫秒）
func GetCurrentMilliTimestamp() int64 {
	return time.Now().UnixMilli()
}

// GetCurrentNanoTimestamp 获取当前Unix 时间戳（精确到纳秒）
func GetCurrentNanoTimestamp() int64 {
	return time.Now().UnixNano()
}

// GetCurrentNano 获取当前纳秒数，单位：纳秒
func GetCurrentNano() int {
	return time.Now().Nanosecond()
}

// GetCurrentYear 获取当前年
func GetCurrentYear() int {
	return time.Now().Year()
}

// GetCurrentMonth 获取当前月份（英文，示例：January）
func GetCurrentMonth() string {
	return time.Now().Month().String()
}

// GetCurrentWeekday 获取当前星期数
func GetCurrentWeekday() int {
	return int(time.Now().Weekday())
}

// GetCurrentDay 获取当前天数
func GetCurrentDay() int {
	return time.Now().Day()
}

// GetCurrentYearMonthDay 分别获取年、月、日
func GetCurrentYearMonthDay() (int, int, int) {
	year, month, day := time.Now().Date()
	return year, int(month), day
}

// GetCurrentHour 获取当前小时数
func GetCurrentHour() int {
	return time.Now().Hour()
}

// GetCurrentMinute 获取当前分钟数
func GetCurrentMinute() int {
	return time.Now().Minute()
}

// GetCurrentSecond 获取当前秒数
func GetCurrentSecond() int {
	return time.Now().Second()
}

// GetCurrentDayOfYear 获取当前是当年中的第几天
func GetCurrentDayOfYear() int {
	return time.Now().YearDay()
}

// GetCurrentDate 获取当前日期（年月日）
func GetCurrentDate() string {
	return time.Unix(time.Now().Unix(), 0).Format(DATE)
}

// GetCurrentDateUntilSecond 获取当前日期（精确到秒）
func GetCurrentDateUntilSecond() string {
	return time.Unix(time.Now().Unix(), 0).Format(DATA_UNTIL_SECOND)
}

// GetCurrentDateUntilMilli 获取当前日期（精确到毫秒）
func GetCurrentDateUntilMilli() string {
	return time.Unix(time.Now().Unix(), 0).Format(DATE_UNTIL_MILLI)
}

// GetCurrentDateUntilNano 获取当前日期（精确到纳秒）
func GetCurrentDateUntilNano() string {
	return time.Unix(time.Now().Unix(), 0).Format(DATE_UNTIL_NANO)
}

// DateUntilSecondToTimeStamp 日期(精确到秒)转时间戳
/**
timeZone 时区
	本地：Local
	上海：Asia/Shanghai
*/
func DateUntilSecondToTimeStamp(value, timeZone string) (int64, error) {

	if timeZone == "" {
		timeZone = "Local"
	}
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return 0, err
	}
	tt, err := time.ParseInLocation(DATA_UNTIL_SECOND, value, loc)
	if err != nil {
		return 0, err
	}
	return tt.Unix(), nil
}

// DateUntilMilliToTimeStamp 日期(精确到毫秒)转时间戳
/**
timeZone 时区
	本地：Local
	上海：Asia/Shanghai
*/
func DateUntilMilliToTimeStamp(value, timeZone string) (int64, error) {
	if timeZone == "" {
		timeZone = "Local"
	}
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return 0, err
	}
	tt, err := time.ParseInLocation(DATE_UNTIL_MILLI, value, loc)
	if err != nil {
		return 0, err
	}
	return tt.UnixMilli(), nil
}

// DateUntilNanoToTimeStamp 日期(精确到毫秒)转时间戳
/**
timeZone 时区
	本地：Local
	上海：Asia/Shanghai
*/
func DateUntilNanoToTimeStamp(value, timeZone string) (int64, error) {
	if timeZone == "" {
		timeZone = "Local"
	}
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return 0, err
	}
	tt, err := time.ParseInLocation(DATE_UNTIL_NANO, value, loc)
	if err != nil {
		return 0, err
	}
	return tt.UnixNano(), nil
}

// TimeStampToDate 时间戳转日期
func TimeStampToDate(timeStamp int64) string {
	return time.Unix(timeStamp, 0).Format(DATE)
}

// TimeStampToDateUntilSecond 时间戳转日期(精确到秒)
func TimeStampToDateUntilSecond(timeStamp int64) string {
	return time.Unix(timeStamp, 0).Format(DATA_UNTIL_SECOND)
}

// TimeStampToDateUntilMilli 时间戳转日期(精确到毫秒)
func TimeStampToDateUntilMilli(timeStamp int64) string {
	return time.Unix(timeStamp, 0).Format(DATE_UNTIL_MILLI)
}

// TimeStampToDateUntilNano 时间戳转日期(精确到纳秒)
func TimeStampToDateUntilNano(timeStamp int64) string {
	return time.Unix(timeStamp, 0).Format(DATE_UNTIL_NANO)
}
