package iniutil

import (
    "os"
    "strings"
    "testing"
)

// TestParseBasic 测试基本INI解析能力
// 参数 t: 测试句柄
// 返回值: 无
// 关键步骤：验证注释跳过、section识别、类型获取
func TestParseBasic(t *testing.T) {
    ini := `; 注释行
# 另一条注释

app_name="Demo 应用"
debug=true
[db]
host=localhost
port=3306
timeout=1.5
[paths]
log_dir=/var/log/demo
`
    // 关键步骤：通过Reader加载
    cfg, err := LoadFromReader(strings.NewReader(ini))
    if err != nil {
        t.Fatalf("LoadFromReader error: %v", err)
    }
    if got := cfg.GetString("", "app_name", ""); got != "Demo 应用" {
        t.Fatalf("app_name = %q", got)
    }
    bv, err := cfg.GetBool("", "debug", false)
    if err != nil || !bv {
        t.Fatalf("debug parse, v=%v err=%v", bv, err)
    }
    if got := cfg.GetString("db", "host", ""); got != "localhost" {
        t.Fatalf("db.host = %q", got)
    }
    iv, err := cfg.GetInt("db", "port", 0)
    if err != nil || iv != 3306 {
        t.Fatalf("db.port = %v err=%v", iv, err)
    }
    fv, err := cfg.GetFloat64("db", "timeout", 0)
    if err != nil || fv != 1.5 {
        t.Fatalf("db.timeout = %v err=%v", fv, err)
    }
}

// TestSaveLoadRoundTrip 测试保存与重新加载保持值一致
// 参数 t: 测试句柄
// 返回值: 无
// 关键步骤：写入临时文件后再次加载校验
func TestSaveLoadRoundTrip(t *testing.T) {
    cfg := New()
    // 关键步骤：设置多种类型键
    cfg.Set("", "app_name", "XApp")
    cfg.Set("core", "enabled", "true")
    cfg.Set("core", "workers", "8")
    cfg.Set("core", "ratio", "0.75")

    dir := t.TempDir()
    path := dir + string(os.PathSeparator) + "cfg.ini"
    if err := cfg.SaveToFile(path); err != nil {
        t.Fatalf("SaveToFile error: %v", err)
    }
    // 关键步骤：从文件重新加载
    cfg2, err := LoadFromFile(path)
    if err != nil {
        t.Fatalf("LoadFromFile error: %v", err)
    }
    if s := cfg2.GetString("", "app_name", ""); s != "XApp" {
        t.Fatalf("app_name roundtrip = %q", s)
    }
    b, err := cfg2.GetBool("core", "enabled", false)
    if err != nil || !b {
        t.Fatalf("enabled roundtrip v=%v err=%v", b, err)
    }
    n, err := cfg2.GetInt("core", "workers", 0)
    if err != nil || n != 8 {
        t.Fatalf("workers roundtrip v=%v err=%v", n, err)
    }
    r, err := cfg2.GetFloat64("core", "ratio", 0)
    if err != nil || r != 0.75 {
        t.Fatalf("ratio roundtrip v=%v err=%v", r, err)
    }
}

// TestMergeOverwrite 测试合并逻辑（覆盖与保留）
// 参数 t: 测试句柄
// 返回值: 无
// 关键步骤：分别验证 overwrite=false/true 行为
func TestMergeOverwrite(t *testing.T) {
    a := New()
    a.Set("", "debug", "false")
    a.Set("db", "host", "a.local")
    a.Set("db", "port", "3306")

    b := New()
    b.Set("", "debug", "true")
    b.Set("db", "host", "b.local")
    b.Set("db", "timeout", "2.0")

    outKeep := Merge(a, b, false)
    if v := outKeep.GetString("db", "host", ""); v != "a.local" {
        t.Fatalf("keep host expected a.local got %q", v)
    }
    if bv, _ := outKeep.GetBool("", "debug", false); bv != false {
        t.Fatalf("keep debug expected false got %v", bv)
    }
    if fv, _ := outKeep.GetFloat64("db", "timeout", 0); fv != 2.0 {
        t.Fatalf("keep timeout expected 2.0 got %v", fv)
    }

    outOverwrite := Merge(a, b, true)
    if v := outOverwrite.GetString("db", "host", ""); v != "b.local" {
        t.Fatalf("overwrite host expected b.local got %q", v)
    }
    if bv, _ := outOverwrite.GetBool("", "debug", false); bv != true {
        t.Fatalf("overwrite debug expected true got %v", bv)
    }
}

// TestDeleteAndKeys 测试删除与键列表
// 参数 t: 测试句柄
// 返回值: 无
// 关键步骤：删除后检查Has与Keys
func TestDeleteAndKeys(t *testing.T) {
    c := New()
    c.Set("x", "a", "1")
    c.Set("x", "b", "2")
    c.Delete("x", "a")
    if c.Has("x", "a") {
        t.Fatalf("a should be deleted")
    }
    keys := c.Keys("x")
    if len(keys) != 1 || keys[0] != "b" {
        t.Fatalf("keys expected [b] got %v", keys)
    }
}