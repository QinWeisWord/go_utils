package timeenv

import (
    "errors"
    "time"
)

// 本文件提供公历（阳历）与农历（阴历）之间的互相转换方法。
// 采用经典的农历编码表（1900-2099），通过位标识解析每年各月大小与闰月信息。
// 数据来源参考：lunarInfo 表的通用实现思路，见资料示例 [StackOverflow: lunarInfo 用法解释]、[JJonline 的 1900-2100 数据整理] 等。

// LunarDate 表示农历日期
// 字段 Year: 农历年
// 字段 Month: 农历月（1-12）
// 字段 Day: 农历日（1-29/30）
// 字段 IsLeap: 是否闰月（true 表示闰月）
type LunarDate struct {
    Year   int  // 农历年
    Month  int  // 农历月（1-12）
    Day    int  // 农历日（1-29/30）
    IsLeap bool // 是否闰月
}

// SolarToLunar 将公历时间转换为农历日期（本地时区）
// 参数 t: 指定的公历时间（会按本地时区的当天 00:00:00 进行计算）
// 返回值: 转换得到的农历日期与错误（年份不在支持范围、或内部计算异常时返回错误）
// 关键步骤：以 1900-01-31 作为农历 1900-正月初一基准，计算与目标日期的天数偏移，再逐年逐月扣减得到农历年月日与是否闰月
func SolarToLunar(t time.Time) (LunarDate, error) {
    loc := t.Location()
    base := time.Date(1900, time.January, 31, 0, 0, 0, 0, loc)
    dayStart := StartOfDay(t)
    // 限定支持范围：1900-01-31 至 2099-12-31（含）
    if dayStart.Before(base) || dayStart.After(time.Date(2099, time.December, 31, 23, 59, 59, int(time.Nanosecond-1), loc)) {
        return LunarDate{}, errors.New("仅支持 1900-01-31 到 2099-12-31 范围的农历转换")
    }

    // 关键步骤：计算与基准的天数偏移
    offset := int(dayStart.Sub(base) / (24 * time.Hour))

    // 逐年扣减，定位目标农历年
    y := 1900
    for y <= 2099 {
        diy := lunarYearDays(y)
        if offset < diy {
            break
        }
        offset -= diy
        y++
    }
    if y > 2099 {
        return LunarDate{}, errors.New("年份超出可支持范围")
    }

    // 逐月扣减，定位目标农历月与闰月状态
    leap := leapMonth(y)
    isLeap := false
    m := 1
    for m = 1; m <= 12; m++ {
        md := monthDays(y, m)
        if offset < md {
            break
        }
        offset -= md
        // 关键步骤：存在闰月时，普通月之后紧跟闰月，需单独判断并扣减闰月天数
        if leap != 0 && m == leap {
            ld := leapDays(y)
            if offset < ld {
                isLeap = true
                break
            }
            offset -= ld
        }
    }

    day := offset + 1
    return LunarDate{Year: y, Month: m, Day: day, IsLeap: isLeap}, nil
}

// LunarToSolar 将农历日期转换为公历时间（本地时区）
// 参数 ld: 农历日期（年范围需在 1900-2099）
// 参数 loc: 目标时区（用于构造返回的公历时间，通常为 time.Local）
// 返回值: 对应的公历时间（当天 00:00:00）与错误（不合法的日期或超出支持范围）
// 关键步骤：先累计基准到农历年的总天数，再累计至目标月（处理闰月），最后加上农历日偏移得到公历日期
func LunarToSolar(ld LunarDate, loc *time.Location) (time.Time, error) {
    if ld.Year < 1900 || ld.Year > 2099 {
        return time.Time{}, errors.New("仅支持 1900-2099 年的农历转换")
    }
    if ld.Month < 1 || ld.Month > 12 {
        return time.Time{}, errors.New("农历月需在 1-12 范围内")
    }
    base := time.Date(1900, time.January, 31, 0, 0, 0, 0, loc)

    // 验证日与闰月合法性
    lm := leapMonth(ld.Year)
    if ld.IsLeap {
        if lm == 0 || lm != ld.Month {
            return time.Time{}, errors.New("该年无此闰月或闰月月份不匹配")
        }
        if ld.Day < 1 || ld.Day > leapDays(ld.Year) {
            return time.Time{}, errors.New("闰月的农历日不合法")
        }
    } else {
        if ld.Day < 1 || ld.Day > monthDays(ld.Year, ld.Month) {
            return time.Time{}, errors.New("农历日不合法")
        }
    }

    // 关键步骤：累计从 1900-正月初一 到目标农历日期的总偏移天数
    offset := 0
    // 累计整年
    for y := 1900; y < ld.Year; y++ {
        offset += lunarYearDays(y)
    }
    // 累计当年月份（考虑闰月插入在对应月份之后）
    for m := 1; m < ld.Month; m++ {
        offset += monthDays(ld.Year, m)
        if lm != 0 && m == lm {
            offset += leapDays(ld.Year)
        }
    }
    // 若目标为闰月，则先加上该月的普通月天数以越过普通月到达闰月
    if ld.IsLeap {
        offset += monthDays(ld.Year, ld.Month)
    }
    // 农历日偏移（从初一开始）
    offset += ld.Day - 1

    // 基于偏移得到公历日期（当天 00:00:00）
    return base.Add(time.Duration(offset) * 24 * time.Hour), nil
}

// 以下为私有辅助方法与数据（置于公有方法之后）：

// lunarInfo 农历编码表（1900-2099）。
// 低4位为闰月月份（1-12，无闰为0）；
// 第5-16位为12个月大小月标识，1表示大月30天，0表示小月29天（从第16位对应一月，至第5位对应十二月）；
// 第17位及以上用于闰月大小标识（0x10000：闰大月30天，否则29天）。
var lunarInfo = []int{
    // 1900-1909
    0x04bd8, 0x04ae0, 0x0a570, 0x054d5, 0x0d260, 0x0d950, 0x16554, 0x056a0, 0x09ad0, 0x055d2,
    // 1910-1919
    0x04ae0, 0x0a5b6, 0x0a4d0, 0x0d250, 0x1d255, 0x0b540, 0x0d6a0, 0x0ada2, 0x095b0, 0x14977,
    // 1920-1929
    0x04970, 0x0a4b0, 0x0b4b5, 0x06a50, 0x06d40, 0x1ab54, 0x02b60, 0x09570, 0x052f2, 0x04970,
    // 1930-1939
    0x06566, 0x0d4a0, 0x0ea50, 0x16a95, 0x05ad0, 0x02b60, 0x186e3, 0x092e0, 0x1c8d7, 0x0c950,
    // 1940-1949
    0x0d4a0, 0x1d8a6, 0x0b550, 0x056a0, 0x1a5b4, 0x025d0, 0x092d0, 0x0d2b2, 0x0a950, 0x0b557,
    // 1950-1959
    0x06ca0, 0x0b550, 0x15355, 0x04da0, 0x0a5b0, 0x14573, 0x052b0, 0x0a9a8, 0x0e950, 0x06aa0,
    // 1960-1969
    0x0aea6, 0x0ab50, 0x04b60, 0x0aae4, 0x0a570, 0x05260, 0x0f263, 0x0d950, 0x05b57, 0x056a0,
    // 1970-1979
    0x096d0, 0x04dd5, 0x04ad0, 0x0a4d0, 0x0d4d4, 0x0d250, 0x0d558, 0x0b540, 0x0b6a0, 0x195a6,
    // 1980-1989
    0x095b0, 0x049b0, 0x0a974, 0x0a4b0, 0x0b27a, 0x06a50, 0x06d40, 0x0af46, 0x0ab60, 0x09570,
    // 1990-1999
    0x04af5, 0x04970, 0x064b0, 0x074a3, 0x0ea50, 0x06b58, 0x05ac0, 0x0ab60, 0x096d5, 0x092e0,
    // 2000-2009
    0x0c960, 0x0d954, 0x0d4a0, 0x0da50, 0x07552, 0x056a0, 0x0abb7, 0x025d0, 0x092d0, 0x0cab5,
    // 2010-2019
    0x0a950, 0x0b4a0, 0x0baa4, 0x0ad50, 0x055d9, 0x04ba0, 0x0a5b0, 0x15176, 0x052b0, 0x0a930,
    // 2020-2029
    0x07954, 0x06aa0, 0x0ad50, 0x05b52, 0x04b60, 0x0a6e6, 0x0a4e0, 0x0d260, 0x0ea65, 0x0d530,
    // 2030-2039
    0x05aa0, 0x076a3, 0x096d0, 0x04afb, 0x04ad0, 0x0a4d0, 0x1d0b6, 0x0d250, 0x0d520, 0x0dd45,
    // 2040-2049
    0x0b5a0, 0x056d0, 0x055b2, 0x049b0, 0x0a577, 0x0a4b0, 0x0aa50, 0x1b255, 0x06d20, 0x0ada0,
    // 2050-2059（扩展数据，便于后续拓展，当前转换仍限制到 2099）
    0x14b63, 0x09370, 0x049f8, 0x04970, 0x064b0, 0x168a6, 0x0ea50, 0x06b20, 0x1a6c4, 0x0aae0,
    // 2060-2069
    0x092e0, 0x0d2e3, 0x0c960, 0x0d557, 0x0d4a0, 0x0da50, 0x05d55, 0x056a0, 0x0a6d0, 0x055d4,
    // 2070-2079
    0x052d0, 0x0a9b8, 0x0a950, 0x0b4a0, 0x0b6a6, 0x0ad50, 0x055a0, 0x0aba4, 0x0a5b0, 0x052b0,
    // 2080-2089
    0x0b273, 0x06930, 0x07337, 0x06aa0, 0x0ad50, 0x14b55, 0x04b60, 0x0a570, 0x054e4, 0x0d160,
    // 2090-2099
    0x0e968, 0x0d520, 0x0daa0, 0x16aa6, 0x056d0, 0x04ae0, 0x0a9d4, 0x0a2d0, 0x0d150, 0x0f252,
}

// lunarYearDays 计算农历某年的总天数（含闰月）
// 参数 y: 年份
// 返回值: 该农历年的总天数
// 关键步骤：基础 12×29=348 天 + 每个大月额外加 1 天，再加上闰月天数
func lunarYearDays(y int) int {
    sum := 348
    // 关键步骤：从第16位到第5位依次判断12个月是否为大月
    for i := 0x8000; i > 0x8; i >>= 1 {
        if (lunarInfo[y-1900] & i) != 0 {
            sum += 1
        }
    }
    return sum + leapDays(y)
}

// leapDays 返回农历闰月的天数（无闰月返回 0）
// 参数 y: 年份
// 返回值: 闰月天数（30 或 29）或无闰时为 0
// 关键步骤：根据 0x10000 位判断闰月大小
func leapDays(y int) int {
    if leapMonth(y) != 0 {
        if (lunarInfo[y-1900] & 0x10000) != 0 {
            return 30
        }
        return 29
    }
    return 0
}

// leapMonth 返回农历闰月的月份（1-12，无闰返回 0）
// 参数 y: 年份
// 返回值: 闰月的月份编号或 0 表示无闰月
// 关键步骤：取编码表的低4位
func leapMonth(y int) int {
    return lunarInfo[y-1900] & 0xf
}

// monthDays 返回农历某年某月的天数（29 或 30）
// 参数 y: 年份
// 参数 m: 月份（1-12）
// 返回值: 该月天数（29/30）
// 关键步骤：判断 (0x10000 >> m) 位是否为 1，1 表示大月 30 天
func monthDays(y int, m int) int {
    if (lunarInfo[y-1900] & (0x10000 >> m)) == 0 {
        return 29
    }
    return 30
}