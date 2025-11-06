package iniutil

import (
    "io"
    "os"
    "path/filepath"
    "strings"
    "testing"
    "go_utils/strutil"
)

// TestInlineCommentAndColon 测试内联注释与冒号分隔支持
// 参数 t: 测试句柄
// 返回值: 无
// 关键步骤：开启 InlineComment 与 AllowColon 选项
func TestInlineCommentAndColon(t *testing.T) {
    ini := `name = Alice ; 用户名注释
age: 30 # 年龄注释
`
    opt := ParseOptions{InlineComment: true, AllowColon: true}
    cfg, err := LoadFromReaderWithOptions(strings.NewReader(ini), opt)
    if err != nil { t.Fatalf("parse error: %v", err) }
    if v := cfg.GetString("", "name", ""); v != "Alice" { t.Fatalf("name=%q", v) }
    n, err := cfg.GetInt("", "age", 0)
    if err != nil || n != 30 { t.Fatalf("age=%v err=%v", n, err) }
}

// TestMultilineValue 测试多行值拼接
// 参数 t: 测试句柄
// 返回值: 无
// 关键步骤：行末'\\'续行，以'\n'连接
func TestMultilineValue(t *testing.T) {
    ini := `desc = first line\\
second line\\
third line
`
    opt := ParseOptions{AllowMultiline: true}
    cfg, err := LoadFromReaderWithOptions(strings.NewReader(ini), opt)
    if err != nil { t.Fatalf("parse error: %v", err) }
    want := "first line\nsecond line\nthird line"
    if v := cfg.GetString("", "desc", ""); v != want { t.Fatalf("desc=%q", v) }
}

// TestIncludeBasic 测试包含文件解析
// 参数 t: 测试句柄
// 返回值: 无
// 关键步骤：使用默认基路径解析 .include 指令
func TestIncludeBasic(t *testing.T) {
    dir := t.TempDir()
    inc := "[db]\nhost=inc.local\n"
    incPath := filepath.Join(dir, "inc.ini")
    if err := os.WriteFile(incPath, []byte(inc), 0644); err != nil { t.Fatalf("write inc: %v", err) }
    main := ".include \"inc.ini\"\n[db]\nport=3306\n"
    mainPath := filepath.Join(dir, "main.ini")
    if err := os.WriteFile(mainPath, []byte(main), 0644); err != nil { t.Fatalf("write main: %v", err) }
    cfg, err := LoadFromFileWithOptions(mainPath, ParseOptions{IncludeOverwrite: false})
    if err != nil { t.Fatalf("load error: %v", err) }
    if h := cfg.GetString("db", "host", ""); h != "inc.local" { t.Fatalf("host=%q", h) }
    if p, _ := cfg.GetInt("db", "port", 0); p != 3306 { t.Fatalf("port=%v", p) }
}

// TestIncludeWithResolver 测试自定义包含解析器
// 参数 t: 测试句柄
// 返回值: 无
// 关键步骤：通过 IncludeResolver 提供读取流
func TestIncludeWithResolver(t *testing.T) {
    inc := "[x]\nname=via-resolver\n"
    resolver := func(base, path string) (io.ReadCloser, error) {
        return io.NopCloser(strings.NewReader(inc)), nil
    }
    opt := ParseOptions{IncludeResolver: resolver, IncludeOverwrite: true}
    cfg, err := LoadFromReaderWithOptions(strings.NewReader("!include dummy\n[x]\nflag=true\n"), opt)
    if err != nil { t.Fatalf("parse error: %v", err) }
    if v := cfg.GetString("x", "name", ""); v != "via-resolver" { t.Fatalf("name=%q", v) }
    b, _ := cfg.GetBool("x", "flag", false)
    if !b { t.Fatalf("flag should be true") }
}

// TestInterpolation 测试占位符插值
// 参数 t: 测试句柄
// 返回值: 无
// 关键步骤：启用 EnableInterpolation，支持同区段与跨区段引用
func TestInterpolation(t *testing.T) {
    ini := "[db]\nhost=local\nurl=postgres://${db.host}:5432\n[svc]\nendpoint=http://${db.host}:8080\n"
    opt := ParseOptions{EnableInterpolation: true}
    cfg, err := LoadFromReaderWithOptions(strings.NewReader(ini), opt)
    if err != nil { t.Fatalf("parse error: %v", err) }
    if v := cfg.GetString("db", "url", ""); v != "postgres://local:5432" { t.Fatalf("url=%q", v) }
    if e := cfg.GetString("svc", "endpoint", ""); e != "http://local:8080" { t.Fatalf("endpoint=%q", e) }
}

// TestDuplicateAppendAndGetStrings 测试重复键拼接与切片读取
// 参数 t: 测试句柄
// 返回值: 无
// 关键步骤：AppendDuplicateKeys开启后逗号拼接，GetStrings按','拆分
func TestDuplicateAppendAndGetStrings(t *testing.T) {
    ini := "[list]\nitems=a\nitems=b\nitems=c\n"
    opt := ParseOptions{AppendDuplicateKeys: true}
    cfg, err := LoadFromReaderWithOptions(strings.NewReader(ini), opt)
    if err != nil { t.Fatalf("parse error: %v", err) }
    // 关键步骤：读取字符串后使用 strutil 进行拆分
    raw := cfg.GetString("list", "items", "")
    arr := strutil.SplitTrimNonEmpty(raw, ",")
    if len(arr) != 3 || arr[0] != "a" || arr[1] != "b" || arr[2] != "c" {
        t.Fatalf("items=%v", arr)
    }
}

// TestSetStrings 测试 SetStrings 与 GetStrings 配套使用
// 参数 t: 测试句柄
// 返回值: 无
// 关键步骤：SetStrings 写入后应可通过 GetStrings 读取到同样的元素
func TestSetStrings(t *testing.T) {
    c := New()
    // 关键步骤：使用 strutil.Join 将切片连接后写入配置
    c.Set("arr", "vals", strutil.Join([]string{"a", "b", "c"}, ","))
    raw := c.GetString("arr", "vals", "")
    arr := strutil.SplitTrimNonEmpty(raw, ",")
    if len(arr) != 3 || arr[0] != "a" || arr[1] != "b" || arr[2] != "c" {
        t.Fatalf("vals=%v", arr)
    }
}