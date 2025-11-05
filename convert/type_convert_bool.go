package convert

import (
    "fmt"
    "strings"
)

// 本文件提供布尔类型转换工具

// ToBool 将常见类型转换为bool
// 参数 v: 源值（支持string/bool/int/uint/float）
// 返回值: bool值与错误（转换失败时返回错误）
// 关键步骤：字符串按常见真值/假值词匹配；数值非零为true
func ToBool(v interface{}) (bool, error) {
    if v == nil {
        return false, newConvertError("nil", "bool", v, "输入为nil")
    }
    switch x := v.(type) {
    case bool:
        return x, nil
    case int:
        return x != 0, nil
    case int8:
        return x != 0, nil
    case int16:
        return x != 0, nil
    case int32:
        return x != 0, nil
    case int64:
        return x != 0, nil
    case uint:
        return x != 0, nil
    case uint8:
        return x != 0, nil
    case uint16:
        return x != 0, nil
    case uint32:
        return x != 0, nil
    case uint64:
        return x != 0, nil
    case float32:
        return x != 0, nil
    case float64:
        return x != 0, nil
    case string:
        s := strings.TrimSpace(strings.ToLower(x))
        if s == "" {
            return false, newConvertError("string", "bool", x, "空字符串")
        }
        // 关键步骤：支持多语言常见词
        switch s {
        case "1", "t", "true", "y", "yes", "on", "ok", "是", "真", "对", "开":
            return true, nil
        case "0", "f", "false", "n", "no", "off", "否", "假", "错", "关":
            return false, nil
        default:
            return false, newConvertError("string", "bool", x, "未知布尔字符串")
        }
    default:
        return false, newConvertError(fmt.Sprintf("%T", v), "bool", v, "不支持的类型")
    }
}