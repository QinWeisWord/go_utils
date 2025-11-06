package strutil

import (
    "reflect"
    "testing"
)

// TestBasic 测试基础字符串判断与转换
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖 IsEmpty/Trim/ToUpper/ToLower/ContainsSubstr/ReplaceAll
func TestBasic(t *testing.T) {
    if !IsEmpty("  \t\n") { t.Fatalf("空白字符串应视为空") }
    if IsEmpty("a") { t.Fatalf("非空字符串不应视为空") }
    if Trim("  a \t") != "a" { t.Fatalf("Trim 失败") }
    if ToUpper("Abc") != "ABC" { t.Fatalf("ToUpper 失败") }
    if ToLower("AbC") != "abc" { t.Fatalf("ToLower 失败") }
    if !ContainsSubstr("hello", "ll") { t.Fatalf("包含判断失败") }
    if ReplaceAll("a-b-a", "-", "+") != "a+b+a" { t.Fatalf("ReplaceAll 失败") }
}

// TestSplitJoinSubstringReversePad 测试拆分/连接/子串/反转/补齐
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：验证 UTF-8 安全子串与补齐长度
func TestSplitJoinSubstringReversePad(t *testing.T) {
    // Split/Join
    parts := Split("a,b,c", ",")
    if !reflect.DeepEqual(parts, []string{"a", "b", "c"}) { t.Fatalf("Split 失败: %v", parts) }
    if Join(parts, ":") != "a:b:c" { t.Fatalf("Join 失败") }

    // UTF-8 子串与反转
    s := "你😀好"
    sub := Substring(s, 1, 2) // 期待取到 😀好
    if sub != "😀好" { t.Fatalf("Substring 失败: %q", sub) }
    rev := Reverse("中国abc")
    if rev != "cba国中" { t.Fatalf("Reverse 失败: %q", rev) }

    // PadLeft/PadRight（按字符长度）
    if PadLeft("ab", "_", 5) != "___ab" { t.Fatalf("PadLeft 失败") }
    if PadRight("ab", "_", 5) != "ab___" { t.Fatalf("PadRight 失败") }
}