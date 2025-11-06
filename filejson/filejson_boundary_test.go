package filejson

import (
    "path/filepath"
    "testing"
)

// TestFromJSONEmpty 测试 FromJSON 对空JSON的错误分支
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：空字符串解析应返回错误
func TestFromJSONEmpty(t *testing.T) {
    type U struct{ X int }
    if _, err := FromJSON[U](""); err == nil {
        t.Fatalf("空JSON应返回错误")
    }
}

// TestReadCopyMissing 测试 ReadFile/CopyFile 对缺失文件的错误分支
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：读取缺失文件与复制缺失源应返回错误
func TestReadCopyMissing(t *testing.T) {
    dir := t.TempDir()
    f := filepath.Join(dir, "missing.txt")
    if _, err := ReadFile(f); err == nil {
        t.Fatalf("读取缺失文件应返回错误")
    }
    src := filepath.Join(dir, "no_src.txt")
    dst := filepath.Join(dir, "dst.txt")
    if err := CopyFile(src, dst); err == nil {
        t.Fatalf("复制缺失源文件应返回错误")
    }
}