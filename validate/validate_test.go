package validate

import "testing"

// TestEmail 测试邮箱格式校验
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖常见有效邮箱与无效格式
func TestEmail(t *testing.T) {
    // 有效
    if !IsEmail("user@example.com") { t.Fatalf("基本邮箱应有效") }
    if !IsEmail("user.name+tag@example.co.uk") { t.Fatalf("带子域与+标签应有效") }
    if !IsEmail("USER@EXAMPLE.COM") { t.Fatalf("大写邮箱应有效") }
    // 无效
    if IsEmail("user@") { t.Fatalf("缺少域名不应有效") }
    if IsEmail("example") { t.Fatalf("无@不应有效") }
    if IsEmail("user@localhost") { t.Fatalf("无TLD不应有效") }
}

// TestMobileCN 测试中国大陆手机号
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖带/不带国家码的有效号码与无效号段
func TestMobileCN(t *testing.T) {
    if !IsMobileCN("13812345678") { t.Fatalf("基本手机号应有效") }
    if !IsMobileCN("+8613812345678") { t.Fatalf("+86前缀应有效") }
    if !IsMobileCN("86-13812345678") { t.Fatalf("86-分隔应有效") }
    if IsMobileCN("12345678901") { t.Fatalf("无效号段不应有效") }
    if IsMobileCN("1381234567") { t.Fatalf("长度不足不应有效") }
}

// TestURL 测试URL格式校验
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖带协议/补全协议、IP主机与localhost差异
func TestURL(t *testing.T) {
    if !IsURL("https://example.com") { t.Fatalf("https URL 应有效") }
    if !IsURL("example.com/path") { t.Fatalf("补全http后应有效") }
    if !IsURL("http://localhost") { t.Fatalf("带协议的localhost应有效") }
    if !IsURL("http://192.168.0.1") { t.Fatalf("IP主机应有效") }
    if IsURL("localhost:8080") { t.Fatalf("无协议的localhost不应有效") }
    if IsURL("ftp://example.com") { t.Fatalf("非http/https协议不应有效") }
}

// TestIP 测试IP地址校验
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖IPv4/IPv6与无效地址
func TestIP(t *testing.T) {
    if !IsIP("192.168.1.1") { t.Fatalf("IPv4 应有效") }
    if !IsIP("2001:db8::1") { t.Fatalf("IPv6 应有效") }
    if IsIP("999.999.999.999") { t.Fatalf("无效IPv4不应有效") }
}

// TestChineseIDCard 测试中国居民身份证号校验
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖合法示例、非法日期与错误校验位
func TestChineseIDCard(t *testing.T) {
    // 合法示例：计算校验位为X
    if !IsChineseIDCard("11010519491231002X") { t.Fatalf("合法身份证号应有效") }
    // 非法日期：不存在的日期
    if IsChineseIDCard("11010519990230002X") { t.Fatalf("非法日期不应有效") }
    // 错误校验位：最后一位不匹配
    if IsChineseIDCard("110105194912310021") { t.Fatalf("错误校验位不应有效") }
}

// TestPatentApplicationCN 测试中国专利申请号格式
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖可选前缀CN/ZL、10~12位数字+校验位的常见格式
func TestPatentApplicationCN(t *testing.T) {
    if !IsCNPatentApplicationNo("CN201430123456.7") { t.Fatalf("带CN前缀的申请号应有效") }
    if !IsCNPatentApplicationNo("201430123456.7") { t.Fatalf("无前缀申请号应有效") }
    if IsCNPatentApplicationNo("CN2014301234567.7") { t.Fatalf("过长数字不应有效") }
    if IsCNPatentApplicationNo("CN201430123456.") { t.Fatalf("缺少校验位不应有效") }
    if IsCNPatentApplicationNo("CN201430123456.78") { t.Fatalf("多位校验不应有效") }
}

// TestPatentNoCN 测试中国专利号/公开号格式
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：覆盖尾部类型字母与数字长度范围
func TestPatentNoCN(t *testing.T) {
    if !IsCNPatentNo("CN102345678A") { t.Fatalf("有效专利号应通过") }
    if !IsCNPatentNo("CN1234567B") { t.Fatalf("7位数字的专利号应通过") }
    if IsCNPatentNo("CN123456") { t.Fatalf("缺少尾部类型字母不应通过") }
    if IsCNPatentNo("CN12345678@") { t.Fatalf("尾部非字母字符不应通过") }
    if IsCNPatentNo("CN1234567890123A") { t.Fatalf("过长数字不应通过") }
}

// TestUnifiedSocialCreditCode 测试统一社会信用代码（USCC）校验
// 参数 t: 测试对象
// 返回值: 无
// 关键步骤：使用合法字符集与权重计算校验位的正例与负例
func TestUnifiedSocialCreditCode(t *testing.T) {
    // 构造合法代码：前17位为数字，最后一位为根据算法计算得到的校验字符（此处应为'8'）
    valid := "123456789012345678"
    if !IsUnifiedSocialCreditCode(valid) { t.Fatalf("合法USCC应有效: %s", valid) }
    // 非法：包含不允许字符I
    if IsUnifiedSocialCreditCode("I23456789012345678") { t.Fatalf("包含非法字符不应有效") }
    // 非法：校验位错误（篡改最后一位）
    if IsUnifiedSocialCreditCode("123456789012345679") { t.Fatalf("校验位错误不应有效") }
    // 非法：长度不为18
    if IsUnifiedSocialCreditCode("12345678901234567") { t.Fatalf("长度不足不应有效") }
}