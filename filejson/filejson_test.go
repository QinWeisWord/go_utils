package filejson

import (
    "path/filepath"
    "strings"
    "testing"
)

// TestExistsAndIsDir 测试 Exists 与 IsDir 在目录与文件场景下的行为
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：使用 t.TempDir 创建临时目录与文件并断言
func TestExistsAndIsDir(t *testing.T) {
    dir := t.TempDir()
    if !Exists(dir) || !IsDir(dir) {
        t.Fatalf("临时目录应存在且为目录")
    }
    f := filepath.Join(dir, "a.txt")
    if err := WriteFile(f, "hi"); err != nil {
        t.Fatalf("写入文件失败: %v", err)
    }
    if !Exists(f) || IsDir(f) {
        t.Fatalf("文件应存在且不是目录")
    }
}

// TestReadWriteCopyList 测试 ReadFile/WriteFile/CopyFile/ListDir 的基本功能
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：写入源文件、复制到目标、读取内容、列出目录
func TestReadWriteCopyList(t *testing.T) {
    dir := t.TempDir()
    src := filepath.Join(dir, "src.txt")
    dst := filepath.Join(dir, "sub", "dst.txt")
    content := "hello world"

    if err := WriteFile(src, content); err != nil {
        t.Fatalf("WriteFile 失败: %v", err)
    }
    if s, err := ReadFile(src); err != nil || s != content {
        t.Fatalf("ReadFile 失败或内容不一致: s=%q err=%v", s, err)
    }
    if err := CopyFile(src, dst); err != nil {
        t.Fatalf("CopyFile 失败: %v", err)
    }
    if s, err := ReadFile(dst); err != nil || s != content {
        t.Fatalf("复制后读取失败或内容不一致: s=%q err=%v", s, err)
    }

    names, err := ListDir(dir)
    if err != nil { t.Fatalf("ListDir 失败: %v", err) }
    // 关键步骤：应包含 src.txt 与 sub 目录
    joined := strings.Join(names, ",")
    if !strings.Contains(joined, "src.txt") || !strings.Contains(joined, "sub") {
        t.Fatalf("ListDir 未包含预期项: %v", names)
    }
}

// TestJSON 编解码测试 ToJSON/ToPrettyJSON/FromJSON 泛型
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：对结构体编码到JSON后再解析回来
func TestJSON(t *testing.T) {
    type User struct{
        ID int    `json:"id"`
        Name string `json:"name"`
    }
    u := User{ID: 7, Name: "Tom"}
    js, err := ToJSON(u)
    if err != nil || !strings.Contains(js, "\"id\":7") || !strings.Contains(js, "\"name\":\"Tom\"") {
        t.Fatalf("ToJSON 失败或内容不包含字段: %v", js)
    }
    pjs, err := ToPrettyJSON(u)
    if err != nil || !strings.Contains(pjs, "\n") {
        t.Fatalf("ToPrettyJSON 失败或未美化: %v", err)
    }
    back, err := FromJSON[User](js)
    if err != nil || back.ID != 7 || back.Name != "Tom" {
        t.Fatalf("FromJSON 失败或字段不一致: %v", back)
    }
}