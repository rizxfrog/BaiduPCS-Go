package pcscommand

import (
	"errors"
	"fmt"
	"os"

	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcsfunctions/pcsdownload"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
)

// RunVerify 校验本地文件与网盘文件的大小和完整文件 MD5。
func RunVerify(remotePath, localPath string) error {
	if err := matchPathByShellPatternOnce(&remotePath); err != nil {
		return fmt.Errorf("解析网盘路径失败: %w", err)
	}

	fileInfo, pcsError := GetBaiduPCS().FilesDirectoriesMeta(remotePath)
	if pcsError != nil {
		return pcsError
	}
	if fileInfo.Isdir {
		return fmt.Errorf("网盘路径是目录，verify 仅支持文件: %s", remotePath)
	}

	localInfo, err := os.Stat(localPath)
	if err != nil {
		return fmt.Errorf("读取本地文件失败: %w", err)
	}
	if localInfo.IsDir() {
		return fmt.Errorf("本地路径是目录，verify 仅支持文件: %s", localPath)
	}

	fmt.Printf("网盘文件: %s\n", fileInfo.Path)
	fmt.Printf("本地文件: %s\n", localPath)
	fmt.Printf("网盘大小: %d B (%s)\n", fileInfo.Size, converter.ConvertFileSize(fileInfo.Size, 2))
	fmt.Printf("本地大小: %d B (%s)\n", localInfo.Size(), converter.ConvertFileSize(localInfo.Size(), 2))
	if localInfo.Size() != fileInfo.Size {
		return fmt.Errorf("校验失败: 文件大小不一致")
	}

	result, err := pcsdownload.CheckFileMD5(localPath, fileInfo)
	if result != nil {
		fmt.Printf("网盘 MD5: %s\n", result.RemoteMD5)
		fmt.Printf("本地 MD5: %s\n", result.LocalMD5)
	}
	if err != nil {
		if errors.Is(err, pcsdownload.ErrDownloadNotSupportChecksum) {
			return fmt.Errorf("无法校验 MD5: 网盘未提供可用的完整文件 MD5")
		}
		return fmt.Errorf("校验失败: %w", err)
	}

	fmt.Println("校验成功: 文件大小和 MD5 均一致")
	return nil
}
