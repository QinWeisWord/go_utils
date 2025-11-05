package convert

import (
    "fmt"
    "strconv"
    "strings"
)

// 本文件提供常见进制（2/8/16）格式化与解析工具，并支持字符串在不同进制间转换

// ToBinString 将常见数值类型转换为二进制字符串（无前缀）
// 参数 v: 源值（支持int/uint/float/bool）
// 返回值: 二进制字符串与错误（若类型不支持则返回错误）
// 关键步骤：对无符号类型使用FormatUint，其余使用FormatInt（负数保留前缀“-”）
func ToBinString(v interface{}) (string, error) {
    switch x := v.(type) {
    case int:
        return strconv.FormatInt(int64(x), 2), nil
    case int8:
        return strconv.FormatInt(int64(x), 2), nil
    case int16:
        return strconv.FormatInt(int64(x), 2), nil
    case int32:
        return strconv.FormatInt(int64(x), 2), nil
    case int64:
        return strconv.FormatInt(x, 2), nil
    case uint:
        return strconv.FormatUint(uint64(x), 2), nil
    case uint8:
        return strconv.FormatUint(uint64(x), 2), nil
    case uint16:
        return strconv.FormatUint(uint64(x), 2), nil
    case uint32:
        return strconv.FormatUint(uint64(x), 2), nil
    case uint64:
        return strconv.FormatUint(x, 2), nil
    case float32:
        return strconv.FormatInt(int64(x), 2), nil
    case float64:
        return strconv.FormatInt(int64(x), 2), nil
    case bool:
        if x { return "1", nil }
        return "0", nil
    default:
        return "", newConvertError(fmt.Sprintf("%T", v), "bin string", v, "不支持的类型")
    }
}

// ToBinStringWithPrefix 将数值转换为二进制字符串并带前缀0b
// 参数 v: 源值（支持int/uint/float/bool）
// 返回值: 带前缀的二进制字符串与错误
// 关键步骤：调用ToBinString并追加"0b"前缀（负数保留在最前）
func ToBinStringWithPrefix(v interface{}) (string, error) {
    s, err := ToBinString(v)
    if err != nil { return "", err }
    if strings.HasPrefix(s, "-") { return "-0b" + s[1:], nil }
    return "0b" + s, nil
}

// ToOctString 将常见数值类型转换为八进制字符串（无前缀）
// 参数 v: 源值（支持int/uint/float/bool）
// 返回值: 八进制字符串与错误
// 关键步骤：与二进制一致，基数为8
func ToOctString(v interface{}) (string, error) {
    switch x := v.(type) {
    case int:
        return strconv.FormatInt(int64(x), 8), nil
    case int8:
        return strconv.FormatInt(int64(x), 8), nil
    case int16:
        return strconv.FormatInt(int64(x), 8), nil
    case int32:
        return strconv.FormatInt(int64(x), 8), nil
    case int64:
        return strconv.FormatInt(x, 8), nil
    case uint:
        return strconv.FormatUint(uint64(x), 8), nil
    case uint8:
        return strconv.FormatUint(uint64(x), 8), nil
    case uint16:
        return strconv.FormatUint(uint64(x), 8), nil
    case uint32:
        return strconv.FormatUint(uint64(x), 8), nil
    case uint64:
        return strconv.FormatUint(x, 8), nil
    case float32:
        return strconv.FormatInt(int64(x), 8), nil
    case float64:
        return strconv.FormatInt(int64(x), 8), nil
    case bool:
        if x { return "1", nil }
        return "0", nil
    default:
        return "", newConvertError(fmt.Sprintf("%T", v), "oct string", v, "不支持的类型")
    }
}

// ToOctStringWithPrefix 将数值转换为八进制字符串并带前缀0o
// 参数 v: 源值（支持int/uint/float/bool）
// 返回值: 带前缀的八进制字符串与错误
// 关键步骤：调用ToOctString并追加"0o"前缀（负数保留在最前）
func ToOctStringWithPrefix(v interface{}) (string, error) {
    s, err := ToOctString(v)
    if err != nil { return "", err }
    if strings.HasPrefix(s, "-") { return "-0o" + s[1:], nil }
    return "0o" + s, nil
}

// ToHexString 将常见数值类型转换为十六进制字符串（无前缀）
// 参数 v: 源值（支持int/uint/float/bool）
// 参数 uppercase: 是否输出大写（A-F）
// 返回值: 十六进制字符串与错误
// 关键步骤：使用FormatInt/FormatUint并按需转大写
func ToHexString(v interface{}, uppercase bool) (string, error) {
    var s string
    switch x := v.(type) {
    case int:
        s = strconv.FormatInt(int64(x), 16)
    case int8:
        s = strconv.FormatInt(int64(x), 16)
    case int16:
        s = strconv.FormatInt(int64(x), 16)
    case int32:
        s = strconv.FormatInt(int64(x), 16)
    case int64:
        s = strconv.FormatInt(x, 16)
    case uint:
        s = strconv.FormatUint(uint64(x), 16)
    case uint8:
        s = strconv.FormatUint(uint64(x), 16)
    case uint16:
        s = strconv.FormatUint(uint64(x), 16)
    case uint32:
        s = strconv.FormatUint(uint64(x), 16)
    case uint64:
        s = strconv.FormatUint(x, 16)
    case float32:
        s = strconv.FormatInt(int64(x), 16)
    case float64:
        s = strconv.FormatInt(int64(x), 16)
    case bool:
        if x { s = "1" } else { s = "0" }
    default:
        return "", newConvertError(fmt.Sprintf("%T", v), "hex string", v, "不支持的类型")
    }
    if uppercase { s = strings.ToUpper(s) }
    return s, nil
}

// ToHexStringWithPrefix 将数值转换为十六进制字符串并带前缀0x
// 参数 v: 源值（支持int/uint/float/bool）
// 参数 uppercase: 是否输出大写（A-F）
// 返回值: 带前缀的十六进制字符串与错误
// 关键步骤：调用ToHexString并追加"0x"前缀（负数保留在最前）
func ToHexStringWithPrefix(v interface{}, uppercase bool) (string, error) {
    s, err := ToHexString(v, uppercase)
    if err != nil { return "", err }
    // 前缀统一使用0x（不区分大小写），仅控制数字大小写
    if strings.HasPrefix(s, "-") { return "-0x" + s[1:], nil }
    return "0x" + s, nil
}

// ParseIntFromBase 解析指定进制的字符串为int64
// 参数 s: 输入字符串（可含前缀"-"、允许空白）
// 参数 base: 进制，范围[2,36]；为0时启用自动检测（支持0x/0o/0b前缀）
// 返回值: 解析得到的int64与错误
// 关键步骤：Trim空白并调用strconv.ParseInt
func ParseIntFromBase(s string, base int) (int64, error) {
    s = strings.TrimSpace(s)
    if base != 0 && (base < 2 || base > 36) {
        return 0, newConvertError("string", "int64", s, "非法进制范围")
    }
    v, err := strconv.ParseInt(s, base, 64)
    if err != nil {
        return 0, newConvertError("string", "int64", s, err.Error())
    }
    return v, nil
}

// ParseUintFromBase 解析指定进制的字符串为uint64
// 参数 s: 输入字符串（不支持负号）
// 参数 base: 进制，范围[2,36]；为0时启用自动检测（支持0x/0o/0b前缀）
// 返回值: 解析得到的uint64与错误
// 关键步骤：Trim空白并调用strconv.ParseUint
func ParseUintFromBase(s string, base int) (uint64, error) {
    s = strings.TrimSpace(s)
    if base != 0 && (base < 2 || base > 36) {
        return 0, newConvertError("string", "uint64", s, "非法进制范围")
    }
    if strings.HasPrefix(s, "-") {
        return 0, newConvertError("string", "uint64", s, "不支持负数")
    }
    v, err := strconv.ParseUint(s, base, 64)
    if err != nil {
        return 0, newConvertError("string", "uint64", s, err.Error())
    }
    return v, nil
}

// ConvertBaseString 将数字字符串从一种进制转换为另一种进制
// 参数 s: 输入字符串（可带符号；支持0x/0o/0b前缀当fromBase为0时）
// 参数 fromBase: 源进制（2~36；为0表示自动检测）
// 参数 toBase: 目标进制（2~36）
// 参数 uppercase: 是否输出大写（十六进制A-F等）
// 参数 withPrefix: 是否附加进制前缀（2→0b，8→0o，16→0x，其它无前缀）
// 返回值: 转换后的字符串与错误
// 关键步骤：先按有符号解析，再用FormatInt按目标进制格式化并追加前缀
func ConvertBaseString(s string, fromBase, toBase int, uppercase, withPrefix bool) (string, error) {
    if toBase < 2 || toBase > 36 {
        return "", newConvertError("string", fmt.Sprintf("base%d string", toBase), s, "非法目标进制范围")
    }
    s = strings.TrimSpace(s)
    v, err := strconv.ParseInt(s, fromBase, 64)
    if err != nil {
        return "", newConvertError("string", fmt.Sprintf("base%d string", toBase), s, err.Error())
    }
    out := strconv.FormatInt(v, toBase)
    if uppercase { out = strings.ToUpper(out) }
    if withPrefix {
        prefix := ""
        switch toBase {
        case 2: prefix = "0b"
        case 8: prefix = "0o"
        case 16: prefix = "0x"
        }
        if len(prefix) > 0 {
            if strings.HasPrefix(out, "-") { out = "-" + prefix + out[1:] } else { out = prefix + out }
        }
    }
    return out, nil
}