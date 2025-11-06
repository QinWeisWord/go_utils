package timeenv

import (
    "math/rand"
    "testing"
    "time"
)

// TestAddMonthsNegativeEdges 测试 AddMonths 负月数在月底的裁剪
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖闰年与平年2月，月底向前减1月的裁剪行为
func TestAddMonthsNegativeEdges(t *testing.T) {
    // 闰年：2024-03-31 减1月 → 2024-02-29
    dt1 := time.Date(2024, time.March, 31, 12, 30, 0, 0, time.Local)
    got1 := AddMonths(dt1, -1)
    if got1.Year() != 2024 || got1.Month() != time.February || got1.Day() != 29 {
        t.Fatalf("闰年裁剪失败: got=%v", got1)
    }
    if got1.Hour() != 12 || got1.Minute() != 30 { t.Fatalf("应保留原时分: %v", got1) }

    // 平年：2023-03-31 减1月 → 2023-02-28
    dt2 := time.Date(2023, time.March, 31, 8, 15, 0, 0, time.Local)
    got2 := AddMonths(dt2, -1)
    if got2.Year() != 2023 || got2.Month() != time.February || got2.Day() != 28 {
        t.Fatalf("平年裁剪失败: got=%v", got2)
    }
    if got2.Hour() != 8 || got2.Minute() != 15 { t.Fatalf("应保留原时分: %v", got2) }

    // 跨年负月：2023-01-31 减1月 → 2022-12-31
    dt3 := time.Date(2023, time.January, 31, 23, 59, 59, 1, time.Local)
    got3 := AddMonths(dt3, -1)
    if got3.Year() != 2022 || got3.Month() != time.December || got3.Day() != 31 {
        t.Fatalf("跨年裁剪失败: got=%v", got3)
    }
}

// TestNextWeekdayRandom 测试 NextWeekday 的随机起点与各星期目标
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：随机生成日期，验证返回为目标星期且严格晚于起始当天
func TestNextWeekdayRandom(t *testing.T) {
    // 关键步骤：使用局部随机生成器以获得确定性序列（避免全局 Seed 弃用）
    rng := rand.New(rand.NewSource(20231111))
    for i := 0; i < 100; i++ {
        // 随机日期（近几年区间）
        y := 2020 + rng.Intn(6)
        m := time.Month(1 + rng.Intn(12))
        d := 1 + rng.Intn(28) // 避免无效日期
        base := time.Date(y, m, d, rng.Intn(24), rng.Intn(60), 0, 0, time.Local)
        target := time.Weekday(rng.Intn(7))
        next := NextWeekday(base, target)
        if next.Weekday() != target {
            t.Fatalf("返回星期不匹配: got=%v want=%v", next.Weekday(), target)
        }
        // 应严格晚于当天起始（不含当天）
        if !next.After(StartOfDay(base)) {
            t.Fatalf("返回时间不晚于当天起始: base=%v next=%v", base, next)
        }
    }
}

// TestNextWeekdaySameDayEdge 测试 NextWeekday 在同一星期时返回下一周
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：当天是目标星期，应返回7天后的该星期起始
func TestNextWeekdaySameDayEdge(t *testing.T) {
    base := time.Date(2025, time.November, 3, 9, 0, 0, 0, time.Local) // 需保证为周一示例
    // 若日期非周一，微调为周一
    for base.Weekday() != time.Monday { base = AddDays(base, 1) }
    next := NextWeekday(base, time.Monday)
    expect := StartOfDay(AddDays(base, 7))
    if !next.Equal(expect) { t.Fatalf("同日应返回下一周起始: next=%v expect=%v", next, expect) }
}