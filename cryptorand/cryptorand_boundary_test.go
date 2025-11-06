package cryptorand

import (
    "strings"
    "testing"
)

// TestRandomIntBoundaries 测试 RandomInt 的边界行为与范围正确性
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖 min==max、max<min 回退、以及多次随机结果均落于闭区间
func TestRandomIntBoundaries(t *testing.T) {
    // min==max 情况应恒等返回
    for i := 0; i < 10; i++ {
        v := RandomInt(5, 5)
        if v != 5 { t.Fatalf("min==max 时应返回自身: %d", v) }
    }
    // max<min 情况按实现回退到 min
    v := RandomInt(10, 3)
    if v != 10 { t.Fatalf("max<min 时应返回 min: %d", v) }
    // 多次随机均应落于闭区间
    min, max := -3, 7
    for i := 0; i < 200; i++ {
        r := RandomInt(min, max)
        if r < min || r > max { t.Fatalf("随机值越界: %d 不在 [%d,%d]", r, min, max) }
    }
}

// TestRandomStringAlphabet 测试 RandomString 的字符集与长度正确性
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：生成长字符串并断言所有字符均在 [a-zA-Z0-9] 集合内
func TestRandomStringAlphabet(t *testing.T) {
    s := RandomString(256)
    if len(s) != 256 { t.Fatalf("长度应为 256: %d", len(s)) }
    for i := 0; i < len(s); i++ {
        c := s[i]
        isDigit := c >= '0' && c <= '9'
        isLower := c >= 'a' && c <= 'z'
        isUpper := c >= 'A' && c <= 'Z'
        if !(isDigit || isLower || isUpper) {
            t.Fatalf("出现非字母数字字符: %q", c)
        }
    }
    // 非正长度应返回空字符串
    if RandomString(0) != "" { t.Fatalf("length<=0 应返回空字符串") }
}

// TestUUIDv4FormatAndVariantUniqueness 测试 UUIDv4 的格式、版本/变体位与近似唯一性
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：
// 1) 断言第三段首字符为 '4'（版本位）；
// 2) 第四段首字符为 8/9/a/b（变体位）；
// 3) 多次生成近似唯一（集合大小占多数）。
func TestUUIDv4FormatAndVariantUniqueness(t *testing.T) {
    seen := make(map[string]struct{}, 200)
    for i := 0; i < 200; i++ {
        u := UUIDv4()
        parts := strings.Split(u, "-")
        if len(parts) != 5 { t.Fatalf("UUID 分段应为5: %v", parts) }
        if len(parts[2]) == 0 || parts[2][0] != '4' {
            t.Fatalf("版本位应为4: %s", u)
        }
        if len(parts[3]) == 0 {
            t.Fatalf("变体段长度异常: %s", u)
        }
        v := parts[3][0]
        if !(v == '8' || v == '9' || v == 'a' || v == 'A' || v == 'b' || v == 'B') {
            t.Fatalf("变体位应为8/9/a/b: %s", u)
        }
        seen[u] = struct{}{}
    }
    if len(seen) < 190 { // 保守阈值，避免极小概率碰撞导致测试不稳定
        t.Fatalf("近似唯一性不足: 总数=200 唯一=%d", len(seen))
    }
}