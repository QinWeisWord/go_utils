package convert

import (
    "fmt"
    "reflect"
    "strings"
)

// 本文件提供结构体与 map 的互相转换工具，支持字段标签（如 json）

// ToMapFromStruct 将结构体转换为 map[string]interface{}
// 参数 v: 源值（结构体或结构体指针）
// 参数 tag: 标签名（如"json"；为空表示直接使用字段名）
// 返回值: 字段名为键、字段值为值的 map 与错误；当类型不匹配或输入为空时返回错误
// 关键步骤：使用反射遍历导出字段，依据标签确定键名；指针字段与嵌套结构体做递归处理
func ToMapFromStruct(v interface{}, tag string) (map[string]interface{}, error) {
    if v == nil {
        return nil, newConvertError("nil", "map[string]interface{}", v, "输入为空")
    }
    rv := reflect.ValueOf(v)
    if rv.Kind() == reflect.Ptr {
        if rv.IsNil() {
            return nil, newConvertError(rv.Type().String(), "map[string]interface{}", v, "指针为空")
        }
        rv = rv.Elem()
    }
    if rv.Kind() != reflect.Struct {
        return nil, newConvertError(rv.Type().String(), "map[string]interface{}", v, "仅支持结构体或结构体指针")
    }
    rt := rv.Type()
    out := make(map[string]interface{}, rt.NumField())
    for i := 0; i < rt.NumField(); i++ {
        sf := rt.Field(i)
        // 跳过未导出字段
        if sf.PkgPath != "" { // 非空表示未导出
            continue
        }
        key, skip := getFieldKey(sf, tag)
        if skip { continue }
        fv := rv.Field(i)
        // 关键步骤：指针字段与嵌套结构体递归处理
        if fv.Kind() == reflect.Ptr {
            if fv.IsNil() {
                out[key] = nil
                continue
            }
            elem := fv.Elem()
            if elem.Kind() == reflect.Struct {
                nested, err := ToMapFromStruct(elem.Interface(), tag)
                if err != nil { return nil, err }
                out[key] = nested
            } else {
                out[key] = elem.Interface()
            }
            continue
        }
        if fv.Kind() == reflect.Struct {
            nested, err := ToMapFromStruct(fv.Interface(), tag)
            if err != nil { return nil, err }
            out[key] = nested
            continue
        }
        out[key] = fv.Interface()
    }
    return out, nil
}

// FillStructFromMap 将 map 的值填充到结构体（支持标签名匹配）
// 参数 m: 源 map（键为字符串，值为任意类型）
// 参数 out: 目标结构体指针（必须为非nil的结构体指针）
// 参数 tag: 标签名（如"json"；为空表示按字段名匹配）
// 返回值: 错误（当类型不匹配、指针为空或赋值失败时返回错误）
// 关键步骤：反射遍历字段，按键查找值，并进行必要的类型转换与递归结构体填充
func FillStructFromMap(m map[string]interface{}, out interface{}, tag string) error {
    if m == nil {
        return newConvertError("map[string]interface{}", "struct", m, "输入map为空")
    }
    if out == nil {
        return newConvertError("nil", "struct", out, "输出目标为空")
    }
    rv := reflect.ValueOf(out)
    if rv.Kind() != reflect.Ptr || rv.IsNil() {
        return newConvertError(fmt.Sprintf("%T", out), "struct", out, "必须为非nil的结构体指针")
    }
    rv = rv.Elem()
    if rv.Kind() != reflect.Struct {
        return newConvertError(rv.Type().String(), "struct", out, "仅支持结构体指针")
    }
    rt := rv.Type()
    for i := 0; i < rt.NumField(); i++ {
        sf := rt.Field(i)
        if sf.PkgPath != "" { // 未导出字段不可设置
            continue
        }
        key, skip := getFieldKey(sf, tag)
        if skip { continue }
        val, ok := m[key]
        if !ok { continue }
        fv := rv.Field(i)
        if err := setFieldValue(fv, val, tag); err != nil {
            return err
        }
    }
    return nil
}

// getFieldKey 根据字段与标签名确定键名（私有辅助函数，置于公有方法之后）
// 参数 sf: 结构体字段信息
// 参数 tag: 标签名（例如"json"；为空表示不用标签）
// 返回值: 键名与是否跳过该字段（当标签为"-"时跳过）
// 关键步骤：解析标签第一个片段（逗号前），为"-"则跳过，否则使用字段名或标签名
func getFieldKey(sf reflect.StructField, tag string) (string, bool) {
    if tag == "" {
        return sf.Name, false
    }
    tv := sf.Tag.Get(tag)
    if tv == "-" {
        return "", true
    }
    if tv == "" {
        return sf.Name, false
    }
    // 关键步骤：仅取逗号前的第一个片段
    name := tv
    if i := strings.IndexByte(tv, ','); i >= 0 {
        name = tv[:i]
    }
    if name == "" { // 标签为",omitempty"等情况，退回字段名
        return sf.Name, false
    }
    return name, false
}

// setFieldValue 设置字段值，支持指针、结构体递归以及常见基础类型转换
// 参数 fv: 目标字段的反射值
// 参数 val: 源值
// 参数 tag: 标签名（用于嵌套结构体递归时）
// 返回值: 错误（当类型不兼容时返回错误）
// 关键步骤：处理指针分配、结构体递归及基础类型（字符串/布尔/数值）转换
func setFieldValue(fv reflect.Value, val interface{}, tag string) error {
    if !fv.CanSet() {
        return newConvertError(fv.Type().String(), "field", val, "字段不可设置")
    }
    if val == nil {
        // 对指针字段赋nil，其它保持零值
        if fv.Kind() == reflect.Ptr {
            fv.Set(reflect.Zero(fv.Type()))
        }
        return nil
    }
    // 指针字段处理
    if fv.Kind() == reflect.Ptr {
        elemType := fv.Type().Elem()
        // 结构体指针：递归填充
        if elemType.Kind() == reflect.Struct {
            mv, ok := val.(map[string]interface{})
            if !ok {
                return newConvertError(fmt.Sprintf("%T", val), elemType.String(), val, "期望map用于结构体指针填充")
            }
            ptr := reflect.New(elemType)
            if err := FillStructFromMap(mv, ptr.Interface(), tag); err != nil { return err }
            fv.Set(ptr)
            return nil
        }
        // 其它基础类型指针：创建新值并设置
        ptr := reflect.New(elemType)
        if err := setFieldValue(ptr.Elem(), val, tag); err != nil { return err }
        fv.Set(ptr)
        return nil
    }

    // 结构体递归（非指针）
    if fv.Kind() == reflect.Struct {
        mv, ok := val.(map[string]interface{})
        if !ok {
            return newConvertError(fmt.Sprintf("%T", val), fv.Type().String(), val, "期望map用于结构体填充")
        }
        // 关键步骤：递归填充子结构体
        return FillStructFromMap(mv, fv.Addr().Interface(), tag)
    }

    // 直接可赋或可转换
    srcVal := reflect.ValueOf(val)
    if srcVal.IsValid() && srcVal.Type().AssignableTo(fv.Type()) {
        fv.Set(srcVal)
        return nil
    }
    if srcVal.IsValid() && srcVal.Type().ConvertibleTo(fv.Type()) {
        fv.Set(srcVal.Convert(fv.Type()))
        return nil
    }

    // 基础类型适配：字符串、布尔、整型、无符号、浮点
    switch fv.Kind() {
    case reflect.String:
        fv.SetString(fmt.Sprintf("%v", val))
        return nil
    case reflect.Bool:
        b, err := ToBool(val)
        if err != nil { return err }
        fv.SetBool(b)
        return nil
    case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
        i64, err := ToInt64(val)
        if err != nil { return err }
        // 关键步骤：窄类型范围校验
        switch fv.Kind() {
        case reflect.Int8:
            if i64 < -128 || i64 > 127 { return newConvertError(fmt.Sprintf("%T", val), "int8", val, "超出int8范围") }
        case reflect.Int16:
            if i64 < -32768 || i64 > 32767 { return newConvertError(fmt.Sprintf("%T", val), "int16", val, "超出int16范围") }
        case reflect.Int32:
            if i64 < -2147483648 || i64 > 2147483647 { return newConvertError(fmt.Sprintf("%T", val), "int32", val, "超出int32范围") }
        }
        fv.SetInt(i64)
        return nil
    case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
        u64, err := ToUint64(val)
        if err != nil { return err }
        // 关键步骤：窄类型范围校验
        switch fv.Kind() {
        case reflect.Uint8:
            if u64 > 255 { return newConvertError(fmt.Sprintf("%T", val), "uint8", val, "超出uint8范围") }
        case reflect.Uint16:
            if u64 > 65535 { return newConvertError(fmt.Sprintf("%T", val), "uint16", val, "超出uint16范围") }
        case reflect.Uint32:
            if u64 > 4294967295 { return newConvertError(fmt.Sprintf("%T", val), "uint32", val, "超出uint32范围") }
        }
        fv.SetUint(u64)
        return nil
    case reflect.Float32, reflect.Float64:
        f64, err := ToFloat64(val)
        if err != nil { return err }
        if fv.Kind() == reflect.Float32 { fv.SetFloat(float64(float32(f64))) } else { fv.SetFloat(f64) }
        return nil
    }

    return newConvertError(fmt.Sprintf("%T", val), fv.Type().String(), val, "不支持的字段类型转换")
}