package convert

import (
    "reflect"
)

// 本文件提供数组/切片与 map 的互相转换工具函数

// ToMapFromSlice 将数组或切片转换为以索引为键的 map
// 参数 v: 源值（必须为数组或切片类型）
// 返回值: 以索引为键、元素为值的 map 和错误；当类型不匹配时返回错误
// 关键步骤：使用反射遍历数组/切片，按索引写入 map
func ToMapFromSlice(v interface{}) (map[int]interface{}, error) {
    if v == nil {
        return nil, newConvertError("nil", "map[int]interface{}", v, "输入为空")
    }
    rv := reflect.ValueOf(v)
    k := rv.Kind()
    if k != reflect.Slice && k != reflect.Array {
        return nil, newConvertError(rv.Type().String(), "map[int]interface{}", v, "仅支持数组或切片")
    }
    n := rv.Len()
    out := make(map[int]interface{}, n)
    for i := 0; i < n; i++ {
        // 关键步骤：取出元素并按索引写入映射
        out[i] = rv.Index(i).Interface()
    }
    return out, nil
}

// ToSliceFromMapValues 将 map 的所有值转换为切片（顺序不保证）
// 参数 v: 源值（必须为 map 类型，键值类型任意）
// 返回值: 值组成的切片与错误；当类型不匹配时返回错误
// 关键步骤：遍历 MapKeys 并收集对应的值
func ToSliceFromMapValues(v interface{}) ([]interface{}, error) {
    if v == nil {
        return nil, newConvertError("nil", "[]interface{}", v, "输入为空")
    }
    rv := reflect.ValueOf(v)
    if rv.Kind() != reflect.Map {
        return nil, newConvertError(rv.Type().String(), "[]interface{}", v, "仅支持 map 类型")
    }
    keys := rv.MapKeys()
    out := make([]interface{}, 0, len(keys))
    for _, k := range keys {
        // 关键步骤：按键读取值并加入切片
        val := rv.MapIndex(k)
        out = append(out, val.Interface())
    }
    return out, nil
}

// ToSliceFromMapKeys 将 map 的所有键转换为切片（顺序不保证）
// 参数 v: 源值（必须为 map 类型，键值类型任意）
// 返回值: 键组成的切片与错误；当类型不匹配时返回错误
// 关键步骤：直接收集 MapKeys 并输出为切片
func ToSliceFromMapKeys(v interface{}) ([]interface{}, error) {
    if v == nil {
        return nil, newConvertError("nil", "[]interface{}", v, "输入为空")
    }
    rv := reflect.ValueOf(v)
    if rv.Kind() != reflect.Map {
        return nil, newConvertError(rv.Type().String(), "[]interface{}", v, "仅支持 map 类型")
    }
    keys := rv.MapKeys()
    out := make([]interface{}, 0, len(keys))
    for _, k := range keys {
        // 关键步骤：将键加入输出切片
        out = append(out, k.Interface())
    }
    return out, nil
}