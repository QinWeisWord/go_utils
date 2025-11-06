package collections

import (
    "reflect"
    "sort"
    "testing"
)

// TestContains 测试 Contains 是否正确判断切片包含元素
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：分别测试存在与不存在、空切片与自定义类型
func TestContains(t *testing.T) {
    // 基本类型
    if !Contains([]int{1, 2, 3}, 2) {
        t.Fatalf("期望包含2")
    }
    if Contains([]int{1, 2, 3}, 9) {
        t.Fatalf("不应包含9")
    }
    // 空切片
    if Contains([]string{}, "a") {
        t.Fatalf("空切片不应包含任何元素")
    }
    // 自定义可比较类型
    type ID int
    if !Contains([]ID{1, 5, 9}, ID(5)) {
        t.Fatalf("期望包含ID(5)")
    }
}

// TestMap 测试 Map 是否正确映射切片元素
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：对整数切片平方映射，并验证输出
func TestMap(t *testing.T) {
    in := []int{1, 2, 3, 4}
    out := Map(in, func(x int) int { return x * x })
    expect := []int{1, 4, 9, 16}
    if !reflect.DeepEqual(out, expect) {
        t.Fatalf("映射结果不一致: got=%v expect=%v", out, expect)
    }
}

// TestKeysValues 测试 Keys 与 Values 是否正确返回键与值
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：构造map后取键和值并排序比较
func TestKeysValues(t *testing.T) {
    m := map[string]int{"a": 1, "b": 2, "c": 3}
    ks := Keys(m)
    vs := Values(m)
    sort.Strings(ks)
    sort.Ints(vs)
    if !reflect.DeepEqual(ks, []string{"a", "b", "c"}) {
        t.Fatalf("键集合不一致: %v", ks)
    }
    if !reflect.DeepEqual(vs, []int{1, 2, 3}) {
        t.Fatalf("值集合不一致: %v", vs)
    }
}

// TestMerge 测试 Merge 是否正确合并两个map
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：b 应覆盖 a 的同名键，同时保留双方键
func TestMerge(t *testing.T) {
    a := map[string]int{"x": 1, "y": 2}
    b := map[string]int{"y": 9, "z": 3}
    m := Merge(a, b)
    expect := map[string]int{"x": 1, "y": 9, "z": 3}
    if !reflect.DeepEqual(m, expect) {
        t.Fatalf("合并结果不一致: got=%v expect=%v", m, expect)
    }
    // 关键步骤：不应修改原map
    if a["y"] != 2 {
        t.Fatalf("原map a 不应被修改")
    }
}