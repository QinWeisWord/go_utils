package timeenv

import (
    "testing"
    "time"
)

// TestDayBoundariesAndAdd 测试起止时间与加减天数
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：验证 StartOfDay/EndOfDay/AddDays 的基本行为
func TestDayBoundariesAndAdd(t *testing.T) {
    loc := time.FixedZone("UTC+8", 8*3600)
    base := time.Date(2024, 3, 14, 10, 30, 5, 123, loc)
    sod := StartOfDay(base)
    eod := EndOfDay(base)
    if sod.Hour() != 0 || sod.Minute() != 0 || sod.Second() != 0 || sod.Nanosecond() != 0 {
        t.Fatalf("当天起始应为 00:00:00")
    }
    if !eod.After(sod) || eod.Sub(sod) != 24*time.Hour-time.Nanosecond {
        t.Fatalf("当天结束与起始间隔应为 24h-1ns")
    }
    next := AddDays(base, 2)
    if DiffDays(base, next) != 2 {
        t.Fatalf("DiffDays 计算错误")
    }
}

// TestWeekMonthYearBoundaries 测试周/月/年起止与范围
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：验证 Start/EndOfWeek/Month/Year 与范围函数返回一致
func TestWeekMonthYearBoundaries(t *testing.T) {
    t0 := time.Date(2024, 6, 6, 12, 0, 0, 0, time.Local) // 周四
    sow := StartOfWeek(t0, time.Monday)
    eow := EndOfWeek(t0, time.Monday)
    s2, e2 := GetWeekRange(t0, time.Monday)
    if sow != s2 || eow != e2 { t.Fatalf("周范围不一致") }
    if sow.Weekday() != time.Monday || eow.Sub(sow) != 7*24*time.Hour-time.Nanosecond {
        t.Fatalf("周起止计算错误")
    }
    som := StartOfMonth(t0)
    eom := EndOfMonth(t0)
    sm, em := GetMonthRange(t0)
    if som != sm || eom != em { t.Fatalf("月范围不一致") }
    soy := StartOfYear(t0)
    eoy := EndOfYear(t0)
    sy, ey := GetYearRange(t0)
    if soy != sy || eoy != ey { t.Fatalf("年范围不一致") }
}

// TestAddMonthsAndDiff 测试按月增加与时间差
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖月底溢出裁剪与小时差
func TestAddMonthsAndDiff(t *testing.T) {
    loc := time.Local
    // 闰年 2024-01-31 +1 月 → 2024-02-29
    a := time.Date(2024, 1, 31, 8, 0, 0, 0, loc)
    b := AddMonths(a, 1)
    if b.Day() != 29 || b.Month() != time.February { t.Fatalf("闰年二月裁剪失败: %v", b) }
    // 非闰年 2023-01-31 +1 月 → 2023-02-28
    c := time.Date(2023, 1, 31, 8, 0, 0, 0, loc)
    d := AddMonths(c, 1)
    if d.Day() != 28 || d.Month() != time.February { t.Fatalf("平年二月裁剪失败: %v", d) }
    // 小时差
    h := DiffHours(a, a.Add(5*time.Hour+30*time.Minute))
    if h != 5 { t.Fatalf("小时差应向下取整为5: %d", h) }
}

// TestSameWeekendNextWeekdayAndTruncate 测试同一天/周末/下个星期与截断
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖 IsSameDay/IsWeekend/NextWeekday/TruncateToHour
func TestSameWeekendNextWeekdayAndTruncate(t *testing.T) {
    a := time.Date(2024, 7, 20, 9, 15, 0, 0, time.Local) // 周六
    b := time.Date(2024, 7, 20, 22, 0, 0, 0, time.Local)
    if !IsSameDay(a, b) { t.Fatalf("同一天判断失败") }
    if !IsWeekend(a) { t.Fatalf("周末判断失败") }
    nw := NextWeekday(a, time.Friday)
    if nw.Weekday() != time.Friday || !nw.After(StartOfDay(a)) { t.Fatalf("下个星期计算错误: %v", nw) }
    th := TruncateToHour(a)
    if th.Minute() != 0 || th.Second() != 0 || th.Nanosecond() != 0 { t.Fatalf("截断到小时失败") }
}

// TestFormatParseLocal 测试格式化/解析与时区转换
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖 FormatRFC3339/ParseRFC3339/ToLocal 及自定义布局
func TestFormatParseLocal(t *testing.T) {
    t0 := time.Date(2024, 8, 1, 12, 34, 56, 0, time.FixedZone("UTC+0", 0))
    s := FormatRFC3339(t0)
    p, err := ParseRFC3339(s)
    if err != nil || !p.Equal(t0) { t.Fatalf("RFC3339 格式解析失败: %v %v", p, err) }
    loc := time.FixedZone("UTC+8", 8*3600)
    tl := ToLocal(t0, loc)
    if tl.Location() != loc || tl.Sub(t0) != 0 { t.Fatalf("ToLocal 应不改变绝对时间") }
    if FormatTime(t0, "2006-01-02 15:04:05") != "2024-08-01 12:34:56" { t.Fatalf("FormatTime 失败") }
    if _, err := ParseTime("2006-01-02", "2024-08-01"); err != nil { t.Fatalf("ParseTime 失败: %v", err) }
}

// TestLunarConvertRoundtrip 测试公历与农历互转的往返一致性
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：选取多个春节日期，断言转换为农历正月初一并能还原
func TestLunarConvertRoundtrip(t *testing.T) {
    loc := time.Local
    cases := []time.Time{
        time.Date(2023, time.January, 22, 0, 0, 0, 0, loc), // 2023 春节
        time.Date(2024, time.February, 10, 0, 0, 0, 0, loc), // 2024 春节
        time.Date(2021, time.February, 12, 0, 0, 0, 0, loc), // 2021 春节
    }
    for _, c := range cases {
        ld, err := SolarToLunar(c)
        if err != nil { t.Fatalf("SolarToLunar 失败: %v", err) }
        if ld.Month != 1 || ld.Day != 1 { t.Fatalf("应为农历正月初一: %+v", ld) }
        s, err := LunarToSolar(ld, loc)
        if err != nil { t.Fatalf("LunarToSolar 失败: %v", err) }
        if !IsSameDay(c, s) { t.Fatalf("往返日期不一致: solar=%v back=%v", c, s) }
    }
}