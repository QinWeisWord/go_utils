package kvcache

import (
    "sync"
    "testing"
)

// TestBasicSetGet 测试基础的Set/Get行为
// 参数 t: 测试对象
// 返回值：无
// 关键步骤：设置后应能读取到相同值
func TestBasicSetGet(t *testing.T) {
    c := New[string, int]()
    c.Set("a", 1)
    c.Set("b", 2)

    if v, ok := c.Get("a"); !ok || v != 1 {
        t.Fatalf("读取键a失败: ok=%v v=%d", ok, v)
    }
    if v, ok := c.Get("b"); !ok || v != 2 {
        t.Fatalf("读取键b失败: ok=%v v=%d", ok, v)
    }
    if _, ok := c.Get("c"); ok {
        t.Fatalf("不存在键c不应返回ok=true")
    }
}

// TestOverwriteAndLen 测试覆盖写入与长度统计
// 参数 t: 测试对象
// 返回值：无
// 关键步骤：覆盖同键应更新值且长度保持不变
func TestOverwriteAndLen(t *testing.T) {
    c := New[string, int]()
    c.Set("k", 1)
    if c.Len() != 1 { t.Fatalf("长度应为1: %d", c.Len()) }
    c.Set("k", 9)
    if c.Len() != 1 { t.Fatalf("覆盖后长度仍应为1: %d", c.Len()) }
    if v, _ := c.Get("k"); v != 9 { t.Fatalf("覆盖后值应为9: %d", v) }
}

// TestDeleteHasClear 测试删除、存在判断与清空
// 参数 t: 测试对象
// 返回值：无
// 关键步骤：删除返回true，清空后长度为0
func TestDeleteHasClear(t *testing.T) {
    c := New[string, string]()
    c.Set("x", "A")
    c.Set("y", "B")
    if !c.Has("x") || !c.Has("y") { t.Fatalf("初始键应存在") }
    if !c.Delete("x") { t.Fatalf("删除x应返回true") }
    if c.Has("x") { t.Fatalf("x删除后不应存在") }
    if c.Len() != 1 { t.Fatalf("删除后长度应为1: %d", c.Len()) }
    c.Clear()
    if c.Len() != 0 { t.Fatalf("清空后长度应为0: %d", c.Len()) }
}

// TestKeysValues 测试Keys与Values返回
// 参数 t: 测试对象
// 返回值：无
// 关键步骤：插入多个键后应能返回全部键与值
func TestKeysValues(t *testing.T) {
    c := New[int, string]()
    for i := 0; i < 5; i++ { c.Set(i, string(rune('A'+i))) }
    ks := c.Keys()
    vs := c.Values()
    if len(ks) != 5 || len(vs) != 5 {
        t.Fatalf("键值长度应为5: ks=%d vs=%d", len(ks), len(vs))
    }
}

// TestConcurrentSafety 并发安全性测试
// 参数 t: 测试对象
// 返回值：无
// 关键步骤：多协程同时写入与读取，断言最终长度与值正确
func TestConcurrentSafety(t *testing.T) {
    c := New[int, int]()
    const workers = 32
    const loops = 1000
    var wg sync.WaitGroup
    wg.Add(workers)
    for w := 0; w < workers; w++ {
        go func(base int) {
            defer wg.Done()
            for i := 0; i < loops; i++ {
                // 关键步骤：使用不同键空间，避免频繁覆盖导致长度不稳定
                key := base*loops + i
                c.Set(key, i)
                if v, ok := c.Get(key); !ok || v != i {
                    t.Errorf("并发读取失败: key=%d ok=%v v=%d want=%d", key, ok, v, i)
                }
            }
        }(w)
    }
    wg.Wait()
    expected := workers * loops
    if c.Len() != expected {
        t.Fatalf("最终长度不匹配: got=%d want=%d", c.Len(), expected)
    }
}