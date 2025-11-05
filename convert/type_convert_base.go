package convert

import (
    "fmt"
)

// 本文件提供类型转换的统一错误类型与构造函数

// ConvertError 表示一次类型转换失败的错误
// 结构体字段 FromType: 源类型名称；ToType: 目标类型名称；Value: 源值的字符串表示；Message: 错误说明
// 关键步骤：实现error接口用于统一错误返回
type ConvertError struct {
    FromType string
    ToType   string
    Value    string
    Message  string
}

// Error 返回错误信息字符串
// 参数: 无
// 返回值: 错误信息字符串
// 关键步骤：格式化错误包含源类型、目标类型、值与说明
func (e *ConvertError) Error() string {
    return fmt.Sprintf("无法将类型 %s 的值 %q 转换为类型 %s：%s", e.FromType, e.Value, e.ToType, e.Message)
}

// newConvertError 创建统一的转换错误
// 参数 fromType: 源类型名称
// 参数 toType: 目标类型名称
// 参数 value: 源值
// 参数 msg: 错误说明
// 返回值: *ConvertError 错误对象
// 关键步骤：将值转字符串并填充错误结构体
func newConvertError(fromType, toType string, value interface{}, msg string) *ConvertError {
    return &ConvertError{
        FromType: fromType,
        ToType:   toType,
        Value:    fmt.Sprintf("%v", value),
        Message:  msg,
    }
}