package filejson

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Exists 判断文件或目录是否存在
// 参数 path: 文件或目录路径
// 返回值: 布尔值，true表示存在，false表示不存在
// 关键步骤：调用os.Stat并根据错误是否为IsNotExist判定
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

// IsDir 判断路径是否为目录
// 参数 path: 文件或目录路径
// 返回值: 布尔值，true表示目录，false表示不是目录或不存在
// 关键步骤：调用os.Stat并检查FileInfo.IsDir
func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

// ReadFile 读取文本文件内容为字符串
// 参数 path: 文件路径
// 返回值: 文件内容字符串与错误信息（若成功则error为nil）
// 关键步骤：调用os.ReadFile读取全部内容
func ReadFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// WriteFile 将字符串写入到文件（自动创建父目录）
// 参数 path: 文件路径
// 参数 content: 要写入的文本内容
// 返回值: 错误信息（若成功则error为nil）
// 关键步骤：确保父目录存在，再写入文件
func WriteFile(path string, content string) error {
	// 关键步骤：创建父目录
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	// 关键步骤：写入文件
	return os.WriteFile(path, []byte(content), 0o644)
}

// CopyFile 拷贝文件（若目标父目录不存在则创建）
// 参数 src: 源文件路径
// 参数 dst: 目标文件路径
// 返回值: 错误信息（若成功则error为nil）
// 关键步骤：打开源文件，创建目标文件，使用io.Copy复制内容
func CopyFile(src, dst string) error {
	// 关键步骤：检查源文件是否存在且非目录
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return errors.New("source path is a directory")
	}

	// 关键步骤：确保目标父目录存在
	if err = os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	// 关键步骤：打开源文件
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// 关键步骤：创建目标文件
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		// 关键步骤：确保关闭并刷写
		_ = out.Close()
	}()

	// 关键步骤：复制内容
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return nil
}

// ListDir 列出目录下的文件/子目录名称（不递归）
// 参数 dir: 目录路径
// 返回值: 名称切片与错误信息（若成功则error为nil）
// 关键步骤：调用os.ReadDir并收集Name
func ListDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names, nil
}

// ToJSON 将任意对象编码为JSON字符串
// 参数 v: 需要编码的任意对象
// 返回值: JSON字符串与错误信息（若成功则error为nil）
// 关键步骤：调用json.Marshal
func ToJSON(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ToPrettyJSON 将对象编码为格式化的JSON字符串
// 参数 v: 需要编码的任意对象
// 返回值: 美化后的JSON字符串与错误信息（若成功则error为nil）
// 关键步骤：调用json.MarshalIndent
func ToPrettyJSON(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FromJSON 解析JSON字符串为目标类型（泛型）
// 参数 data: JSON字符串
// 返回值: 解析得到的目标类型实例与错误信息（若成功则error为nil）
// 关键步骤：调用json.Unmarshal并返回泛型类型
func FromJSON[T any](data string) (T, error) {
	var v T
	if len(data) == 0 {
		var zero T
		return zero, errors.New("empty json data")
	}
	err := json.Unmarshal([]byte(data), &v)
	return v, err
}
