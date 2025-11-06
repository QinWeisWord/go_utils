package strutil

import (
    "reflect"
    "testing"
)

// TestSplitTrimNonEmpty 基本行为测试
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：去除空白与过滤空项
func TestSplitTrimNonEmpty(t *testing.T) {
    s := " a , b ,  , c ,  "
    got := SplitTrimNonEmpty(s, ",")
    want := []string{"a", "b", "c"}
    if !reflect.DeepEqual(got, want) {
        t.Fatalf("SplitTrimNonEmpty 失败: got=%v want=%v", got, want)
    }
}