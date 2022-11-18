package snake

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
)

var Len = utf8.RuneCountInString

// ucfirst 英文首字母大写 ...
func ucfirst(src string) string {
	for i, v := range src {
		return string(unicode.ToUpper(v)) + src[i+1:]
	}
	return src
}

// lcfirst 英文首字母小写 ...
func lcfirst(src string) string {
	for i, v := range src {
		return string(unicode.ToLower(v)) + src[i+1:]
	}
	return src
}

// WalkPath Files……
// 遍历目录查找文件
func walkPath(path string, dst ...string) []string {
	var res []string
	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() {
			for _, v := range dst {
				if l, err := filepath.Glob(filepath.Join(p, filepath.Base(v))); len(l) != 0 && err == nil {
					res = append(res, l...)
				}
			}
		}
		return err
	})
	return res
}

// ls 路径目录下内容
func ls(path string, dst ...string) []string {
	var res []string
	for _, v := range dst {
		if l, err := filepath.Glob(filepath.Join(path, v)); err == nil {
			res = append(res, l...)
		}
	}
	return res
}

// _owcpfile 路径目录下内容
func _owcpfile(src FileSystem, dst FileSystem) bool {
	// 覆盖拷贝
	if f, ok := dst.MkFile(); ok {
		defer f.Get().Close()
		if s, ok := src.Open(); ok {
			defer s.Get().Close()
			_, err := io.Copy(f.Get(), s.Get())
			return err == nil
		}
	}
	return false
}

func getEncoding(charset string) encoding.Encoding {
	if e, err := ianaindex.MIB.Encoding(charset); err == nil && e != nil {
		return e
	}
	return nil
}

func prenum(data byte) int {
	str := fmt.Sprintf("%b", data)
	var i int = 0
	for i < len(str) {
		if str[i] != '1' {
			break
		}
		i++
	}
	return i
}
