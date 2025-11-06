package collections

import (
    "math/rand"
    "testing"
)

// TestUniqueOrderStability 测试 Unique 的顺序稳定与随机重复场景
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：构造小域随机序列以产生大量重复，验证保留首次出现顺序
func TestUniqueOrderStability(t *testing.T) {
    // 关键步骤：使用局部随机生成器以获得确定性序列（避免全局 Seed 弃用）
    rng := rand.New(rand.NewSource(12345))
    in := make([]int, 0, 1000)
    for i := 0; i < 1000; i++ { in = append(in, rng.Intn(50)) }
    out := Unique(in)

    // 计算期望：记录首次出现位置并按该顺序输出
    seen := make(map[int]bool, 50)
    expect := make([]int, 0, 50)
    for _, v := range in {
        if !seen[v] { seen[v] = true; expect = append(expect, v) }
    }
    if len(out) != len(expect) { t.Fatalf("长度不一致: %d != %d", len(out), len(expect)) }
    for i := range expect {
        if out[i] != expect[i] { t.Fatalf("顺序不稳定: idx=%d got=%d want=%d", i, out[i], expect[i]) }
    }
}

// TestFilterEdges 测试 Filter 在全保留/全丢弃与随机选择场景
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：分别验证长度与元素一致性
func TestFilterEdges(t *testing.T) {
    in := []int{1,2,3,4,5}
    // 全保留
    all := Filter(in, func(x int) bool { return true })
    if len(all) != len(in) { t.Fatalf("全保留长度应一致: %d != %d", len(all), len(in)) }
    for i := range in { if all[i] != in[i] { t.Fatalf("全保留元素不一致") } }
    // 全丢弃
    none := Filter(in, func(x int) bool { return false })
    if len(none) != 0 { t.Fatalf("全丢弃应为空: got=%d", len(none)) }
    // 随机选择（固定seed，使用局部随机生成器）
    rng := rand.New(rand.NewSource(7))
    sel := Filter(in, func(x int) bool { return rng.Intn(2) == 0 })
    // 长度应在[0,len(in)]且元素来自原集合
    if len(sel) > len(in) { t.Fatalf("随机过滤长度异常: %d", len(sel)) }
    for _, v := range sel {
        found := false
        for _, w := range in { if v == w { found = true; break } }
        if !found { t.Fatalf("过滤结果出现非原集合元素: %d", v) }
    }
}

// TestIndexOfBoundaries 测试 IndexOf 在空切片与重复元素时的行为
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：空切片返回-1；重复元素应返回首次出现位置
func TestIndexOfBoundaries(t *testing.T) {
    if idx := IndexOf([]int{}, 1); idx != -1 { t.Fatalf("空切片应返回-1: %d", idx) }
    if idx := IndexOf([]int{9,9,1,9}, 9); idx != 0 { t.Fatalf("重复元素应返回首次索引: %d", idx) }
}

// TestKeysValuesEmpty 测试 Keys/Values 在空map时的返回
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：应返回空切片（非nil长度为0）
func TestKeysValuesEmpty(t *testing.T) {
    m := map[string]int{}
    ks := Keys(m)
    vs := Values(m)
    if len(ks) != 0 || len(vs) != 0 { t.Fatalf("空map键值长度应为0: ks=%d vs=%d", len(ks), len(vs)) }
}

// TestGetOrDefaultMissing 测试 GetOrDefault 缺失键返回默认值
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：存在键返回真实值，缺失键返回默认值
func TestGetOrDefaultMissing(t *testing.T) {
    m := map[string]int{"x": 1}
    if v := GetOrDefault(m, "x", 9); v != 1 { t.Fatalf("存在键应返回真实值: %d", v) }
    if v := GetOrDefault(m, "y", 9); v != 9 { t.Fatalf("缺失键应返回默认值: %d", v) }
}