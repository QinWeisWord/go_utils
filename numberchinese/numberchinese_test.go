package numberchinese

import "testing"

// TestLowerInt 测试整数中文小写转换
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖正数、负数、十~十九的省略“一”规则
func TestLowerInt(t *testing.T) {
    if s := ToChineseLowerInt(12345); s != "一万二千三百四十五" {
        t.Fatalf("12345 小写应为 一万二千三百四十五: %s", s)
    }
    if s := ToChineseLowerInt(-10); s != "负十" {
        t.Fatalf("-10 小写应为 负十: %s", s)
    }
    if s := ToChineseLowerInt(11); s != "十零一" {
        t.Fatalf("11 小写应为 十零一: %s", s)
    }
    if s := ToChineseLowerInt(20); s != "二十" {
        t.Fatalf("20 小写应为 二十: %s", s)
    }
}

// TestUpperInt 测试整数中文大写转换
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖十位不省略“壹”、多组单位与负数
func TestUpperInt(t *testing.T) {
    if s := ToChineseUpperInt(10010); s != "壹万零壹十" {
        t.Fatalf("10010 大写应为 壹万零壹十: %s", s)
    }
    if s := ToChineseUpperInt(-10); s != "负壹十" {
        t.Fatalf("-10 大写应为 负壹十: %s", s)
    }
}

// TestLowerUpperFloat 测试浮点中文大小写转换
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖四舍五入、小数部分逐位读出与负数
func TestLowerUpperFloat(t *testing.T) {
    if s := ToChineseLowerFloat(10.25, 2); s != "十点二五" {
        t.Fatalf("10.25 小写应为 十点二五: %s", s)
    }
    if s := ToChineseUpperFloat(1001.05, 2); s != "壹千零壹点零伍" {
        t.Fatalf("1001.05 大写应为 壹千零壱点零伍: %s", s)
    }
    if s := ToChineseLowerFloat(-3.0, 0); s != "负三" {
        t.Fatalf("-3 小写应为 负三: %s", s)
    }
}

// TestLowerUpperNumber 测试通用数字到中文的大小写入口
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖int/uint/float三类
func TestLowerUpperNumber(t *testing.T) {
    if s, err := ToChineseLowerNumber(uint32(12), 0); err != nil || s != "十零二" {
        t.Fatalf("Lower uint32 应为 十零二: %s err=%v", s, err)
    }
    if s, err := ToChineseUpperNumber(12.34, 2); err != nil || s != "壹十零贰点叁肆" {
        t.Fatalf("Upper float 应为 壹十零贰点叁肆: %s err=%v", s, err)
    }
}

// TestRMBUpper 测试人民币金额中文大写
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖角分与“整”的输出
func TestRMBUpper(t *testing.T) {
    if s := ToChineseRMBUpper(123.45); s != "壹百零贰十叁元肆角伍分" {
        t.Fatalf("123.45 RMB 大写应为 壹百零贰十叁元肆角伍分: %s", s)
    }
    if s := ToChineseRMBUpper(300); s != "叁百元整" {
        t.Fatalf("300 RMB 大写应为 叁百元整: %s", s)
    }
}