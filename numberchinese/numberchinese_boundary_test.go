package numberchinese

import "testing"

// TestLowerFloatRoundingBoundaries 测试小写浮点的四舍五入边界
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖 .005 向上进位与进位到整数的场景
func TestLowerFloatRoundingBoundaries(t *testing.T) {
    if s := ToChineseLowerFloat(1.005, 2); s != "一" {
        t.Fatalf("1.005 保留2位按当前四舍五入应为 一: %q", s)
    }
    if s := ToChineseLowerFloat(9.999, 2); s != "十" {
        t.Fatalf("9.999 保留2位应进位为 十: %q", s)
    }
    if s := ToChineseLowerFloat(-0.5, 0); s != "负一" {
        t.Fatalf("-0.5 保留0位应为 负一: %q", s)
    }
    if s := ToChineseLowerFloat(-1.234, 2); s != "负一点二三" {
        t.Fatalf("-1.234 保留2位应为 负一点二三: %q", s)
    }
}

// TestUpperFloatRoundingBoundaries 测试大写浮点的四舍五入边界
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖 .005 向上进位与进位到整数的场景（大写数字）
func TestUpperFloatRoundingBoundaries(t *testing.T) {
    if s := ToChineseUpperFloat(1.005, 2); s != "壹" {
        t.Fatalf("1.005 保留2位按当前四舍五入应为 壹: %q", s)
    }
    if s := ToChineseUpperFloat(9.995, 2); s != "玖点玖玖" {
        t.Fatalf("9.995 保留2位按当前四舍五入为 玖点玖玖: %q", s)
    }
    if s := ToChineseUpperFloat(-0.5, 0); s != "负壹" {
        t.Fatalf("-0.5 保留0位应为 负壹: %q", s)
    }
    if s := ToChineseUpperFloat(-2.345, 2); s != "负贰点叁伍" {
        t.Fatalf("-2.345 保留2位应为 负贰点叁伍: %q", s)
    }
}

// TestLowerIntBigUnitsZeros 测试小写整数在大单位与零插入上的边界
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖万/亿等组边界以及跨组零插入规则
func TestLowerIntBigUnitsZeros(t *testing.T) {
    if s := ToChineseLowerInt(10000); s != "一万" {
        t.Fatalf("10000 应为 一万: %q", s)
    }
    if s := ToChineseLowerInt(10001); s != "一万零一" {
        t.Fatalf("10001 应为 一万零一: %q", s)
    }
    if s := ToChineseLowerInt(10010); s != "一万零十" {
        t.Fatalf("10010 应为 一万零十: %q", s)
    }
    if s := ToChineseLowerInt(1000000); s != "一百万" {
        t.Fatalf("1000000 应为 一百万: %q", s)
    }
    if s := ToChineseLowerInt(1001000); s != "一百万一千" {
        t.Fatalf("1001000 应为 一百万一千: %q", s)
    }
    if s := ToChineseLowerInt(100010001); s != "一亿零一万零一" {
        t.Fatalf("100010001 应为 一亿零一万零一: %q", s)
    }
}

// TestUpperIntBigUnitsZeros 测试大写整数在大单位与零插入上的边界
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：与小写规则一致，但10~19不省略“壹”
func TestUpperIntBigUnitsZeros(t *testing.T) {
    if s := ToChineseUpperInt(10); s != "壹十" {
        t.Fatalf("10（大写）应为 壹十: %q", s)
    }
    if s := ToChineseUpperInt(10010); s != "壹万零壹十" {
        t.Fatalf("10010（大写）应为 壹万零壹十: %q", s)
    }
    if s := ToChineseUpperInt(100010001); s != "壹亿零壹万零壹" {
        t.Fatalf("100010001（大写）应为 壹亿零壹万零壹: %q", s)
    }
}