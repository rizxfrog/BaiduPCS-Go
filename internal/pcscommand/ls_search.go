package pcscommand

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/pcstable"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/converter"
	"github.com/qjfoidnh/BaiduPCS-Go/pcsutil/pcstime"
)

type (
	// LsOptions 列目录可选项
	LsOptions struct {
		Total           bool
		SortByExtension bool
	}

	// SearchOptions 搜索可选项
	SearchOptions struct {
		Total   bool
		Recurse bool
	}
)

const (
	opLs int = iota
	opSearch
)

// RunLs 执行列目录
func RunLs(pcspath string, lsOptions *LsOptions, orderOptions *baidupcs.OrderOptions) {
	err := matchPathByShellPatternOnce(&pcspath)
	if err != nil {
		fmt.Println(err)
		return
	}

	files, err := GetBaiduPCS().FilesDirectoriesList(pcspath, orderOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	if lsOptions == nil {
		lsOptions = &LsOptions{}
	}
	if lsOptions.SortByExtension {
		order := baidupcs.OrderAsc
		if orderOptions != nil {
			order = orderOptions.Order
		}
		sortFileDirectoryListByExtension(files, order)
	}

	fmt.Printf("\n当前目录: %s\n----\n", pcspath)

	renderTable(opLs, lsOptions.Total, pcspath, files)
	return
}

func sortFileDirectoryListByExtension(files baidupcs.FileDirectoryList, order baidupcs.Order) {
	sort.SliceStable(files, func(i, j int) bool {
		left, right := files[i], files[j]
		if left == nil || right == nil {
			return left != nil
		}
		if left.Isdir != right.Isdir {
			return left.Isdir
		}

		leftName := strings.ToLower(left.Filename)
		rightName := strings.ToLower(right.Filename)
		leftExtension := ""
		rightExtension := ""
		if !left.Isdir {
			leftExtension = strings.ToLower(path.Ext(left.Filename))
			rightExtension = strings.ToLower(path.Ext(right.Filename))
		}

		comparison := strings.Compare(leftExtension, rightExtension)
		if comparison == 0 {
			comparison = strings.Compare(leftName, rightName)
		}
		if order == baidupcs.OrderDesc {
			return comparison > 0
		}
		return comparison < 0
	})
}

// RunSearch 执行搜索
func RunSearch(targetPath, keyword string, opt *SearchOptions) {
	err := matchPathByShellPatternOnce(&targetPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	if opt == nil {
		opt = &SearchOptions{}
	}

	files, err := GetBaiduPCS().Search(targetPath, keyword, opt.Recurse)
	if err != nil {
		fmt.Println(err)
		return
	}

	renderTable(opSearch, opt.Total, targetPath, files)
	return
}

func renderTable(op int, isTotal bool, path string, files baidupcs.FileDirectoryList) {
	tb := pcstable.NewTable(os.Stdout)
	var (
		fN, dN   int64
		showPath string
	)

	switch op {
	case opLs:
		showPath = "文件(目录)"
	case opSearch:
		showPath = "路径"
	}

	if isTotal {
		tb.SetHeader([]string{"#", "fs_id", "app_id", "文件大小", "创建日期", "修改日期", "md5(截图请打码)", showPath})
		tb.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
		for k, file := range files {
			if file.Isdir {
				tb.Append([]string{strconv.Itoa(k), strconv.FormatInt(file.FsID, 10), strconv.FormatInt(file.AppID, 10), "-", pcstime.FormatTime(file.Ctime), pcstime.FormatTime(file.Mtime), file.MD5, file.Filename + baidupcs.PathSeparator})
				continue
			}

			var md5 string
			if len(file.BlockList) > 1 {
				md5 = "(可能不正确)" + file.MD5
			} else {
				md5 = file.MD5
			}

			switch op {
			case opLs:
				tb.Append([]string{strconv.Itoa(k), strconv.FormatInt(file.FsID, 10), strconv.FormatInt(file.AppID, 10), converter.ConvertFileSize(file.Size, 2), pcstime.FormatTime(file.Ctime), pcstime.FormatTime(file.Mtime), md5, file.Filename})
			case opSearch:
				tb.Append([]string{strconv.Itoa(k), strconv.FormatInt(file.FsID, 10), strconv.FormatInt(file.AppID, 10), converter.ConvertFileSize(file.Size, 2), pcstime.FormatTime(file.Ctime), pcstime.FormatTime(file.Mtime), md5, file.Path})
			}
		}
		fN, dN = files.Count()
		tb.Append([]string{"", "", "总: " + converter.ConvertFileSize(files.TotalSize(), 2), "", "", "", fmt.Sprintf("文件总数: %d, 目录总数: %d", fN, dN)})
	} else {
		tb.SetHeader([]string{"#", "文件大小", "修改日期", showPath})
		tb.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
		for k, file := range files {
			if file.Isdir {
				tb.Append([]string{strconv.Itoa(k), "-", pcstime.FormatTime(file.Mtime), file.Filename + baidupcs.PathSeparator})
				continue
			}

			switch op {
			case opLs:
				tb.Append([]string{strconv.Itoa(k), converter.ConvertFileSize(file.Size, 2), pcstime.FormatTime(file.Mtime), file.Filename})
			case opSearch:
				tb.Append([]string{strconv.Itoa(k), converter.ConvertFileSize(file.Size, 2), pcstime.FormatTime(file.Mtime), file.Path})
			}
		}
		fN, dN = files.Count()
		tb.Append([]string{"", "总: " + converter.ConvertFileSize(files.TotalSize(), 2), "", fmt.Sprintf("文件总数: %d, 目录总数: %d", fN, dN)})
	}

	tb.Render()

	if fN+dN >= 50 {
		fmt.Printf("\n当前目录: %s\n", path)
	}

	fmt.Printf("----\n")
}
