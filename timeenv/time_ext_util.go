package timeenv

import "time"

// 本文件提供日期与时间的扩展方法（周/月/年边界、差值、截断与解析）

// StartOfWeek 获取给定时间所在周的起始时间（本地时区）
// 参数 t: 指定时间
// 参数 firstWeekday: 周的起始星期（例如 time.Monday 或 time.Sunday）
// 返回值: 该周起始的 00:00:00 时间对象
// 关键步骤：根据给定起始星期计算差值，回退到本周起点并清零时分秒
func StartOfWeek(t time.Time, firstWeekday time.Weekday) time.Time {
    sd := StartOfDay(t)
    wd := sd.Weekday()
    // 关键步骤：计算距离周起点的天数差
    delta := (7 + int(wd) - int(firstWeekday)) % 7
    return sd.Add(-time.Duration(delta) * 24 * time.Hour)
}

// EndOfWeek 获取给定时间所在周的结束时间（本地时区）
// 参数 t: 指定时间
// 参数 firstWeekday: 周的起始星期（例如 time.Monday 或 time.Sunday）
// 返回值: 该周结束的 23:59:59.999999999 时间对象
// 关键步骤：在周起点基础上加7天减去1纳秒
func EndOfWeek(t time.Time, firstWeekday time.Weekday) time.Time {
    sow := StartOfWeek(t, firstWeekday)
    return sow.Add(7*24*time.Hour - time.Nanosecond)
}

// StartOfMonth 获取给定时间所在月的起始时间（本地时区）
// 参数 t: 指定时间
// 返回值: 当月第一天的 00:00:00 时间对象
// 关键步骤：使用 time.Date 重建时间为当月1号
func StartOfMonth(t time.Time) time.Time {
    y, m, _ := t.Date()
    loc := t.Location()
    return time.Date(y, m, 1, 0, 0, 0, 0, loc)
}

// EndOfMonth 获取给定时间所在月的结束时间（本地时区）
// 参数 t: 指定时间
// 返回值: 当月最后一天的 23:59:59.999999999 时间对象
// 关键步骤：下月1号减去1纳秒
func EndOfMonth(t time.Time) time.Time {
    som := StartOfMonth(t)
    next := AddMonths(som, 1)
    return next.Add(-time.Nanosecond)
}

// StartOfYear 获取给定时间所在年的起始时间（本地时区）
// 参数 t: 指定时间
// 返回值: 当年1月1日的 00:00:00 时间对象
// 关键步骤：使用 time.Date 构造年起始
func StartOfYear(t time.Time) time.Time {
    y := t.Year()
    loc := t.Location()
    return time.Date(y, time.January, 1, 0, 0, 0, 0, loc)
}

// EndOfYear 获取给定时间所在年的结束时间（本地时区）
// 参数 t: 指定时间
// 返回值: 当年12月31日的 23:59:59.999999999 时间对象
// 关键步骤：下一年1月1日减去1纳秒
func EndOfYear(t time.Time) time.Time {
    soy := StartOfYear(t)
    next := time.Date(soy.Year()+1, time.January, 1, 0, 0, 0, 0, soy.Location())
    return next.Add(-time.Nanosecond)
}

// StartOfQuarter 获取给定时间所在季度的起始时间（本地时区）
// 参数 t: 指定时间
// 返回值: 当季第一天的 00:00:00 时间对象
// 关键步骤：根据月份计算季度起始月份（1/4/7/10），构造该月1号的起始时间
func StartOfQuarter(t time.Time) time.Time {
    y := t.Year()
    m := t.Month()
    loc := t.Location()
    // 关键步骤：((m-1)/3)*3+1 → 1/4/7/10
    qm := time.Month(((int(m)-1)/3)*3 + 1)
    return time.Date(y, qm, 1, 0, 0, 0, 0, loc)
}

// EndOfQuarter 获取给定时间所在季度的结束时间（本地时区）
// 参数 t: 指定时间
// 返回值: 当季最后一天的 23:59:59.999999999 时间对象
// 关键步骤：季度起始 +3 个月减去 1 纳秒
func EndOfQuarter(t time.Time) time.Time {
    soq := StartOfQuarter(t)
    next := AddMonths(soq, 3)
    return next.Add(-time.Nanosecond)
}

// GetWeekRange 获取指定时间所在周的时间区间（起止）
// 参数 t: 指定时间
// 参数 firstWeekday: 周的起始星期（例如 time.Monday 或 time.Sunday）
// 返回值: 起始时间与结束时间（本地时区）
// 关键步骤：调用 StartOfWeek/EndOfWeek 获取范围
func GetWeekRange(t time.Time, firstWeekday time.Weekday) (time.Time, time.Time) {
    s := StartOfWeek(t, firstWeekday)
    e := EndOfWeek(t, firstWeekday)
    return s, e
}

// GetMonthRange 获取指定时间所在月的时间区间（起止）
// 参数 t: 指定时间
// 返回值: 起始时间与结束时间（本地时区）
// 关键步骤：调用 StartOfMonth/EndOfMonth 获取范围
func GetMonthRange(t time.Time) (time.Time, time.Time) {
    s := StartOfMonth(t)
    e := EndOfMonth(t)
    return s, e
}

// GetQuarterRange 获取指定时间所在季度的时间区间（起止）
// 参数 t: 指定时间
// 返回值: 起始时间与结束时间（本地时区）
// 关键步骤：调用 StartOfQuarter/EndOfQuarter 获取范围
func GetQuarterRange(t time.Time) (time.Time, time.Time) {
    s := StartOfQuarter(t)
    e := EndOfQuarter(t)
    return s, e
}

// GetYearRange 获取指定时间所在年的时间区间（起止）
// 参数 t: 指定时间
// 返回值: 起始时间与结束时间（本地时区）
// 关键步骤：调用 StartOfYear/EndOfYear 获取范围
func GetYearRange(t time.Time) (time.Time, time.Time) {
    s := StartOfYear(t)
    e := EndOfYear(t)
    return s, e
}

// AddMonths 在给定时间上增加指定的月数（处理月底溢出）
// 参数 t: 起始时间
// 参数 months: 增加的月数（可为负数）
// 返回值: 增加月数后的新时间（尽量保持原时分秒与时区）
// 关键步骤：计算目标年月并将“日”裁剪到目标月最大天数以避免溢出
func AddMonths(t time.Time, months int) time.Time {
    y, m, d := t.Date()
    loc := t.Location()
    // 关键步骤：计算新的年份与月份
    total := int(m) - 1 + months
    ny := y + total/12
    nm := time.Month(total%12 + 1)
    // 关键步骤：裁剪日期到目标月的最大天数
    md := daysInMonth(ny, nm)
    if d > md { d = md }
    return time.Date(ny, nm, d, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
}

// DiffDays 计算两个时间的相差天数（按各自本地起始时间计算，结果带符号）
// 参数 a: 第一个时间
// 参数 b: 第二个时间
// 返回值: b - a 的天数差（向下取整，可能受夏令时影响）
// 关键步骤：分别取各自的当天起始时间，求差并除以24小时
func DiffDays(a, b time.Time) int {
    sa := StartOfDay(a)
    sb := StartOfDay(b)
    dur := sb.Sub(sa)
    return int(dur / (24 * time.Hour))
}

// DiffHours 计算两个时间的相差小时数（结果带符号）
// 参数 a: 第一个时间
// 参数 b: 第二个时间
// 返回值: b - a 的小时数差（向下取整）
// 关键步骤：直接以时间差除以1小时
func DiffHours(a, b time.Time) int64 {
    dur := b.Sub(a)
    return int64(dur / time.Hour)
}

// IsSameDay 判断两个时间是否为同一天（各自本地时区）
// 参数 a: 第一个时间
// 参数 b: 第二个时间
// 返回值: 布尔值；true 表示同一天
// 关键步骤：比较各自的年月日
func IsSameDay(a, b time.Time) bool {
    ya, ma, da := a.Date()
    yb, mb, db := b.Date()
    return ya == yb && ma == mb && da == db
}

// IsWeekend 判断给定时间是否为周末（星期六或星期日）
// 参数 t: 指定时间
// 返回值: 布尔值；true 表示周末
// 关键步骤：判断 Weekday 是否为 time.Saturday 或 time.Sunday
func IsWeekend(t time.Time) bool {
    wd := t.Weekday()
    return wd == time.Saturday || wd == time.Sunday
}

// NextWeekday 获取下一个指定星期的日期（不含当天）
// 参数 t: 起始时间
// 参数 weekday: 目标星期（例如 time.Friday）
// 返回值: 下一个该星期的当天 00:00:00 时间对象（本地时区）
// 关键步骤：以当天起始为基准计算到目标星期的偏移
func NextWeekday(t time.Time, weekday time.Weekday) time.Time {
    sd := StartOfDay(t)
    wd := sd.Weekday()
    // 关键步骤：与当天同星期时返回下一周的该星期
    delta := (7 + int(weekday) - int(wd)) % 7
    if delta == 0 { delta = 7 }
    return sd.Add(time.Duration(delta) * 24 * time.Hour)
}

// TruncateToHour 将时间截断到小时（清零分钟/秒/纳秒）
// 参数 t: 指定时间
// 返回值: 截断后的时间对象
// 关键步骤：使用 time.Date 重建并清零分秒
func TruncateToHour(t time.Time) time.Time {
    y, m, d := t.Date()
    loc := t.Location()
    return time.Date(y, m, d, t.Hour(), 0, 0, 0, loc)
}

// FormatRFC3339 将时间格式化为 RFC3339 字符串（ISO8601）
// 参数 t: 指定时间
// 返回值: RFC3339 字符串（包含时区偏移）
// 关键步骤：调用 t.Format(time.RFC3339)
func FormatRFC3339(t time.Time) string {
    return t.Format(time.RFC3339)
}

// ParseRFC3339 解析 RFC3339 格式时间字符串（ISO8601）
// 参数 s: RFC3339 字符串
// 返回值: 解析得到的时间对象与错误（若成功则error为nil）
// 关键步骤：调用 time.Parse(time.RFC3339, s)
func ParseRFC3339(s string) (time.Time, error) {
    return time.Parse(time.RFC3339, s)
}

// ToLocal 将时间转换到指定时区（仅改变显示时区，不改变绝对时间）
// 参数 t: 原始时间
// 参数 loc: 目标时区（例如 time.Local 或 time.FixedZone("UTC+8", 8*3600)）
// 返回值: 转换到指定时区的时间对象
// 关键步骤：调用 t.In(loc)
func ToLocal(t time.Time, loc *time.Location) time.Time {
    return t.In(loc)
}

// daysInMonth 获取指定年/月的天数（私有辅助函数，置于公有方法之后）
// 参数 year: 年份
// 参数 month: 月份
// 返回值: 该月天数（28~31）
// 关键步骤：通过下月第0天技巧获得当月最后一天
func daysInMonth(year int, month time.Month) int {
    t := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local)
    return t.Day()
}