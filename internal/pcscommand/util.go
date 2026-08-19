package pcscommand

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type pathIdentifierKind uint8

const (
	pathIdentifierNone pathIdentifierKind = iota
	pathIdentifierFsID
	pathIdentifierMD5
)

var (
	// ErrShellPatternMultiRes 多条通配符匹配结果
	ErrShellPatternMultiRes = errors.New("多条通配符匹配结果")
	// ErrShellPatternNoHit 未匹配到路径
	ErrShellPatternNoHit = errors.New("未匹配到路径, 请检测通配符")
)

// ListTask 队列状态 (基类)
type ListTask struct {
	ID       int // 任务id
	MaxRetry int // 最大重试次数
	retry    int // 任务失败的重试次数
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RunTestShellPattern 执行测试通配符
func RunTestShellPattern(pattern string) {
	pcs := GetBaiduPCS()
	paths, err := pcs.MatchPathByShellPattern(GetActiveUser().PathJoin(pattern))
	if err != nil {
		fmt.Println(err)
		return
	}
	for k := range paths {
		fmt.Printf("%s\n", paths[k])
	}
	return
}

func matchPathByShellPatternOnce(pattern *string) error {
	resolvedPath, err := resolveCurrentPathIdentifier(*pattern)
	if err != nil {
		return err
	}
	*pattern = resolvedPath

	paths, err := GetBaiduPCS().MatchPathByShellPattern(GetActiveUser().PathJoin(*pattern))
	if err != nil {
		return err
	}
	switch len(paths) {
	case 0:
		return ErrShellPatternNoHit
	case 1:
		*pattern = paths[0]
	default:
		return ErrShellPatternMultiRes
	}

	return nil
}

func matchPathByShellPattern(patterns ...string) (pcspaths []string, err error) {
	acUser, pcs := GetActiveUser(), GetBaiduPCS()
	for k := range patterns {
		resolvedPath, err := resolveCurrentPathIdentifier(patterns[k])
		if err != nil {
			return nil, err
		}

		ps, err := pcs.MatchPathByShellPattern(acUser.PathJoin(resolvedPath))
		if err != nil {
			return nil, err
		}

		pcspaths = append(pcspaths, ps...)
	}
	return pcspaths, nil
}

func resolveCurrentPathIdentifier(target string) (string, error) {
	kind, _, err := parsePathIdentifier(target)
	if err != nil || kind == pathIdentifierNone {
		return target, err
	}

	activeUser := GetActiveUser()
	files, pcsError := GetBaiduPCS().CacheFilesDirectoriesList(activeUser.Workdir, baidupcs.DefaultOrderOptions)
	if pcsError != nil {
		return "", pcsError
	}

	return resolvePathIdentifier(target, activeUser.Workdir, files)
}

func resolvePathIdentifier(target, currentDirectory string, files baidupcs.FileDirectoryList) (string, error) {
	kind, identifier, err := parsePathIdentifier(target)
	if err != nil || kind == pathIdentifierNone {
		return target, err
	}

	matches := make(baidupcs.FileDirectoryList, 0, 1)
	for _, file := range files {
		if file == nil {
			continue
		}

		matched := false
		switch kind {
		case pathIdentifierFsID:
			matched = strconv.FormatInt(file.FsID, 10) == identifier
		case pathIdentifierMD5:
			matched = strings.EqualFold(file.MD5, identifier)
		}
		if matched {
			matches = append(matches, file)
		}
	}

	identifierName := "fs_id"
	if kind == pathIdentifierMD5 {
		identifierName = "MD5"
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("错误: 当前目录 %s 中未找到 %s=%s", currentDirectory, identifierName, identifier)
	}
	if len(matches) > 1 {
		if kind == pathIdentifierMD5 {
			return "", fmt.Errorf("错误: 当前目录 %s 中有 %d 个项目匹配 MD5=%s, 请改用 fs_id", currentDirectory, len(matches), identifier)
		}
		return "", fmt.Errorf("错误: 当前目录 %s 中有 %d 个项目匹配 fs_id=%s, 目录数据异常", currentDirectory, len(matches), identifier)
	}
	if matches[0].Path != "" {
		return matches[0].Path, nil
	}
	return path.Join(currentDirectory, matches[0].Filename), nil
}

func parsePathIdentifier(target string) (pathIdentifierKind, string, error) {
	if target == "" {
		return pathIdentifierNone, "", nil
	}

	switch target[0] {
	case '$':
		identifier := target[1:]
		fsID, err := strconv.ParseInt(identifier, 10, 64)
		if err != nil || fsID <= 0 {
			return pathIdentifierFsID, "", fmt.Errorf("错误: 无效的 fs_id 标识符 %q", target)
		}
		return pathIdentifierFsID, strconv.FormatInt(fsID, 10), nil
	case '%':
		identifier := target[1:]
		decoded, err := hex.DecodeString(identifier)
		if err != nil || len(decoded) != 16 {
			return pathIdentifierMD5, "", fmt.Errorf("错误: 无效的 MD5 标识符 %q, 应为 32 位十六进制字符串", target)
		}
		return pathIdentifierMD5, strings.ToLower(identifier), nil
	default:
		return pathIdentifierNone, "", nil
	}
}

func randReplaceStr(s string, rname bool) string {
	if !rname {
		return s
	}
	filenameAll := path.Base(s)
	fileSuffix := path.Ext(s)
	filePrefix := filenameAll[0 : len(filenameAll)-len(fileSuffix)]
	runes := []rune(filePrefix)

	for i := 0; i < len(filePrefix); i++ {
		runes[i] = rune(letters[rand.Int63()%int64(len(letters))])
		if i == 3 {
			break
		}
	}
	return path.Join(path.Dir(s), string(runes)+fileSuffix)
}
