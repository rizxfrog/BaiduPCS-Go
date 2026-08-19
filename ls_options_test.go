package main

import (
	"reflect"
	"testing"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
	"github.com/qjfoidnh/BaiduPCS-Go/internal/pcscommand"
	"github.com/urfave/cli"
)

func TestResolveLsOptions(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		wantLsOptions   pcscommand.LsOptions
		wantOrderOption baidupcs.OrderOptions
	}{
		{
			name:            "default name ascending",
			wantLsOptions:   pcscommand.LsOptions{},
			wantOrderOption: baidupcs.OrderOptions{By: baidupcs.OrderByName, Order: baidupcs.OrderAsc},
		},
		{
			name:            "combined long listing and time",
			args:            []string{"-lt"},
			wantLsOptions:   pcscommand.LsOptions{Total: true},
			wantOrderOption: baidupcs.OrderOptions{By: baidupcs.OrderByTime, Order: baidupcs.OrderDesc},
		},
		{
			name:            "combined time reversed",
			args:            []string{"-ltr"},
			wantLsOptions:   pcscommand.LsOptions{Total: true},
			wantOrderOption: baidupcs.OrderOptions{By: baidupcs.OrderByTime, Order: baidupcs.OrderAsc},
		},
		{
			name:            "combined size",
			args:            []string{"-lS"},
			wantLsOptions:   pcscommand.LsOptions{Total: true},
			wantOrderOption: baidupcs.OrderOptions{By: baidupcs.OrderBySize, Order: baidupcs.OrderDesc},
		},
		{
			name:            "combined size reversed",
			args:            []string{"-lSr"},
			wantLsOptions:   pcscommand.LsOptions{Total: true},
			wantOrderOption: baidupcs.OrderOptions{By: baidupcs.OrderBySize, Order: baidupcs.OrderAsc},
		},
		{
			name:          "combined extension",
			args:          []string{"-lX"},
			wantLsOptions: pcscommand.LsOptions{Total: true, SortByExtension: true},
			wantOrderOption: baidupcs.OrderOptions{
				By: baidupcs.OrderByName, Order: baidupcs.OrderAsc,
			},
		},
		{
			name:          "combined extension reversed",
			args:          []string{"-lXr"},
			wantLsOptions: pcscommand.LsOptions{Total: true, SortByExtension: true},
			wantOrderOption: baidupcs.OrderOptions{
				By: baidupcs.OrderByName, Order: baidupcs.OrderDesc,
			},
		},
		{
			name:            "legacy time remains ascending",
			args:            []string{"-time"},
			wantLsOptions:   pcscommand.LsOptions{},
			wantOrderOption: baidupcs.OrderOptions{By: baidupcs.OrderByTime, Order: baidupcs.OrderAsc},
		},
		{
			name:            "legacy size descending",
			args:            []string{"-size", "-desc"},
			wantLsOptions:   pcscommand.LsOptions{},
			wantOrderOption: baidupcs.OrderOptions{By: baidupcs.OrderBySize, Order: baidupcs.OrderDesc},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var gotLsOptions *pcscommand.LsOptions
			var gotOrderOption *baidupcs.OrderOptions
			app := cli.NewApp()
			app.Commands = []cli.Command{
				{
					Name:                   "ls",
					UseShortOptionHandling: true,
					Flags:                  lsFlags(),
					Action: func(c *cli.Context) error {
						gotLsOptions, gotOrderOption = resolveLsOptions(c)
						return nil
					},
				},
			}

			err := app.Run(append([]string{"BaiduPCS-Go", "ls"}, test.args...))
			if err != nil {
				t.Fatalf("parse ls options: %v", err)
			}
			if !reflect.DeepEqual(*gotLsOptions, test.wantLsOptions) {
				t.Fatalf("LsOptions = %+v, want %+v", *gotLsOptions, test.wantLsOptions)
			}
			if !reflect.DeepEqual(*gotOrderOption, test.wantOrderOption) {
				t.Fatalf("OrderOptions = %+v, want %+v", *gotOrderOption, test.wantOrderOption)
			}
		})
	}
}
