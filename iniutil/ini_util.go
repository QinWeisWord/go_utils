package iniutil

import (
    "bufio"
    "errors"
    "io"
    "os"
    "sort"
    "strconv"
    "strings"
)

// Config 表示一个INI配置对象
// 结构体字段 data: 内部数据结构，按 section→key→value 存储
// 关键步骤：使用嵌套map以便快速读写与合并
type Config struct {
    data map[string]map[string]string
}

// New 创建一个空的INI配置对象
// 参数: 无
// 返回值: 新的配置对象指针
// 关键步骤：初始化嵌套map
func New() *Config {
    return &Config{data: make(map[string]map[string]string)}
}

// LoadFromReader 从Reader解析INI文本
// 参数 r: 输入流Reader
// 返回值: 配置对象与错误；若解析成功错误为nil
// 关键步骤：逐行读取，支持 [section] 与 key=value，忽略以 ';' 或 '#' 开头的注释
func LoadFromReader(r io.Reader) (*Config, error) {
    cfg := New()
    scanner := bufio.NewScanner(r)
    section := ""
    lineNo := 0
    for scanner.Scan() {
        lineNo++
        raw := scanner.Text()
        // 关键步骤：去除UTF-8 BOM（仅首行可能出现）
        if lineNo == 1 {
            raw = strings.TrimPrefix(raw, "\uFEFF")
        }
        s := strings.TrimSpace(raw)
        if s == "" { continue }
        if strings.HasPrefix(s, ";") || strings.HasPrefix(s, "#") { continue }
        // 关键步骤：处理Section行
        if s[0] == '[' && strings.HasSuffix(s, "]") {
            name := strings.TrimSpace(s[1:len(s)-1])
            if name == "" { name = "" }
            section = name
            if _, ok := cfg.data[section]; !ok {
                cfg.data[section] = make(map[string]string)
            }
            continue
        }
        // 关键步骤：处理 key=value 行
        eq := strings.IndexByte(s, '=')
        if eq <= 0 { // 无效行：没有'='或key为空
            continue
        }
        key := strings.TrimSpace(s[:eq])
        val := strings.TrimSpace(s[eq+1:])
        // 关键步骤：支持双引号包裹的值
        if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
            val = strings.TrimPrefix(strings.TrimSuffix(val, "\""), "\"")
        }
        if _, ok := cfg.data[section]; !ok {
            cfg.data[section] = make(map[string]string)
        }
        cfg.data[section][key] = val
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    return cfg, nil
}

// LoadFromFile 从文件路径加载INI配置
// 参数 path: 文件路径
// 返回值: 配置对象与错误
// 关键步骤：打开文件并委托给 LoadFromReader
func LoadFromFile(path string) (*Config, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    return LoadFromReader(f)
}

// SaveToWriter 将配置保存到Writer（写为标准INI文本）
// 参数 w: 输出流Writer
// 返回值: 错误；成功时为nil
// 关键步骤：按section与key排序以获得稳定输出
func (c *Config) SaveToWriter(w io.Writer) error {
    if c == nil { return errors.New("nil config") }
    // 关键步骤：收集并排序section
    secs := make([]string, 0, len(c.data))
    for s := range c.data { secs = append(secs, s) }
    sort.Strings(secs)
    bw := bufio.NewWriter(w)
    for _, s := range secs {
        if s != "" {
            if _, err := bw.WriteString("[" + s + "]\n"); err != nil { return err }
        }
        kv := c.data[s]
        // 关键步骤：排序key以稳定写出
        keys := make([]string, 0, len(kv))
        for k := range kv { keys = append(keys, k) }
        sort.Strings(keys)
        for _, k := range keys {
            v := kv[k]
            // 关键步骤：若包含空格或特殊字符，使用双引号包裹
            if strings.ContainsAny(v, "\t \";#[]=") {
                v = "\"" + v + "\""
            }
            if _, err := bw.WriteString(k + "=" + v + "\n"); err != nil { return err }
        }
        // 关键步骤：section间空行分隔
        if _, err := bw.WriteString("\n"); err != nil { return err }
    }
    if err := bw.Flush(); err != nil { return err }
    return nil
}

// SaveToFile 将配置保存到文件路径
// 参数 path: 文件路径
// 返回值: 错误；成功时为nil
// 关键步骤：创建/覆盖文件并委托给 SaveToWriter
func (c *Config) SaveToFile(path string) error {
    f, err := os.Create(path)
    if err != nil { return err }
    defer f.Close()
    return c.SaveToWriter(f)
}

// GetString 获取字符串值（若缺失返回默认值）
// 参数 section: 区段名称（空字符串表示默认区段）
// 参数 key: 键名称
// 参数 def: 缺失时返回的默认值
// 返回值: 字符串值；缺失返回def
// 关键步骤：安全读取嵌套map
func (c *Config) GetString(section, key, def string) string {
    if c == nil { return def }
    if kv, ok := c.data[section]; ok {
        if v, ok2 := kv[key]; ok2 { return v }
    }
    return def
}

// GetInt 获取整数值（缺失或解析失败返回默认值并附带错误）
// 参数 section: 区段名称
// 参数 key: 键名称
// 参数 def: 缺失或解析失败时返回的默认值
// 返回值: 整数值与错误；成功时错误为nil
// 关键步骤：使用strconv.Atoi解析
func (c *Config) GetInt(section, key string, def int) (int, error) {
    v := c.GetString(section, key, "")
    if v == "" { return def, nil }
    n, err := strconv.Atoi(strings.TrimSpace(v))
    if err != nil { return def, err }
    return n, nil
}

// GetFloat64 获取浮点值（缺失或解析失败返回默认值并附带错误）
// 参数 section: 区段名称
// 参数 key: 键名称
// 参数 def: 缺失或解析失败时返回的默认值
// 返回值: 浮点值与错误；成功时错误为nil
// 关键步骤：使用strconv.ParseFloat解析
func (c *Config) GetFloat64(section, key string, def float64) (float64, error) {
    v := c.GetString(section, key, "")
    if v == "" { return def, nil }
    f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
    if err != nil { return def, err }
    return f, nil
}

// GetBool 获取布尔值（支持 true/false/yes/no/on/off/1/0）
// 参数 section: 区段名称
// 参数 key: 键名称
// 参数 def: 缺失或解析失败时返回的默认值
// 返回值: 布尔值与错误；成功时错误为nil
// 关键步骤：统一小写后匹配常见布尔字面量
func (c *Config) GetBool(section, key string, def bool) (bool, error) {
    v := strings.ToLower(strings.TrimSpace(c.GetString(section, key, "")))
    if v == "" { return def, nil }
    switch v {
    case "true", "yes", "on", "1":
        return true, nil
    case "false", "no", "off", "0":
        return false, nil
    default:
        return def, errors.New("invalid bool value: " + v)
    }
}

// Set 设置键值（若区段或键不存在则创建）
// 参数 section: 区段名称
// 参数 key: 键名称
// 参数 value: 值字符串
// 返回值: 无
// 关键步骤：确保区段存在后赋值
func (c *Config) Set(section, key, value string) {
    if _, ok := c.data[section]; !ok {
        c.data[section] = make(map[string]string)
    }
    c.data[section][key] = value
}

// Delete 删除指定键（若键或区段不存在则忽略）
// 参数 section: 区段名称
// 参数 key: 键名称
// 返回值: 无
// 关键步骤：安全检查后删除键
func (c *Config) Delete(section, key string) {
    if kv, ok := c.data[section]; ok {
        delete(kv, key)
    }
}

// Sections 返回所有区段名称（排序）
// 参数: 无
// 返回值: 已排序的区段名称切片
// 关键步骤：收集并排序
func (c *Config) Sections() []string {
    secs := make([]string, 0, len(c.data))
    for s := range c.data { secs = append(secs, s) }
    sort.Strings(secs)
    return secs
}

// Keys 返回指定区段的所有键（排序）
// 参数 section: 区段名称
// 返回值: 已排序的键切片
// 关键步骤：收集并排序，缺失区段返回空切片
func (c *Config) Keys(section string) []string {
    kv, ok := c.data[section]
    if !ok { return []string{} }
    keys := make([]string, 0, len(kv))
    for k := range kv { keys = append(keys, k) }
    sort.Strings(keys)
    return keys
}

// Has 判断是否存在指定键
// 参数 section: 区段名称
// 参数 key: 键名称
// 返回值: 布尔值；存在返回true
// 关键步骤：安全读取嵌套map
func (c *Config) Has(section, key string) bool {
    if kv, ok := c.data[section]; ok {
        _, ok2 := kv[key]
        return ok2
    }
    return false
}

// Merge 合并两个配置并返回新配置
// 参数 a: 配置A（作为基础）
// 参数 b: 配置B（作为增量）
// 参数 overwrite: 当键冲突时是否以B覆盖A（true覆盖；false保留A）
// 返回值: 新的合并后配置对象
// 关键步骤：逐section与key合并，遵循覆盖策略
func Merge(a, b *Config, overwrite bool) *Config {
    out := New()
    // 关键步骤：复制A
    if a != nil {
        for s, kv := range a.data {
            if _, ok := out.data[s]; !ok { out.data[s] = make(map[string]string) }
            for k, v := range kv { out.data[s][k] = v }
        }
    }
    // 关键步骤：合并B
    if b != nil {
        for s, kv := range b.data {
            if _, ok := out.data[s]; !ok { out.data[s] = make(map[string]string) }
            for k, v := range kv {
                if overwrite || !out.Has(s, k) {
                    out.data[s][k] = v
                }
            }
        }
    }
    return out
}