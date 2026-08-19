package main

import (
	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcscommand"
	"github.com/urfave/cli"
)

func lsFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{Name: "l", Usage: "详细显示"},
		cli.BoolFlag{Name: "t", Usage: "按修改时间排序, 最新的在前"},
		cli.BoolFlag{Name: "S", Usage: "按文件大小排序, 大的在前"},
		cli.BoolFlag{Name: "X", Usage: "按扩展名排序"},
		cli.BoolFlag{Name: "r", Usage: "反转排序结果"},
		cli.BoolFlag{Name: "asc", Usage: "升序排序"},
		cli.BoolFlag{Name: "desc", Usage: "降序排序"},
		cli.BoolFlag{Name: "time", Usage: "根据时间排序"},
		cli.BoolFlag{Name: "name", Usage: "根据文件名排序"},
		cli.BoolFlag{Name: "size", Usage: "根据大小排序"},
	}
}

func resolveLsOptions(c *cli.Context) (*pcscommand.LsOptions, *baidupcs.OrderOptions) {
	orderOptions := &baidupcs.OrderOptions{
		By:    baidupcs.OrderByName,
		Order: baidupcs.OrderAsc,
	}
	sortByExtension := false

	switch {
	case c.Bool("t"):
		orderOptions.By = baidupcs.OrderByTime
		orderOptions.Order = baidupcs.OrderDesc
	case c.Bool("S"):
		orderOptions.By = baidupcs.OrderBySize
		orderOptions.Order = baidupcs.OrderDesc
	case c.Bool("X"):
		sortByExtension = true
	case c.IsSet("time"):
		orderOptions.By = baidupcs.OrderByTime
	case c.IsSet("name"):
		orderOptions.By = baidupcs.OrderByName
	case c.IsSet("size"):
		orderOptions.By = baidupcs.OrderBySize
	}

	if c.Bool("r") {
		orderOptions.Order = reverseOrder(orderOptions.Order)
	}

	// 显式的旧式排序方向参数优先, 保持现有命令兼容。
	switch {
	case c.IsSet("asc"):
		orderOptions.Order = baidupcs.OrderAsc
	case c.IsSet("desc"):
		orderOptions.Order = baidupcs.OrderDesc
	}

	total := c.Bool("l")
	if c.Parent() != nil && c.Parent().Args().Get(0) == "ll" {
		total = true
	}

	return &pcscommand.LsOptions{
		Total:           total,
		SortByExtension: sortByExtension,
	}, orderOptions
}

func reverseOrder(order baidupcs.Order) baidupcs.Order {
	if order == baidupcs.OrderDesc {
		return baidupcs.OrderAsc
	}
	return baidupcs.OrderDesc
}
