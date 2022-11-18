package snake

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/mycalf/snake/pkg"
	"github.com/yuin/charsetutil"
	"golang.org/x/text/transform"
	"golang.org/x/text/width"
)

type SnakeString struct {
	Input string
}

// ---------------------------------------
// 输入 :

// Text 初始化...
func String(str ...interface{}) *SnakeString {
	t := &SnakeString{}
	if len(str) > 0 {
		t.Add(str...)
	}
	return t
}

// Add 在字符串中追加文字...
func (t *SnakeString) Add(str ...interface{}) *SnakeString {
	b := bytes.NewBufferString(t.Input)
	if len(str) > 0 {
		for _, v := range str {
			b.WriteString(fmt.Sprint(v))
		}
	}
	t.Input = b.String()
	return t
}

// AddSlice 通过Slice在字符串中追加文字...
func (t *SnakeString) AddSlice(dts []interface{}) *SnakeString {
	for _, v := range dts {
		t.Add(v)
	}
	return t
}

// LN 回车...
func (t *SnakeString) Ln(line ...int) *SnakeString {
	if len(line) > 0 {
		for i := 0; i < line[0]; i++ {
			t.Add("\n")
		}
		return t
	}

	return t.Add("\n")
}

// ---------------------------------------
// 处理 :

// Replace 替换字符串或符合正则规则的字符串 ...
// snake.Text("http://www.dedecms.com").Replace("(http://).*(dedecms.com)", "${1}${2}")
// out: http://dedecms.com
// snake.Text("http://www.example.com").Replace("example", "dedecms")
// out: http://www.dedecms.com
// 如需替换$等字符，请使用\\$
// snake.Text("http://$1example.com").Replace("\\$1.*(.com)", "www.dedecms${1}")
func (t *SnakeString) Replace(src, dst string, noreg ...bool) *SnakeString {
	if len(noreg) > 0 && noreg[0] {
		t.Input = strings.Replace(t.Input, src, dst, -1)
		return t
	}
	t.Input = regexp.MustCompile(src).ReplaceAllString(t.Input, dst)
	return t
}

// ReplaceOne 替换出现的第一个字符串 ...
func (t *SnakeString) ReplaceOne(src, dst string) *SnakeString {
	t.Input = strings.Replace(t.Input, src, dst, 1)
	return t
}

// Find 判断字符串或符合正则规则的字符串是否存在 ...
func (t *SnakeString) Find(dst string, noreg ...bool) bool {

	if len(noreg) > 0 && noreg[0] {
		return strings.Contains(t.Input, dst)
	}

	if d := regexp.MustCompile(dst).FindAll([]byte(t.Input), -1); len(d) > 0 {
		return true
	}
	return false
}

// Remove 根据正则规则删除字符串 ...
func (t *SnakeString) Remove(dst ...string) *SnakeString {
	if len(dst) > 0 {
		for _, v := range dst {
			t.Input = regexp.MustCompile(v).ReplaceAllString(t.Input, "")
		}
	}

	return t
}

// ExistSlice 字符串是否存在于数组中 ...
func (t *SnakeString) ExistSlice(dst []string) bool {
	if len(dst) > 0 {
		for _, v := range dst {
			if t.Get() == v {
				return true
			}
		}
	}
	return false
}

// Keep 根据正则规则保留字符串 ...
func (t *SnakeString) Keep(dst string) *SnakeString {

	if t.Find(dst) {
		p := String()
		d := regexp.MustCompile(dst).FindAll([]byte(t.Get()), -1)

		for _, v := range d {
			p.Add(string(v))
		}

		t.Input = p.Get()
	}

	return t
}

// Extract 根据正则规则提取字符数组 ...
func (t *SnakeString) Extract(dst string, out ...string) []string {
	arr := []string{}
	if t.Find(dst) {
		d := regexp.MustCompile(dst).FindAll([]byte(t.Get()), -1)
		if len(out) > 0 && out[0] != "" {
			for _, s := range out {
				for _, v := range d {
					arr = append(arr, String(string(v)).Replace(dst, s).Get())
				}
			}
			return arr
		}

		for _, v := range d {
			arr = append(arr, string(v))
		}

	}
	return arr
}

// Narrow 全角字符转半角字符 ...
func (t *SnakeString) Narrow() *SnakeString {
	t.Input = width.Narrow.String(t.Input)
	return t
}

// Widen 半角字符转全角字符 ...
func (t *SnakeString) Widen() *SnakeString {
	t.Input = width.Narrow.String(t.Input)
	return t
}

// ReComment 去除代码注解...
func (t *SnakeString) ReComment() *SnakeString {
	t.Remove(
		`\/\/.*`,
		`\/\*(\s|.)*?\*\/`,
		`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
	return t
}

// Trim 去除开始及结束出现的字符 ...
func (t *SnakeString) Trim(sep string) *SnakeString {
	t.Input = strings.Trim(t.Input, sep)
	return t
}

// ToLower 字符串全部小写 ...
func (t *SnakeString) ToLower() *SnakeString {
	t.Input = strings.ToLower(t.Input)
	return t
}

// ToUpper 字符串全部小写 ...
func (t *SnakeString) ToUpper() *SnakeString {
	t.Input = strings.ToUpper(t.Input)
	return t
}

// UCFirst 字符串首字母大写 ...
func (t *SnakeString) UcFirst() *SnakeString {
	t.Input = ucfirst(t.Input)
	return t
}

// LCFirst 字符串首字母小写 ...
func (t *SnakeString) LcFirst() *SnakeString {
	t.Input = lcfirst(t.Input)
	return t
}

// Between 截取区间内容 ...
func (t *SnakeString) Between(start, end string) *SnakeString {
	if (start == "" && end == "") || t.Input == "" {
		return t
	}
	// 处理数据，将所有文字转为小写 .
	input := strings.ToLower(t.Input)
	lowerStart := strings.ToLower(start)
	lowerEnd := strings.ToLower(end)

	var startIndex, endIndex int

	if len(start) > 0 && strings.Contains(input, lowerStart) {
		startIndex = len(start)
	}
	if len(end) > 0 && strings.Contains(input, lowerEnd) {
		endIndex = strings.Index(input, lowerEnd)
	} else if len(input) > 0 {
		endIndex = len(input)
	}
	// 输出字符A与字符B之间的字符 .
	t.Input = strings.TrimSpace(t.Input[startIndex:endIndex])
	return t
}

// EnBase Text to Base-x:  2 < base > 36 ...
// 将Text转为2～36进制编码
func (t *SnakeString) EnBase(base int) *SnakeString {
	var r []string
	for _, i := range t.Input {
		r = append(r, strconv.FormatInt(int64(i), base))
	}
	t.Input = strings.Join(r, " ")
	return t
}

// DeBase Text Base-x to Text:  2 < base > 36 ...
// 将2～36进制解码为Text
func (t *SnakeString) DeBase(base int) *SnakeString {
	var r []rune
	for _, i := range t.Split(" ") {
		i64, err := strconv.ParseInt(i, base, 64)
		if err != nil {
			panic(err)
		}
		r = append(r, rune(i64))
	}
	t.Input = string(r)
	return t
}

// ---------------------------------------
// 分词 :

// CamelCase 驼峰分词: HelloWord ...
func (t *SnakeString) CamelCase() *SnakeString {
	caseWords := t.caseWords(true)
	for i, word := range caseWords {
		caseWords[i] = ucfirst(word)
	}
	t.Input = strings.Join(caseWords, "")
	return t
}

// SnakeCase 贪吃蛇分词: hello_word ...
func (t *SnakeString) SnakeCase() *SnakeString {
	caseWords := t.caseWords(false)
	t.Input = strings.Join(caseWords, "_")
	return t
}

// KebabCase "烤串儿"分词: hello-word ...
func (t *SnakeString) KebabCase() *SnakeString {
	caseWords := t.caseWords(false)
	t.Input = strings.Join(caseWords, "-")
	return t
}

// ---------------------------------------
// 输出 :

// Get 获取文本...
func (t *SnakeString) Get() string {
	return t.Input
}

// GetOneLine 在多行字符串中获取第一行字符串...
func (t *SnakeString) GetOneLine() string {
	for _, v := range t.Lines() {
		return v
	}
	return t.Input
}

// 以LF格式输出...
func (t *SnakeString) LF() *SnakeString {
	return t.Replace("\r\n", "\n", true)
}

// Byte Function
// 获取字符串的Byte ...
func (t *SnakeString) Byte() []byte {
	return []byte(t.Input)
}

// Split 根据字符串进行文本分割 ...
func (t *SnakeString) Split(sep string) []string {
	return strings.Split(t.Input, sep)
}

// SplitPlace 根据字符串的位置进行分割
// Text("abcdefg").SpltPlace([]int{1,3,4})
// Out: []string{"a", "bc", "d", "efg"}
func (t *SnakeString) SplitPlace(sep []int) []string {
	var a []string
	b := String()
	for k, v := range []rune(t.Input) {
		b.Add(string(v))
		for _, i := range sep {
			if i == k+1 {
				a = append(a, b.Get())
				b = String()
			}
		}

		if len(t.Input) == k+1 {
			a = append(a, b.Get())
		}
	}
	return a
}

// SplitInt 根据设置对字符串等分
// Text("abcdefg").SpltPlace([]int{1,3,4})
// Out: []string{"a", "bc", "d", "efg"}
func (t *SnakeString) SplitInt(sep int) []string {
	var a []string
	b := String()
	i := 0
	for _, v := range t.Input {
		b.Add(string(v))

		i = i + len(string(v))

		bl := len(b.Get())

		if bl >= sep || i == len(t.Get()) {
			a = append(a, b.Get())
			b = String()
		}
	}

	return a
}

// ss 根据行进行分割字符 ...
func (t *SnakeString) Lines() []string {
	return strings.Split(strings.TrimSuffix(t.Input, "\n"), "\n")
}

// MD5 获取文件的MD5
func (t *SnakeString) MD5() string {
	return fmt.Sprintf("%x", md5.Sum([]byte(t.Get())))
}

// 根据length循环复制字符串儿
func (t *SnakeString) Copy(length int) string {
	str := String()
	for i := 0; i < length; i++ {
		str.Add(t.Input)
	}
	return str.Get()
}

// 根据文字自动绘制代码提示框.
func (t *SnakeString) DrawBox(width int, chars ...pkg.Box9Slice) *SnakeString {

	res := String()
	char := pkg.DefaultBox9Slice()
	if len(chars) == 1 {
		char = chars[0]
	}

	var topInsideWidth = width - Len(char.TopLeft) - Len(char.TopRight)
	var bottomInsideWidth = width - Len(char.BottomLeft) - Len(char.BottomRight)

	if topInsideWidth < 1 || bottomInsideWidth < 1 {
		topInsideWidth = 60
		bottomInsideWidth = 60
	}

	lines := t.Lines()

	//top
	res.Add(char.TopLeft).
		Add(String(char.Top).Copy(topInsideWidth)).
		Add(char.TopRight).Ln()

	//middle
	for _, line := range lines {
		res.Add(char.Left).Add(" ").Add(String(line).Trim(" ").Get()).Ln()
	}

	//bottom
	res.Add(char.BottomLeft).
		Add(String(char.Bottom).Copy(bottomInsideWidth)).
		Add(char.BottomRight)

	t.Input = res.Get()

	return t
}

func (t *SnakeString) Unescape() string {
	if html, err := url.QueryUnescape(String(t.Get()).Replace(`%u(.{4})`, "/u$1/").Get()); err == nil {
		temp := String(html)
		for _, v := range t.Extract(`/u(.{4})?/`) {
			if w, err := strconv.Unquote(`"` + String(v).Replace(`/u(.{4})?/`, "\\u$1").Get() + `"`); err == nil {
				temp.Replace(v, w, true)
			}
		}
		return temp.Get()
	}
	return t.Get()
}

// ---------------------------------------
// 字符集 :

// Charset Function
// 返回当前进程的字符集 ...
func (t *SnakeString) Charset() (string, bool) {

	// 自动获取编码 ...
	encoding, err := charsetutil.GuessBytes(t.Byte())

	// 如果自动获取成功或encoding不为空
	// 则输出编码格式 ...
	if err == nil {
		return strings.ToUpper(encoding.Charset()), true
	}

	if t.IsGBK() {
		return "GBK", true
	}

	if encoding != nil && encoding.Charset() != "WINDOWS-1252" {
		return strings.ToUpper(encoding.Charset()), true
	}

	// 如果内容中出现汉字
	// 则输出GB18030 ...
	if t.ExistHan() {
		return "GBK", true
	}

	// 不符合上述条件
	// 则返回空 ...
	return "", false
}

// ExistHan Function
// 判断是否存在中文 ...
func (t *SnakeString) ExistHan() bool {
	hanLen := len(regexp.MustCompile(`[\P{Han}]`).ReplaceAllString(t.Input, ""))
	for _, r := range t.Input {
		if unicode.Is(unicode.Scripts["Han"], r) || hanLen > 0 {
			return true
		}
	}
	return false
}

// ExistGBK Function
// 判断是否为GBK ...
func (t *SnakeString) IsGBK() bool {
	arr := t.Byte()
	var i int = 0
	for i < len(t.Byte()) {
		if arr[i] <= 0xff {
			i++
			continue
		} else {
			if arr[i] >= 0x81 &&
				arr[i] <= 0xfe &&
				arr[i+1] >= 0x40 &&
				arr[i+1] <= 0xfe &&
				arr[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}

// 判断是否为UTF8
func (t *SnakeString) IsUTF8() bool {
	data := t.Byte()
	for i := 0; i < len(data); {
		if data[i]&0x80 == 0x00 {
			i++
			continue
		} else if num := prenum(data[i]); num > 2 {
			i++
			for j := 0; j < num-1; j++ {
				if data[i]&0xc0 != 0x80 {
					return false
				}
				i++
			}
		} else {
			return false
		}
	}
	return true
}

// ToUTF8 Function
// 运行对当前进程进行编码转换成UTF-8 ...
func (t *SnakeString) ToUTF8() (string, bool) {

	// 自动获取资源编码 ...
	charset, ok := t.Charset()

	// 未获取到资源编码 ...
	if !ok {
		return t.Input, false
	}

	// UTF-8无需转换 ...
	if charset == "UTF-8" {
		return t.Input, true
	}

	if encode := getEncoding(charset); encode != nil {
		if reader, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader(t.Byte()), encode.NewDecoder())); err == nil {
			t.Input = string(reader)
			t.Input = html.UnescapeString(t.Input)
			return t.Input, true
		}
	}

	// 转码失败
	return t.Input, false
}

// LCFirst 字符串首字母小写 ...
func (t *SnakeString) Write(dst string, add ...bool) bool {
	return FS(dst).Write(t.Get(), add...)
}

// ---------------------------------------
// 辅助函数 :

// 根据规则字符进行分词 ...
func (t *SnakeString) caseWords(isCamel bool, rule ...string) []string {
	src := t.Input
	if !isCamel {
		re := regexp.MustCompile("([a-z])([A-Z])")
		src = re.ReplaceAllString(src, "$1 $2")
	}
	src = strings.Join(strings.Fields(strings.TrimSpace(src)), " ")
	rule = append(rule, ".", " ", "_", " ", "-", " ")
	replacer := strings.NewReplacer(rule...)
	src = replacer.Replace(src)
	return strings.Fields(src)
}
