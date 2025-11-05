package convert

import (
    "fmt"
    "strconv"
)

// 本文件提供字符串相关的转换工具

// ToString 将任意常见类型转换为字符串
// 参数 v: 源值（支持int/uint/float/bool/string以及实现fmt.Stringer的类型）
// 返回值: 转换得到的字符串与错误（若不支持转换则返回错误）
// 关键步骤：根据具体类型分支处理，数字使用strconv格式化，布尔转换为"true"/"false"
func ToString(v interface{}) (string, error) {
    if v == nil {
        return "", newConvertError("nil", "string", v, "输入为nil")
    }
    switch x := v.(type) {
    case string:
        return x, nil
    case bool:
        if x { return "true", nil }
        return "false", nil
    case int:
        return strconv.Itoa(x), nil
    case int8:
        return strconv.FormatInt(int64(x), 10), nil
    case int16:
        return strconv.FormatInt(int64(x), 10), nil
    case int32:
        return strconv.FormatInt(int64(x), 10), nil
    case int64:
        return strconv.FormatInt(x, 10), nil
    case uint:
        return strconv.FormatUint(uint64(x), 10), nil
    case uint8:
        return strconv.FormatUint(uint64(x), 10), nil
    case uint16:
        return strconv.FormatUint(uint64(x), 10), nil
    case uint32:
        return strconv.FormatUint(uint64(x), 10), nil
    case uint64:
        return strconv.FormatUint(x, 10), nil
    case float32:
        return strconv.FormatFloat(float64(x), 'f', -1, 64), nil
    case float64:
        return strconv.FormatFloat(x, 'f', -1, 64), nil
    default:
        // 关键步骤：尝试fmt.Stringer接口
        type stringer interface{ String() string }
        if s, ok := v.(stringer); ok {
            return s.String(), nil
        }
        return "", newConvertError(fmt.Sprintf("%T", v), "string", v, "不支持的类型")
    }
}