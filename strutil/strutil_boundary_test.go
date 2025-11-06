package strutil

import (
    "math/rand"
    "testing"
    "unicode/utf8"
)

// TestSubstringBoundaries 测试 Substring 在负起点/越界/零长度/负长度的行为
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：验证UTF-8安全切割与各边界分支
func TestSubstringBoundaries(t *testing.T) {
    s := "你😀好abc"
    // 负起点→视为0
    if sub := Substring(s, -2, 2); sub != "你😀" { t.Fatalf("负起点失败: %q", sub) }
    // 起点越界→空串
    if sub := Substring(s, 100, 1); sub != "" { t.Fatalf("越界应为空: %q", sub) }
    // 长度为0→空串
    if sub := Substring(s, 1, 0); sub != "" { t.Fatalf("长度为0应为空: %q", sub) }
    // 负长度→截到字符串末尾
    if sub := Substring(s, 2, -1); sub != "好abc" { t.Fatalf("负长度应到末尾: %q", sub) }
}

// TestReverseTwiceRandom 测试 Reverse 在随机字符串上的双反转等价性
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：随机生成包含中文与Emoji的字符串，验证 Reverse(Reverse(s)) == s
func TestReverseTwiceRandom(t *testing.T) {
    rand.Seed(123456)
    pool := []rune("abcXYZ你我他😀🚀🧡")
    for i := 0; i < 100; i++ {
        n := rand.Intn(40)
        rs := make([]rune, n)
        for j := 0; j < n; j++ { rs[j] = pool[rand.Intn(len(pool))] }
        s := string(rs)
        if Reverse(Reverse(s)) != s { t.Fatalf("双反转不等价: %q", s) }
    }
}

// TestPadLeftRightUnicode 测试 PadLeft/PadRight 在Unicode与精确长度上的行为
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：单字符中文补齐应达到目标长度；已达长度不应变化
func TestPadLeftRightUnicode(t *testing.T) {
    // 左补齐
    s1 := PadLeft("你a", "字", 5)
    if utf8.RuneCountInString(s1) != 5 { t.Fatalf("PadLeft 目标长度应为5: %d", utf8.RuneCountInString(s1)) }
    // 右补齐
    s2 := PadRight("你a", "字", 5)
    if utf8.RuneCountInString(s2) != 5 { t.Fatalf("PadRight 目标长度应为5: %d", utf8.RuneCountInString(s2)) }
    // 已达或超过长度不变化
    s3 := PadLeft("你好世界", "字", 4)
    if s3 != "你好世界" { t.Fatalf("达到长度时不应变化: %q", s3) }
}

// TestContainsSubstrEmpty 测试空子串的包含行为
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：空子串应视为包含（与strings.Contains一致）
func TestContainsSubstrEmpty(t *testing.T) {
    if !ContainsSubstr("abc", "") { t.Fatalf("空子串应返回包含") }
}