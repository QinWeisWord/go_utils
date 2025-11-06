package iniutil

import (
    "bufio"
    "errors"
    "io"
    "os"
    "path/filepath"
    "strings"
)

// ParseOptions 高级解析选项
// 结构体字段解释：
// - InlineComment: 是否支持值后的内联注释（非引号内的 ';' 或 '#' 起始内容被忽略）
// - AllowColon: 是否允许使用 ':' 作为键值分隔符（与 '=' 并存，取第一个分隔符）
// - AllowMultiline: 是否允许使用行末 '\\' 进行多行值拼接（以'\n'连接）
// - IncludeResolver: 包含文件解析器，入参为基路径与包含路径，返回可读的内容流
// - BaseDir: 当前主文件的基路径，用于默认包含解析
// - IncludeOverwrite: 包含文件与当前配置合并时是否覆盖同名键
// - EnableInterpolation: 是否启用占位符插值（支持 ${key} 与 ${section.key}）
// - AppendDuplicateKeys: 是否在同一section下遇到重复键时以逗号拼接而非覆盖
type ParseOptions struct {
    InlineComment       bool
    AllowColon          bool
    AllowMultiline      bool
    IncludeResolver     func(baseDir, includePath string) (io.ReadCloser, error)
    BaseDir             string
    IncludeOverwrite    bool
    EnableInterpolation bool
    AppendDuplicateKeys bool
}

// LoadFromReaderWithOptions 从Reader解析INI文本（高级选项版）
// 参数 r: 输入流Reader
// 参数 opt: 解析选项
// 返回值: 配置对象与错误；若解析成功错误为nil
// 关键步骤：支持内联注释、冒号分隔、多行值与包含文件
func LoadFromReaderWithOptions(r io.Reader, opt ParseOptions) (*Config, error) {
    cfg := New()
    scanner := bufio.NewScanner(r)
    section := ""
    lineNo := 0
    inMultiline := false
    multiKey := ""
    multiVal := ""

    for scanner.Scan() {
        lineNo++
        raw := scanner.Text()
        if lineNo == 1 {
            raw = strings.TrimPrefix(raw, "\uFEFF")
        }
        s := strings.TrimSpace(raw)
        if s == "" { continue }
        if strings.HasPrefix(s, ";") || strings.HasPrefix(s, "#") { continue }

        // 关键步骤：处理多行值续行
        if inMultiline {
            seg := strings.TrimRight(s, "\\")
            multiVal += "\n" + seg
            if strings.HasSuffix(s, "\\") {
                // 继续多行模式
                continue
            }
            // 结束多行，写入键值
            if _, ok := cfg.data[section]; !ok { cfg.data[section] = make(map[string]string) }
            if old, ok := cfg.data[section][multiKey]; ok && opt.AppendDuplicateKeys {
                cfg.data[section][multiKey] = strings.TrimSpace(old) + "," + strings.TrimSpace(multiVal)
            } else {
                cfg.data[section][multiKey] = strings.TrimSpace(multiVal)
            }
            inMultiline = false
            multiKey, multiVal = "", ""
            continue
        }

        // 关键步骤：处理包含指令 .include 或 !include
        if strings.HasPrefix(s, ".include ") || strings.HasPrefix(s, "!include ") {
            arg := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(s, ".include "), "!include "))
            if len(arg) >= 2 && arg[0] == '"' && arg[len(arg)-1] == '"' {
                arg = strings.TrimSuffix(strings.TrimPrefix(arg, "\""), "\"")
            }
            var rc io.ReadCloser
            var err error
            if opt.IncludeResolver != nil {
                rc, err = opt.IncludeResolver(opt.BaseDir, arg)
            } else {
                full := arg
                if opt.BaseDir != "" && !filepath.IsAbs(arg) {
                    full = filepath.Join(opt.BaseDir, arg)
                }
                f, e := os.Open(full)
                if e != nil { return nil, e }
                rc = f
            }
            includeOpt := opt
            // 关键步骤：若可解析到绝对路径，则更新基路径
            if f, ok := rc.(*os.File); ok {
                dir := filepath.Dir(f.Name())
                includeOpt.BaseDir = dir
            }
            incCfg, err := LoadFromReaderWithOptions(rc, includeOpt)
            rc.Close()
            if err != nil { return nil, err }
            cfg = Merge(cfg, incCfg, opt.IncludeOverwrite)
            continue
        }

        // 关键步骤：Section识别
        if s[0] == '[' && strings.HasSuffix(s, "]") {
            name := strings.TrimSpace(s[1:len(s)-1])
            section = name
            if _, ok := cfg.data[section]; !ok { cfg.data[section] = make(map[string]string) }
            continue
        }

        // 关键步骤：键值分隔，支持 '=' 与可选 ':'
        idx := strings.IndexByte(s, '=')
        colonIdx := -1
        if opt.AllowColon {
            colonIdx = strings.IndexByte(s, ':')
        }
        if idx == -1 || (colonIdx != -1 && colonIdx < idx) {
            idx = colonIdx
        }
        if idx <= 0 { // 无效行
            continue
        }
        key := strings.TrimSpace(s[:idx])
        val := strings.TrimSpace(s[idx+1:])

        // 关键步骤：内联注释剥离（非引号内的 ';' '#'）
        if opt.InlineComment {
            val = stripInlineComment(val)
        }
        // 关键步骤：处理双引号包裹的值
        if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
            val = strings.TrimSuffix(strings.TrimPrefix(val, "\""), "\"")
        }
        // 关键步骤：多行值启用
        if opt.AllowMultiline && strings.HasSuffix(s, "\\") {
            inMultiline = true
            multiKey = key
            multiVal = strings.TrimRight(val, "\\")
            if _, ok := cfg.data[section]; !ok { cfg.data[section] = make(map[string]string) }
            continue
        }
        if _, ok := cfg.data[section]; !ok { cfg.data[section] = make(map[string]string) }
        if old, ok := cfg.data[section][key]; ok && opt.AppendDuplicateKeys {
            cfg.data[section][key] = strings.TrimSpace(old) + "," + strings.TrimSpace(val)
        } else {
            cfg.data[section][key] = val
        }
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    // 关键步骤：插值处理
    if opt.EnableInterpolation {
        if err := cfg.Interpolate(); err != nil {
            return nil, err
        }
    }
    return cfg, nil
}

// LoadFromFileWithOptions 从文件路径解析INI（高级选项版）
// 参数 path: 文件路径
// 参数 opt: 解析选项
// 返回值: 配置对象与错误
// 关键步骤：设置BaseDir并委托给 LoadFromReaderWithOptions
func LoadFromFileWithOptions(path string, opt ParseOptions) (*Config, error) {
    f, err := os.Open(path)
    if err != nil { return nil, err }
    defer f.Close()
    opt.BaseDir = filepath.Dir(path)
    return LoadFromReaderWithOptions(f, opt)
}

// Interpolate 对配置中的占位符进行插值替换
// 参数: 无
// 返回值: 错误；循环引用时返回错误
// 关键步骤：支持 ${key} 与 ${section.key}，限制最大递归深度避免死循环
func (c *Config) Interpolate() error {
    if c == nil { return nil }
    const maxDepth = 10
    visiting := make(map[string]bool)

    for s, kv := range c.data {
        for k := range kv {
            _, err := c.resolveValue(s, k, maxDepth, visiting)
            if err != nil { return err }
        }
    }
    return nil
}

// resolveValue 解析并返回键的最终字符串（含插值）
// 参数 section: 区段名称
// 参数 key: 键名称
// 参数 depth: 剩余递归深度
// 参数 visiting: 当前递归访问栈，用于检测循环
// 返回值: 解析后的字符串与错误
// 关键步骤：扫描并替换 ${...} 占位符
func (c *Config) resolveValue(section, key string, depth int, visiting map[string]bool) (string, error) {
    path := section + "|" + key
    if visiting[path] { return "", errors.New("interpolation cycle: " + path) }
    visiting[path] = true
    defer delete(visiting, path)
    v := c.GetString(section, key, "")
    if v == "" || depth <= 0 { return v, nil }

    // 关键步骤：逐字符扫描占位符
    out := strings.Builder{}
    for i := 0; i < len(v); {
        if v[i] == '$' && i+1 < len(v) && v[i+1] == '{' {
            j := i + 2
            for j < len(v) && v[j] != '}' { j++ }
            if j >= len(v) { // 未闭合，原样输出
                out.WriteByte(v[i])
                i++
                continue
            }
            token := v[i+2 : j]
            refSec := section
            refKey := token
            if p := strings.IndexByte(token, '.'); p != -1 {
                refSec = token[:p]
                refKey = token[p+1:]
            }
            rv, err := c.resolveValue(refSec, refKey, depth-1, visiting)
            if err != nil { return "", err }
            out.WriteString(rv)
            i = j + 1
        } else {
            out.WriteByte(v[i])
            i++
        }
    }
    // 关键步骤：写回解析结果
    c.Set(section, key, out.String())
    return out.String(), nil
}

// 注意：字符串数组的拆分/连接属于通用字符串操作，已迁移到 strutil 包。
// 若需要将配置中的值读为切片，可读取为字符串后使用 strutil 包的工具函数处理。

// stripInlineComment 内联注释剥离（不处理引号内的分号与井号）
// 参数 s: 原始值字符串
// 返回值: 去除内联注释后的值字符串
// 关键步骤：简易状态机识别双引号范围
func stripInlineComment(s string) string {
    inQuote := false
    for i := 0; i < len(s); i++ {
        ch := s[i]
        if ch == '"' { inQuote = !inQuote; continue }
        if !inQuote && (ch == ';' || ch == '#') {
            return strings.TrimSpace(s[:i])
        }
    }
    return strings.TrimSpace(s)
}