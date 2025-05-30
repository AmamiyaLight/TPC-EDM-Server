package parse_utils

import (
	"strconv"
	"time"
)

func ParseUint(s string) uint {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}
	return uint(val)
}
func ParseIntUtil(s string) int {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return int(val)
}

func ParseFloat64(s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val
}

func ParseTimeUtil(s string) time.Time {
	val, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}
	}
	return val
}

func StrConvUInt(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}
func StrConvInt(i int) string {
	return strconv.Itoa(i)
}

func StrConvFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
func StrConvTime(t time.Time) string {
	return t.Format("2006-01-02")
}
