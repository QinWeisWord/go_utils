package convert

import (
    "fmt"
    "strconv"
    "strings"
)

// 本文件提供数值类型的转换工具

// ToInt 将常见类型转换为int
// 参数 v: 源值（支持string/int/uint/float/bool）
// 返回值: int值与错误（转换失败时返回错误）
// 关键步骤：字符串使用strconv.ParseInt，浮点数截断取整，布尔true为1/false为0
func ToInt(v interface{}) (int, error) {
    if v == nil {
        return 0, newConvertError("nil", "int", v, "输入为nil")
    }
    switch x := v.(type) {
    case int:
        return x, nil
    case int8:
        return int(x), nil
    case int16:
        return int(x), nil
    case int32:
        return int(x), nil
    case int64:
        // 关键步骤：可能溢出int范围，但Go的转换会截断为平台int宽度
        return int(x), nil
    case uint, uint8, uint16, uint32, uint64:
        // 关键步骤：统一转为uint64再转int
        var u uint64
        switch y := v.(type) {
        case uint: u = uint64(y)
        case uint8: u = uint64(y)
        case uint16: u = uint64(y)
        case uint32: u = uint64(y)
        case uint64: u = y
        }
        return int(u), nil
    case float32:
        return int(x), nil
    case float64:
        return int(x), nil
    case bool:
        if x { return 1, nil }
        return 0, nil
    case string:
        s := strings.TrimSpace(x)
        if len(s) == 0 {
            return 0, newConvertError("string", "int", x, "空字符串")
        }
        i64, err := strconv.ParseInt(s, 10, 64)
        if err != nil {
            return 0, newConvertError("string", "int", x, err.Error())
        }
        return int(i64), nil
    default:
        return 0, newConvertError(fmt.Sprintf("%T", v), "int", v, "不支持的类型")
    }
}

// ToInt64 将常见类型转换为int64
// 参数 v: 源值（支持string/int/uint/float/bool）
// 返回值: int64值与错误（转换失败时返回错误）
// 关键步骤：字符串使用strconv.ParseInt，浮点数截断取整，布尔true为1/false为0
func ToInt64(v interface{}) (int64, error) {
    if v == nil {
        return 0, newConvertError("nil", "int64", v, "输入为nil")
    }
    switch x := v.(type) {
    case int:
        return int64(x), nil
    case int8:
        return int64(x), nil
    case int16:
        return int64(x), nil
    case int32:
        return int64(x), nil
    case int64:
        return x, nil
    case uint:
        return int64(x), nil
    case uint8:
        return int64(x), nil
    case uint16:
        return int64(x), nil
    case uint32:
        return int64(x), nil
    case uint64:
        return int64(x), nil
    case float32:
        return int64(x), nil
    case float64:
        return int64(x), nil
    case bool:
        if x { return 1, nil }
        return 0, nil
    case string:
        s := strings.TrimSpace(x)
        if len(s) == 0 {
            return 0, newConvertError("string", "int64", x, "空字符串")
        }
        i64, err := strconv.ParseInt(s, 10, 64)
        if err != nil {
            return 0, newConvertError("string", "int64", x, err.Error())
        }
        return i64, nil
    default:
        return 0, newConvertError(fmt.Sprintf("%T", v), "int64", v, "不支持的类型")
    }
}

// ToUint64 将常见类型转换为uint64
// 参数 v: 源值（支持string/int/uint/float/bool）
// 返回值: uint64值与错误（转换失败时返回错误）
// 关键步骤：字符串使用strconv.ParseUint，负数与负浮点视为错误
func ToUint64(v interface{}) (uint64, error) {
    if v == nil {
        return 0, newConvertError("nil", "uint64", v, "输入为nil")
    }
    switch x := v.(type) {
    case uint:
        return uint64(x), nil
    case uint8:
        return uint64(x), nil
    case uint16:
        return uint64(x), nil
    case uint32:
        return uint64(x), nil
    case uint64:
        return x, nil
    case int, int8, int16, int32, int64:
        var i int64
        switch y := v.(type) {
        case int: i = int64(y)
        case int8: i = int64(y)
        case int16: i = int64(y)
        case int32: i = int64(y)
        case int64: i = y
        }
        if i < 0 {
            return 0, newConvertError("int", "uint64", v, "负数不能转换为无符号")
        }
        return uint64(i), nil
    case float32:
        if x < 0 { return 0, newConvertError("float32", "uint64", x, "负数不能转换为无符号") }
        return uint64(x), nil
    case float64:
        if x < 0 { return 0, newConvertError("float64", "uint64", x, "负数不能转换为无符号") }
        return uint64(x), nil
    case bool:
        if x { return 1, nil }
        return 0, nil
    case string:
        s := strings.TrimSpace(x)
        if len(s) == 0 {
            return 0, newConvertError("string", "uint64", x, "空字符串")
        }
        u64, err := strconv.ParseUint(s, 10, 64)
        if err != nil {
            return 0, newConvertError("string", "uint64", x, err.Error())
        }
        return u64, nil
    default:
        return 0, newConvertError(fmt.Sprintf("%T", v), "uint64", v, "不支持的类型")
    }
}

// ToFloat64 将常见类型转换为float64
// 参数 v: 源值（支持string/int/uint/float/bool）
// 返回值: float64值与错误（转换失败时返回错误）
// 关键步骤：字符串使用strconv.ParseFloat，布尔true为1/false为0
func ToFloat64(v interface{}) (float64, error) {
    if v == nil {
        return 0, newConvertError("nil", "float64", v, "输入为nil")
    }
    switch x := v.(type) {
    case float32:
        return float64(x), nil
    case float64:
        return x, nil
    case int:
        return float64(x), nil
    case int8:
        return float64(x), nil
    case int16:
        return float64(x), nil
    case int32:
        return float64(x), nil
    case int64:
        return float64(x), nil
    case uint:
        return float64(x), nil
    case uint8:
        return float64(x), nil
    case uint16:
        return float64(x), nil
    case uint32:
        return float64(x), nil
    case uint64:
        return float64(x), nil
    case bool:
        if x { return 1, nil }
        return 0, nil
    case string:
        s := strings.TrimSpace(x)
        if len(s) == 0 {
            return 0, newConvertError("string", "float64", x, "空字符串")
        }
        f64, err := strconv.ParseFloat(s, 64)
        if err != nil {
            return 0, newConvertError("string", "float64", x, err.Error())
        }
        return f64, nil
    default:
        return 0, newConvertError(fmt.Sprintf("%T", v), "float64", v, "不支持的类型")
    }
}

// ToFloat32 将常见类型转换为float32
// 参数 v: 源值（支持string/int/uint/float/bool）
// 返回值: float32值与错误（转换失败时返回错误）
// 关键步骤：先调用ToFloat64统一解析，再转换为float32（可能存在精度损失）
func ToFloat32(v interface{}) (float32, error) {
    f64, err := ToFloat64(v)
    if err != nil {
        return 0, err
    }
    return float32(f64), nil
}

// ToInt32 将常见类型转换为int32
// 参数 v: 源值（支持string/int/uint/float/bool）
// 返回值: int32值与错误（转换失败或溢出时返回错误）
// 关键步骤：先转为int64并检查范围[-2147483648, 2147483647]
func ToInt32(v interface{}) (int32, error) {
    i64, err := ToInt64(v)
    if err != nil {
        return 0, err
    }
    if i64 < -2147483648 || i64 > 2147483647 {
        return 0, newConvertError(fmt.Sprintf("%T", v), "int32", v, "超出int32范围")
    }
    return int32(i64), nil
}

// ToUint32 将常见类型转换为uint32
// 参数 v: 源值（支持string/int/uint/float/bool）
// 返回值: uint32值与错误（转换失败或溢出时返回错误）
// 关键步骤：先转为uint64并检查范围[0, 4294967295]
func ToUint32(v interface{}) (uint32, error) {
    u64, err := ToUint64(v)
    if err != nil {
        return 0, err
    }
    if u64 > 4294967295 {
        return 0, newConvertError(fmt.Sprintf("%T", v), "uint32", v, "超出uint32范围")
    }
    return uint32(u64), nil
}

// ToInt16 将常见类型转换为int16
// 参数 v: 源值（支持string/int/uint/float/bool）
// 返回值: int16值与错误（转换失败或溢出时返回错误）
// 关键步骤：先转为int64并检查范围[-32768, 32767]
func ToInt16(v interface{}) (int16, error) {
    i64, err := ToInt64(v)
    if err != nil {
        return 0, err
    }
    if i64 < -32768 || i64 > 32767 {
        return 0, newConvertError(fmt.Sprintf("%T", v), "int16", v, "超出int16范围")
    }
    return int16(i64), nil
}

// ToUint16 将常见类型转换为uint16
// 参数 v: 源值（支持string/int/uint/float/bool）
// 返回值: uint16值与错误（转换失败或溢出时返回错误）
// 关键步骤：先转为uint64并检查范围[0, 65535]
func ToUint16(v interface{}) (uint16, error) {
    u64, err := ToUint64(v)
    if err != nil {
        return 0, err
    }
    if u64 > 65535 {
        return 0, newConvertError(fmt.Sprintf("%T", v), "uint16", v, "超出uint16范围")
    }
    return uint16(u64), nil
}