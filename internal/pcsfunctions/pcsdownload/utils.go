package pcsdownload

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsconfig"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/checksum"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

// FileMD5CheckResult 保存本地文件与网盘文件的 MD5 校验值。
type FileMD5CheckResult struct {
	LocalMD5  string
	RemoteMD5 string
}

// CheckFileValid 检测文件有效性
func CheckFileValid(filePath string, fileInfo *baidupcs.FileDirectory) error {
	_, err := CheckFileMD5(filePath, fileInfo)
	return err
}

// CheckFileMD5 计算本地文件 MD5，并与网盘记录的完整文件 MD5 比较。
func CheckFileMD5(filePath string, fileInfo *baidupcs.FileDirectory) (*FileMD5CheckResult, error) {
	if fileInfo == nil {
		return nil, ErrDownloadFileInfoNil
	}
	if len(fileInfo.BlockList) != 1 {
		return nil, ErrDownloadNotSupportChecksum
	}

	remoteMD5 := strings.ToLower(fileInfo.MD5)
	decodedMD5, err := hex.DecodeString(remoteMD5)
	if err != nil || len(decodedMD5) != 16 {
		return nil, ErrDownloadNotSupportChecksum
	}

	f := checksum.NewLocalFileChecksum(filePath, int(baidupcs.SliceMD5Size))
	err = f.OpenPath()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = f.Sum(checksum.CHECKSUM_MD5)
	if err != nil {
		return nil, err
	}
	md5Str := hex.EncodeToString(f.MD5)
	result := &FileMD5CheckResult{
		LocalMD5:  md5Str,
		RemoteMD5: remoteMD5,
	}

	if md5Str != remoteMD5 { // md5不一致
		// 检测是否为违规文件
		if IsSkipMd5Checksum(f.Length, md5Str) {
			return result, ErrDownloadFileBanned
		}
		return result, ErrDownloadChecksumFailed
	}
	return result, nil
}

// FileExist 检查文件是否存在,
// 只有当文件存在, 文件大小不为0或断点续传文件不存在时, 才判断为存在
func FileExist(path string) bool {
	if info, err := os.Stat(path); err == nil {
		if info.Size() == 0 {
			return false
		}
		if _, err = os.Stat(path + DownloadSuffix); err != nil {
			return true
		}
	}

	return false
}

// FixHTTPLinkURL 通过配置, 确定链接使用的协议(http,https)
func FixHTTPLinkURL(linkURL *url.URL) {
	if pcsconfig.Config.EnableHTTPS {
		if linkURL.Scheme == "http" {
			linkURL.Scheme = "https"
		}
	}
}

func CloneJarWithDomain(srcJar http.CookieJar, newURL string) (http.CookieJar, error) {
	if srcJar == nil {
		return nil, fmt.Errorf("srcJar is nil")
	}
	dstJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(newURL)
	newDomain := u.Hostname()
	u, _ = url.Parse("https://" + pcsconfig.Config.PCSAddr + "/")
	cookies := srcJar.Cookies(u)
	for _, c := range cookies {
		nc := *c
		nc.Domain = newDomain
		newURL, _ := url.Parse("https://" + newDomain + "/")
		dstJar.SetCookies(newURL, []*http.Cookie{&nc})
	}
	return dstJar, nil
}
