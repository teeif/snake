package snake

import (
	"bytes"
	"os"
)

type snakefile struct {
	Input *os.File
}

// FileOperate ...
type FileOperate interface {
	Get() *os.File
	String() *SnakeString
	Byte() []byte
	Close() error // 关闭文件链接
}

// ---------------------------------------
// 输入 :

// File 初始化...
func File(f *os.File) FileOperate {
	return &snakefile{Input: f}
}

// ---------------------------------------
// 输出 :

// Get 获取文本...
func (sk *snakefile) Get() *os.File {
	return sk.Input
}

// Add 在字符串中追加文字...
func (sk *snakefile) Close() error {
	return sk.Input.Close()
}

// Text 获取文本...
func (sk *snakefile) String() *SnakeString {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(sk.Input)
	if err != nil {
		// todo: 字符串转化错误消息
		return String()
	}

	return String(buf.String())
}

// Text 获取文本 []byte ...
func (sk *snakefile) Byte() []byte {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(sk.Input)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}
