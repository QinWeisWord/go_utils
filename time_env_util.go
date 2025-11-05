package utils

import (
    "os"
    "time"
)

// NowUnix 获取当前时间的Unix秒时间戳
// 参数: 无
// 返回值: 当前时间的Unix秒（int64）
// 关键步骤：调用time.Now().Unix()
func NowUnix() int64 {
    return time.Now().Unix()
}

// NowUnixMilli 获取当前时间的Unix毫秒时间戳
// 参数: 无
// 返回值: 当前时间的Unix毫秒（int64）
// 关键步骤：调用time.Now().UnixMilli()
func NowUnixMilli() int64 {
    return time.Now().UnixMilli()
}

// FormatNow 按布局格式化当前时间
// 参数 layout: 时间布局字符串（例如 2006-01-02 15:04:05）
// 返回值: 格式化后的时间字符串
// 关键步骤：调用time.Now().Format(layout)
func FormatNow(layout string) string {
    return time.Now().Format(layout)
}

// FormatTime 按布局格式化指定时间
// 参数 t: 需要格式化的时间对象
// 参数 layout: 时间布局字符串
// 返回值: 格式化后的时间字符串
// 关键步骤：调用t.Format(layout)
func FormatTime(t time.Time, layout string) string {
    return t.Format(layout)
}

// ParseTime 解析时间字符串为时间对象
// 参数 layout: 时间布局字符串
// 参数 value: 需要解析的时间字符串
// 返回值: 解析得到的时间对象与错误信息（若成功则error为nil）
// 关键步骤：调用time.Parse(layout, value)
func ParseTime(layout, value string) (time.Time, error) {
    return time.Parse(layout, value)
}

// AddDays 在给定时间上增加指定天数
// 参数 t: 起始时间
// 参数 days: 增加的天数（可为负数表示减少）
// 返回值: 增加天数后的新时间
// 关键步骤：用time.Duration计算天数对应的小时
func AddDays(t time.Time, days int) time.Time {
    return t.Add(time.Duration(days) * 24 * time.Hour)
}

// StartOfDay 获取给定时间的当天起始时间（本地时区）
// 参数 t: 指定时间
// 返回值: 当天00:00:00的时间对象
// 关键步骤：取年、月、日再构造新的时间对象
func StartOfDay(t time.Time) time.Time {
    y, m, d := t.Date()
    loc := t.Location()
    return time.Date(y, m, d, 0, 0, 0, 0, loc)
}

// EndOfDay 获取给定时间的当天结束时间（本地时区）
// 参数 t: 指定时间
// 返回值: 当天23:59:59.999999999的时间对象
// 关键步骤：用下一天的起始时间减去1纳秒
func EndOfDay(t time.Time) time.Time {
    nextDay := AddDays(StartOfDay(t), 1)
    return nextDay.Add(-time.Nanosecond)
}

// GetEnv 获取环境变量的值
// 参数 key: 环境变量名称
// 返回值: 环境变量值；若不存在则返回空字符串
// 关键步骤：调用os.Getenv
func GetEnv(key string) string {
    return os.Getenv(key)
}

// GetEnvDefault 获取环境变量值，不存在时返回默认值
// 参数 key: 环境变量名称
// 参数 def: 默认值，当环境变量不存在或为空时返回
// 返回值: 环境变量值或默认值
// 关键步骤：获取后判断是否为空
func GetEnvDefault(key, def string) string {
    val := os.Getenv(key)
    if len(val) == 0 {
        return def
    }
    return val
}