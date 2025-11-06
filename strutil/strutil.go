package strutil

import (
    "strings"
    "unicode/utf8"
)

// IsEmpty 判断字符串是否为空（包含去除空白后为空）
// 参数 s: 待判断的字符串
// 返回值: 布尔值，true表示为空或仅包含空白，false表示非空
// 关键步骤：先Trim空白再判断长度
func IsEmpty(s string) bool {
    // 关键步骤：去除首尾空白
    trimmed := strings.TrimSpace(s)
    // 关键步骤：判断长度
    return len(trimmed) == 0
}

// Trim 去除字符串首尾空白字符
// 参数 s: 原始字符串
// 返回值: 去除空白后的新字符串
// 关键步骤：调用标准库strings.TrimSpace
func Trim(s string) string {
    return strings.TrimSpace(s)
}

// ToUpper 将字符串转换为大写
// 参数 s: 原始字符串
// 返回值: 转换为大写后的字符串
// 关键步骤：调用strings.ToUpper
func ToUpper(s string) string {
    return strings.ToUpper(s)
}

// ToLower 将字符串转换为小写
// 参数 s: 原始字符串
// 返回值: 转换为小写后的字符串
// 关键步骤：调用strings.ToLower
func ToLower(s string) string {
    return strings.ToLower(s)
}

// ContainsSubstr 判断子串是否存在于字符串中
// 参数 s: 原始字符串
// 参数 sub: 需要查找的子串
// 返回值: 布尔值，true表示包含，false表示不包含
// 关键步骤：调用strings.Contains
func ContainsSubstr(s, sub string) bool {
    return strings.Contains(s, sub)
}

// ReplaceAll 替换字符串中的所有匹配子串
// 参数 s: 原始字符串
// 参数 old: 需要被替换的子串
// 参数 new: 替换为的新子串
// 返回值: 完成替换后的新字符串
// 关键步骤：调用strings.ReplaceAll
func ReplaceAll(s, old, new string) string {
    return strings.ReplaceAll(s, old, new)
}

// Split 按分隔符拆分字符串为切片
// 参数 s: 原始字符串
// 参数 sep: 分隔符字符串
// 返回值: 拆分后的字符串切片
// 关键步骤：调用strings.Split
func Split(s, sep string) []string {
    return strings.Split(s, sep)
}

// Join 将字符串切片通过分隔符连接成字符串
// 参数 parts: 字符串切片
// 参数 sep: 分隔符字符串
// 返回值: 连接后的新字符串
// 关键步骤：调用strings.Join
func Join(parts []string, sep string) string {
    return strings.Join(parts, sep)
}

// SplitTrimNonEmpty 按分隔符拆分并去除空白与空项
// 参数 s: 原始字符串
// 参数 sep: 分隔符字符串
// 返回值: 去除空白后的非空字符串切片
// 关键步骤：先Split，再逐项TrimSpace，过滤空字符串
func SplitTrimNonEmpty(s, sep string) []string {
    // 关键步骤：基础拆分
    parts := strings.Split(s, sep)
    // 关键步骤：Trim并过滤空项
    out := make([]string, 0, len(parts))
    for _, p := range parts {
        p = strings.TrimSpace(p)
        if len(p) > 0 {
            out = append(out, p)
        }
    }
    return out
}

// Substring 按字符索引截取子串（支持UTF-8）
// 参数 s: 原始字符串
// 参数 start: 起始字符索引（基于Unicode字符）
// 参数 length: 截取的字符数量；若为负数或超出范围，则尽量安全返回
// 返回值: 截取得到的子串
// 关键步骤：将字符串按UTF-8遍历再截取，避免字节切割导致乱码
func Substring(s string, start, length int) string {
    if start < 0 {
        start = 0
    }
    if length == 0 {
        return ""
    }

    // 关键步骤：遍历Unicode字符，记录每个字符的起止位置
    var runePositions []int
    for i := range s {
        runePositions = append(runePositions, i)
    }
    runePositions = append(runePositions, len(s))

    totalRunes := utf8.RuneCountInString(s)
    if start >= totalRunes {
        return ""
    }

    end := totalRunes
    if length > 0 && start+length < end {
        end = start + length
    }

    // 关键步骤：计算字节边界并安全切割
    byteStart := runePositions[start]
    byteEnd := runePositions[end]
    return s[byteStart:byteEnd]
}

// Reverse 反转字符串（按Unicode字符）
// 参数 s: 原始字符串
// 返回值: 字符反转后的新字符串
// 关键步骤：将字符串转换为rune切片后反转
func Reverse(s string) string {
    // 关键步骤：转换为rune切片
    r := []rune(s)
    // 关键步骤：双指针反转
    for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r)
}

// PadLeft 在字符串左侧补齐到指定长度
// 参数 s: 原始字符串
// 参数 pad: 用于补齐的字符
// 参数 totalLen: 目标总长度（按Unicode字符计）
// 返回值: 补齐后的新字符串
// 关键步骤：计算当前字符长度与差值，并重复补齐字符
func PadLeft(s string, pad string, totalLen int) string {
    curLen := utf8.RuneCountInString(s)
    if curLen >= totalLen {
        return s
    }
    need := totalLen - curLen
    return strings.Repeat(pad, need) + s
}

// PadRight 在字符串右侧补齐到指定长度
// 参数 s: 原始字符串
// 参数 pad: 用于补齐的字符
// 参数 totalLen: 目标总长度（按Unicode字符计）
// 返回值: 补齐后的新字符串
// 关键步骤：计算当前字符长度与差值，并重复补齐字符
func PadRight(s string, pad string, totalLen int) string {
    curLen := utf8.RuneCountInString(s)
    if curLen >= totalLen {
        return s
    }
    need := totalLen - curLen
    return s + strings.Repeat(pad, need)
}