package convert

import (
    "reflect"
    "sort"
    "testing"
)

// S 类型实现 fmt.Stringer 接口用于 ToString 测试
// 结构体字段 X: 数值字段
// 关键步骤：提供 String 方法返回格式化字符串
type S struct{ X int }

// String 返回格式化字符串
// 参数: 无
// 返回值: 字符串形式 "S(<X>)"
// 关键步骤：调用 ToStringMust 生成 X 的字符串
func (s S) String() string { return "S(" + ToStringMust(s.X) + ")" }

// TestToInt 测试 ToInt 的常见类型转换与错误分支
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖 string/int64/float64/uint64/bool/非法字符串
func TestToInt(t *testing.T) {
    if v, err := ToInt("42"); v != 42 || err != nil {
        t.Fatalf("string->int 失败: v=%v err=%v", v, err)
    }
    if v, err := ToInt(int64(7)); v != 7 || err != nil {
        t.Fatalf("int64->int 失败: v=%v err=%v", v, err)
    }
    if v, err := ToInt(3.9); v != 3 || err != nil {
        t.Fatalf("float64->int 截断失败: v=%v err=%v", v, err)
    }
    if v, err := ToInt(uint64(5)); v != 5 || err != nil {
        t.Fatalf("uint64->int 失败: v=%v err=%v", v, err)
    }
    if v, err := ToInt(true); v != 1 || err != nil {
        t.Fatalf("bool->int 失败: v=%v err=%v", v, err)
    }
    if _, err := ToInt("x"); err == nil {
        t.Fatalf("非法字符串应返回错误")
    } else {
        if _, ok := err.(*ConvertError); !ok {
            t.Fatalf("错误类型应为*ConvertError: %T", err)
        }
    }
}

// TestToBool 测试 ToBool 的多语言真值/假值与错误分支
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖 "是"/"否"/数字/未知字符串
func TestToBool(t *testing.T) {
    if v, err := ToBool("是"); v != true || err != nil {
        t.Fatalf("\"是\" 应为 true: v=%v err=%v", v, err)
    }
    if v, err := ToBool("否"); v != false || err != nil {
        t.Fatalf("\"否\" 应为 false: v=%v err=%v", v, err)
    }
    if v, err := ToBool(0); v != false || err != nil {
        t.Fatalf("0 应为 false: v=%v err=%v", v, err)
    }
    if _, err := ToBool("maybe"); err == nil {
        t.Fatalf("未知布尔字符串应返回错误")
    }
}

// TestToString 测试 ToString 的常见类型与 Stringer 接口
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖 int/bool/float32/自定义Stringer
func TestToString(t *testing.T) {
    if s, err := ToString(123); s != "123" || err != nil {
        t.Fatalf("int->string 失败: s=%v err=%v", s, err)
    }
    if s, err := ToString(true); s != "true" || err != nil {
        t.Fatalf("bool->string 失败: s=%v err=%v", s, err)
    }
    if s, err := ToString(float32(1.25)); s != "1.25" || err != nil {
        t.Fatalf("float32->string 失败: s=%v err=%v", s, err)
    }
    if s, err := ToString(S{X: 10}); s != "S(10)" || err != nil {
        t.Fatalf("Stringer->string 失败: s=%v err=%v", s, err)
    }
}

// ToStringMust 测试辅助：将任意值转换为字符串（失败直接t.Fatal）
// 参数 v: 任意值
// 返回值: 字符串（失败时在测试中终止）
// 关键步骤：封装 ToString 用于构造期望字符串
func ToStringMust(v interface{}) string {
    s, err := ToString(v)
    if err != nil {
        // 关键步骤：此处仅用于测试辅助，panic以中止当前用例
        panic(err)
    }
    return s
}

// TestToBinString 测试 ToBinString 的多类型格式化
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖 int/uint/float/bool/负数
func TestToBinString(t *testing.T) {
    if s, err := ToBinString(10); s != "1010" || err != nil {
        t.Fatalf("10 的二进制应为 1010: s=%v err=%v", s, err)
    }
    if s, err := ToBinString(uint8(3)); s != "11" || err != nil {
        t.Fatalf("3 的二进制应为 11: s=%v err=%v", s, err)
    }
    if s, err := ToBinString(2.5); s != "10" || err != nil {
        t.Fatalf("2.5 转为整数后应为 10: s=%v err=%v", s, err)
    }
    if s, err := ToBinString(true); s != "1" || err != nil {
        t.Fatalf("true 的二进制应为 1: s=%v err=%v", s, err)
    }
    if s, err := ToBinString(-2); s != "-10" || err != nil {
        t.Fatalf("-2 的二进制应为 -10: s=%v err=%v", s, err)
    }
}

// TestUniqueSlice 测试 UniqueSlice 的去重与保持顺序
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖可比较与不可比较元素
func TestUniqueSlice(t *testing.T) {
    in := []int{1, 2, 1, 3, 2, 4}
    outAny, err := UniqueSlice(in)
    if err != nil { t.Fatalf("int 去重失败: %v", err) }
    out := outAny.([]int)
    expect := []int{1, 2, 3, 4}
    if !reflect.DeepEqual(out, expect) {
        t.Fatalf("去重结果不一致: %v != %v", out, expect)
    }
    // 不可比较元素（切片）
    type Box struct{ A []int }
    in2 := []Box{{A: []int{1}}, {A: []int{1}}, {A: []int{2}}}
    outAny2, err := UniqueSlice(in2)
    if err != nil { t.Fatalf("struct 切片去重失败: %v", err) }
    out2 := outAny2.([]Box)
    if len(out2) != 2 {
        t.Fatalf("去重后长度应为2: got=%d", len(out2))
    }
}

// TestUniqueSliceByField 测试 UniqueSliceByField 的结构体字段与标签匹配
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖字段名、json标签、nil指针与缺失键
func TestUniqueSliceByField(t *testing.T) {
    type User struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    }
    in := []User{{ID: 1, Name: "a"}, {ID: 2, Name: "b"}, {ID: 1, Name: "c"}}
    outAny, err := UniqueSliceByField(in, "ID", "")
    if err != nil { t.Fatalf("按字段去重失败: %v", err) }
    out := outAny.([]User)
    expect := []User{{ID: 1, Name: "a"}, {ID: 2, Name: "b"}}
    if !reflect.DeepEqual(out, expect) {
        t.Fatalf("按字段去重结果不一致: %v != %v", out, expect)
    }

    // 按标签去重
    outAny2, err := UniqueSliceByField(in, "name", "json")
    if err != nil { t.Fatalf("按标签去重失败: %v", err) }
    out2 := outAny2.([]User)
    expect2 := []User{{ID: 1, Name: "a"}, {ID: 2, Name: "b"}, {ID: 1, Name: "c"}}
    if !reflect.DeepEqual(out2, expect2) {
        t.Fatalf("按标签去重结果不应去重不同name: %v", out2)
    }

    // 指针与nil特殊键
    var pnil *User = nil
    inPtrs := []*User{pnil, pnil, {ID: 1, Name: "a"}, {ID: 1, Name: "b"}}
    outAny3, err := UniqueSliceByField(inPtrs, "ID", "")
    if err != nil { t.Fatalf("指针切片去重失败: %v", err) }
    out3 := outAny3.([]*User)
    if len(out3) != 2 {
        t.Fatalf("指针切片去重后长度应为2: %d", len(out3))
    }
}

// TestMapConversions 测试数组与map互转
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：验证 ToMapFromSlice、ToSliceFromMapValues、ToSliceFromMapKeys
func TestMapConversions(t *testing.T) {
    // 数组->map
    arr := []string{"x", "y"}
    m, err := ToMapFromSlice(arr)
    if err != nil { t.Fatalf("数组转map失败: %v", err) }
    if m[0] != "x" || m[1] != "y" {
        t.Fatalf("索引映射错误: %v", m)
    }
    // map值->切片
    mv := map[string]int{"a": 1, "b": 2}
    vals, err := ToSliceFromMapValues(mv)
    if err != nil { t.Fatalf("map值转切片失败: %v", err) }
    sort.Slice(vals, func(i, j int) bool { return vals[i].(int) < vals[j].(int) })
    if !reflect.DeepEqual(vals, []interface{}{1, 2}) {
        t.Fatalf("值集合不一致: %v", vals)
    }
    // map键->切片
    keys, err := ToSliceFromMapKeys(mv)
    if err != nil { t.Fatalf("map键转切片失败: %v", err) }
    sort.Slice(keys, func(i, j int) bool { return keys[i].(string) < keys[j].(string) })
    if !reflect.DeepEqual(keys, []interface{}{"a", "b"}) {
        t.Fatalf("键集合不一致: %v", keys)
    }
}