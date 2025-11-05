package convert

import (
    "fmt"
    "reflect"
    "strings"
)

// 本文件提供数组/切片的去重方法，以及对象数组按指定字段（或标签）去重的方法

// UniqueSlice 对数组或切片进行去重（保持首次出现顺序）
// 参数 v: 源值（必须为数组或切片类型；元素可为任意类型）
// 返回值: 去重后的切片（元素类型与输入一致）与错误；当类型不匹配时返回错误
// 关键步骤：
// 1) 使用反射遍历元素；
// 2) 若元素类型可比较（Comparable），使用map[interface{}]记录已出现键；
// 3) 若不可比较，使用元素的格式化字符串作为键（可能存在碰撞风险，在注释中说明）。
func UniqueSlice(v interface{}) (interface{}, error) {
    if v == nil {
        return nil, newConvertError("nil", "slice(unique)", v, "输入为空")
    }
    rv := reflect.ValueOf(v)
    k := rv.Kind()
    if k != reflect.Slice && k != reflect.Array {
        return nil, newConvertError(rv.Type().String(), "slice(unique)", v, "仅支持数组或切片类型")
    }

    elemType := rv.Type().Elem()
    // 输出统一为切片（即使输入是数组）
    out := reflect.MakeSlice(reflect.SliceOf(elemType), 0, rv.Len())

    // 可比较元素使用接口键，不可比较元素使用字符串键
    useComparable := elemType.Comparable()
    seenI := make(map[interface{}]struct{}, rv.Len())
    seenS := make(map[string]struct{}, rv.Len())

    for i := 0; i < rv.Len(); i++ {
        ev := rv.Index(i)
        if useComparable {
            key := ev.Interface()
            if _, ok := seenI[key]; ok {
                continue
            }
            seenI[key] = struct{}{}
        } else {
            // 关键步骤：不可比较类型采用字符串键（可能存在碰撞风险）
            key := fmt.Sprintf("%#v", ev.Interface())
            if _, ok := seenS[key]; ok {
                continue
            }
            seenS[key] = struct{}{}
        }
        out = reflect.Append(out, ev)
    }
    return out.Interface(), nil
}

// UniqueSliceByField 对对象数组按指定字段（或标签）去重（保持首次出现顺序）
// 参数 v: 源值（必须为数组或切片，元素为结构体/结构体指针，或map[string]T）
// 参数 field: 字段名（结构体字段名或map键名；当提供tag时可匹配标签名）
// 参数 tag: 标签名（如"json"；为空表示仅按字段名匹配）
// 返回值: 去重后的切片（元素类型与输入一致）与错误；当类型不匹配或字段不存在时返回错误
// 关键步骤：
// 1) 对结构体元素：先解析字段索引（支持按标签名匹配），再按字段值去重；
// 2) 对map元素：按给定键取值去重；
// 3) 键值可比较时使用接口键；不可比较时使用字符串键；缺失键统一视为特殊键，仅保留其首次出现的元素。
func UniqueSliceByField(v interface{}, field string, tag string) (interface{}, error) {
    if v == nil {
        return nil, newConvertError("nil", "slice(uniqueByField)", v, "输入为空")
    }
    if field == "" {
        return nil, newConvertError("unknown", "slice(uniqueByField)", v, "字段名不能为空")
    }
    rv := reflect.ValueOf(v)
    k := rv.Kind()
    if k != reflect.Slice && k != reflect.Array {
        return nil, newConvertError(rv.Type().String(), "slice(uniqueByField)", v, "仅支持数组或切片类型")
    }

    elemType := rv.Type().Elem()
    out := reflect.MakeSlice(reflect.SliceOf(elemType), 0, rv.Len())

    // 场景1：结构体或结构体指针
    if elemType.Kind() == reflect.Struct || (elemType.Kind() == reflect.Ptr && elemType.Elem().Kind() == reflect.Struct) {
        isPtr := elemType.Kind() == reflect.Ptr
        structType := elemType
        if isPtr { structType = elemType.Elem() }

        idx, ok := fieldIndexByNameOrTag(structType, field, tag)
        if !ok {
            return nil, newConvertError(rv.Type().String(), "slice(uniqueByField)", v, fmt.Sprintf("未找到字段或标签：%s", field))
        }
        ft := structType.Field(idx).Type
        useComparable := ft.Comparable()
        seenI := make(map[interface{}]struct{}, rv.Len())
        seenS := make(map[string]struct{}, rv.Len())

        for i := 0; i < rv.Len(); i++ {
            iv := rv.Index(i)
            var sv reflect.Value
            if isPtr {
                if iv.IsNil() {
                    // 关键步骤：nil指针作为特殊键，仅保留首次出现
                    key := uniqueMissing{}
                    if _, ok := seenI[key]; ok { continue }
                    seenI[key] = struct{}{}
                    out = reflect.Append(out, iv)
                    continue
                }
                sv = iv.Elem()
            } else {
                sv = iv
            }
            fv := sv.Field(idx)
            if useComparable {
                key := fv.Interface()
                if _, ok := seenI[key]; ok { continue }
                seenI[key] = struct{}{}
            } else {
                key := fmt.Sprintf("%#v", fv.Interface())
                if _, ok := seenS[key]; ok { continue }
                seenS[key] = struct{}{}
            }
            out = reflect.Append(out, iv)
        }
        return out.Interface(), nil
    }

    // 场景2：map元素（仅支持string键的map）
    if elemType.Kind() == reflect.Map {
        if elemType.Key().Kind() != reflect.String {
            return nil, newConvertError(rv.Type().String(), "slice(uniqueByField)", v, "仅支持键为string的map元素")
        }
        seenI := make(map[interface{}]struct{}, rv.Len())
        seenS := make(map[string]struct{}, rv.Len())
        fieldKey := reflect.ValueOf(field)
        for i := 0; i < rv.Len(); i++ {
            mv := rv.Index(i)
            val := mv.MapIndex(fieldKey)
            if !val.IsValid() {
                // 关键步骤：缺失键作为特殊键，仅保留首次出现
                key := uniqueMissing{}
                if _, ok := seenI[key]; ok { continue }
                seenI[key] = struct{}{}
                out = reflect.Append(out, mv)
                continue
            }
            if val.Type().Comparable() {
                key := val.Interface()
                if _, ok := seenI[key]; ok { continue }
                seenI[key] = struct{}{}
            } else {
                key := fmt.Sprintf("%#v", val.Interface())
                if _, ok := seenS[key]; ok { continue }
                seenS[key] = struct{}{}
            }
            out = reflect.Append(out, mv)
        }
        return out.Interface(), nil
    }

    return nil, newConvertError(rv.Type().String(), "slice(uniqueByField)", v, "仅支持结构体/结构体指针或map元素的切片")
}

// uniqueMissing 缺失键/空指针的占位符（私有辅助结构体）
// 结构体: 空结构体，作为map键以表示“缺失键”或“nil指针”
// 关键步骤：该类型可比较，用于seen集合的唯一键
type uniqueMissing struct{}

// fieldIndexByNameOrTag 查找结构体字段索引（支持按标签匹配）
// 参数 st: 结构体类型
// 参数 field: 字段名或标签值（当提供标签时优先匹配标签值）
// 参数 tag: 标签名（例如"json"；为空表示仅按字段名匹配）
// 返回值: 字段索引与是否找到
// 关键步骤：遍历导出字段，解析标签逗号前片段，与字段名或标签名进行匹配
func fieldIndexByNameOrTag(st reflect.Type, field string, tag string) (int, bool) {
    for i := 0; i < st.NumField(); i++ {
        sf := st.Field(i)
        if sf.PkgPath != "" { // 未导出字段跳过
            continue
        }
        if tag != "" {
            tv := sf.Tag.Get(tag)
            name := tv
            if idx := strings.IndexByte(tv, ','); idx >= 0 {
                name = tv[:idx]
            }
            if name == field && name != "" { return i, true }
        }
        if sf.Name == field { return i, true }
    }
    return -1, false
}