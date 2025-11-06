package main

import (
    "log"
    "net/http"
    "os"
    "strconv"

    "go_utils/captcha"
)

// parseInt 从查询参数解析整数（带默认值）
// 参数 r: HTTP请求
// 参数 key: 查询参数键名
// 参数 defVal: 默认值
// 返回值: 解析后的整数（失败时返回默认值）
func parseInt(r *http.Request, key string, defVal int) int {
    v := r.URL.Query().Get(key)
    if v == "" { return defVal }
    n, err := strconv.Atoi(v)
    if err != nil { return defVal }
    return n
}

// parseFloat 从查询参数解析浮点数（带默认值）
// 参数 r: HTTP请求
// 参数 key: 查询参数键名
// 参数 defVal: 默认值
// 返回值: 解析后的浮点数（失败时返回默认值）
func parseFloat(r *http.Request, key string, defVal float64) float64 {
    v := r.URL.Query().Get(key)
    if v == "" { return defVal }
    n, err := strconv.ParseFloat(v, 64)
    if err != nil { return defVal }
    return n
}

// textHandler 文本图片验证码预览处理器
// 参数 w: 响应写入器
// 参数 r: HTTP请求，支持查询参数：text/width/height/lines/dots/ttf
// 返回值: 无（直接写出PNG或错误）
// 关键步骤：读取参数→生成图片→写出PNG响应
func textHandler(w http.ResponseWriter, r *http.Request) {
    text := r.URL.Query().Get("text")
    if text == "" { text = "Abc9Z" }
    width := parseInt(r, "width", 180)
    height := parseInt(r, "height", 60)
    lines := parseInt(r, "lines", 4)
    dots := parseInt(r, "dots", 200)
    // 新增：字体缩放比例（相对高度的默认尺寸），范围建议 0.6~2.0
    scale := parseFloat(r, "scale", 1.0)
    if scale <= 0 { scale = 1.0 }
    if scale < 0.6 { scale = 0.6 }
    if scale > 2.0 { scale = 2.0 }

    // 可选TTF字体：通过 ttf 参数传入文件路径
    var fontBytes []byte
    if fp := r.URL.Query().Get("ttf"); fp != "" {
        if b, err := os.ReadFile(fp); err == nil {
            fontBytes = b
        }
    }

    img, err := captcha.GenerateTextCaptchaImagePNGWithScale(text, width, height, lines, dots, fontBytes, scale)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.Header().Set("Content-Type", "image/png")
    _, _ = w.Write(img)
}

// main 启动本地预览服务器
// 参数：无
// 返回值：无
// 关键步骤：注册路由并监听 8080 端口
func main() {
    http.HandleFunc("/text.png", textHandler)
    log.Println("文本验证码预览: http://localhost:8080/text.png?text=Abc9Z&width=180&height=60")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}