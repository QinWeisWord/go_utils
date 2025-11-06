package kvcache

import "sync"

// KVCache 一个简易的并发安全键值缓存器（泛型）
// 参数：无
// 返回值：无
// 关键步骤：使用读写锁保护内部map，提供基础的Set/Get/Delete/Has等操作
type KVCache[K comparable, V any] struct {
    mu    sync.RWMutex      // 关键步骤：并发读写保护
    store map[K]V           // 关键步骤：底层存储使用map
}

// New 创建一个新的键值缓存器实例
// 参数：无
// 返回值：缓存器指针
// 关键步骤：初始化内部map
func New[K comparable, V any]() *KVCache[K, V] {
    return &KVCache[K, V]{
        store: make(map[K]V),
    }
}

// Set 设置键的值（覆盖同名键）
// 参数 key: 键（必须可比较类型）
// 参数 value: 值
// 返回值：无
// 关键步骤：写入时加写锁，确保并发安全
func (c *KVCache[K, V]) Set(key K, value V) {
    c.mu.Lock()
    // 关键步骤：直接覆盖写入
    c.store[key] = value
    c.mu.Unlock()
}

// Get 获取键对应的值
// 参数 key: 键
// 返回值 value: 键对应的值（不存在时为零值）
// 返回值 ok: 是否存在该键
// 关键步骤：读取时加读锁，避免阻塞写操作的同时提升并发读性能
func (c *KVCache[K, V]) Get(key K) (value V, ok bool) {
    c.mu.RLock()
    v, ok := c.store[key]
    c.mu.RUnlock()
    return v, ok
}

// Delete 删除指定键
// 参数 key: 键
// 返回值 deleted: 删除是否成功（键是否存在）
// 关键步骤：写入时加写锁，使用多值删除语义判断是否存在
func (c *KVCache[K, V]) Delete(key K) (deleted bool) {
    c.mu.Lock()
    _, ok := c.store[key]
    if ok {
        // 关键步骤：存在时执行删除
        delete(c.store, key)
    }
    c.mu.Unlock()
    return ok
}

// Has 判断是否存在指定键
// 参数 key: 键
// 返回值 exists: 存在返回true
// 关键步骤：读锁保护并使用多值获取判断
func (c *KVCache[K, V]) Has(key K) (exists bool) {
    c.mu.RLock()
    _, ok := c.store[key]
    c.mu.RUnlock()
    return ok
}

// Len 返回当前缓存中的键数量
// 参数：无
// 返回值 n: 键的数量
// 关键步骤：读锁保护并返回map长度
func (c *KVCache[K, V]) Len() (n int) {
    c.mu.RLock()
    n = len(c.store)
    c.mu.RUnlock()
    return n
}

// Clear 清空缓存中的所有键
// 参数：无
// 返回值：无
// 关键步骤：用写锁将store替换为新的空map，避免逐个删除开销
func (c *KVCache[K, V]) Clear() {
    c.mu.Lock()
    c.store = make(map[K]V)
    c.mu.Unlock()
}

// Keys 返回当前所有键的切片
// 参数：无
// 返回值 keys: 键切片（顺序不保证）
// 关键步骤：在读锁下遍历map构建切片
func (c *KVCache[K, V]) Keys() (keys []K) {
    c.mu.RLock()
    keys = make([]K, 0, len(c.store))
    for k := range c.store {
        keys = append(keys, k)
    }
    c.mu.RUnlock()
    return keys
}

// Values 返回当前所有值的切片
// 参数：无
// 返回值 values: 值切片（顺序不保证）
// 关键步骤：在读锁下遍历map构建切片
func (c *KVCache[K, V]) Values() (values []V) {
    c.mu.RLock()
    values = make([]V, 0, len(c.store))
    for _, v := range c.store {
        values = append(values, v)
    }
    c.mu.RUnlock()
    return values
}