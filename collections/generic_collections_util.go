package collections

// 本文件包含通用切片与映射的泛型工具函数

// Contains 判断切片中是否包含目标元素
// 参数 slice: 任意类型的切片（元素需可比较）
// 参数 target: 需要查找的目标元素
// 返回值: 布尔值，true表示找到，false表示未找到
// 关键步骤：遍历切片并按值比较
func Contains[T comparable](slice []T, target T) bool {
    for _, v := range slice {
        if v == target {
            return true
        }
    }
    return false
}

// IndexOf 查找目标元素在切片中的索引
// 参数 slice: 任意类型的切片（元素需可比较）
// 参数 target: 需要查找的目标元素
// 返回值: 索引位置；若未找到则返回-1
// 关键步骤：遍历切片并按值比较
func IndexOf[T comparable](slice []T, target T) int {
    for i, v := range slice {
        if v == target {
            return i
        }
    }
    return -1
}

// Unique 返回切片的去重副本（保持首次出现顺序）
// 参数 slice: 任意类型的切片（元素需可比较）
// 返回值: 去重后的新切片
// 关键步骤：使用map记录已出现元素
func Unique[T comparable](slice []T) []T {
    seen := make(map[T]struct{}, len(slice))
    out := make([]T, 0, len(slice))
    for _, v := range slice {
        if _, ok := seen[v]; !ok {
            seen[v] = struct{}{}
            out = append(out, v)
        }
    }
    return out
}

// Map 对切片进行映射转换
// 参数 slice: 输入切片
// 参数 fn: 将元素T转换为R的映射函数
// 返回值: 映射后的新切片
// 关键步骤：遍历并应用映射函数
func Map[T any, R any](slice []T, fn func(T) R) []R {
    out := make([]R, 0, len(slice))
    for _, v := range slice {
        out = append(out, fn(v))
    }
    return out
}

// Filter 过滤切片元素
// 参数 slice: 输入切片
// 参数 fn: 过滤函数，返回true表示保留，false表示丢弃
// 返回值: 过滤后的新切片
// 关键步骤：遍历并根据函数结果决定是否加入
func Filter[T any](slice []T, fn func(T) bool) []T {
    out := make([]T, 0, len(slice))
    for _, v := range slice {
        if fn(v) {
            out = append(out, v)
        }
    }
    return out
}

// Keys 返回map的所有键
// 参数 m: 任意键值类型的map（键需可比较）
// 返回值: 键的切片
// 关键步骤：遍历map收集键
func Keys[K comparable, V any](m map[K]V) []K {
    out := make([]K, 0, len(m))
    for k := range m {
        out = append(out, k)
    }
    return out
}

// Values 返回map的所有值
// 参数 m: 任意键值类型的map
// 返回值: 值的切片
// 关键步骤：遍历map收集值
func Values[K comparable, V any](m map[K]V) []V {
    out := make([]V, 0, len(m))
    for _, v := range m {
        out = append(out, v)
    }
    return out
}

// Merge 合并两个map（b覆盖a的同名键）
// 参数 a: 基础map
// 参数 b: 将被合并进来的map
// 返回值: 新的合并后map（不修改原map）
// 关键步骤：创建新map并依次拷贝/覆盖
func Merge[K comparable, V any](a, b map[K]V) map[K]V {
    out := make(map[K]V, len(a)+len(b))
    for k, v := range a {
        out[k] = v
    }
    for k, v := range b {
        out[k] = v
    }
    return out
}

// GetOrDefault 从map中获取值，不存在时返回默认值
// 参数 m: map对象
// 参数 key: 键
// 参数 def: 默认值
// 返回值: 键对应的值或默认值
// 关键步骤：使用多值获取判断键是否存在
func GetOrDefault[K comparable, V any](m map[K]V, key K, def V) V {
    if v, ok := m[key]; ok {
        return v
    }
    return def
}