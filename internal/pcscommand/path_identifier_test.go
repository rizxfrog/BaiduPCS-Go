package pcscommand

import (
	"strings"
	"testing"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
)

func TestResolvePathIdentifier(t *testing.T) {
	const (
		currentDirectory = "/apps/bypy"
		directoryMD5     = "0123456789abcdef0123456789abcdef"
	)
	files := baidupcs.FileDirectoryList{
		{FsID: 101, Path: "/apps/bypy/course", Filename: "course", Isdir: true, MD5: directoryMD5},
		{FsID: 102, Path: "/apps/bypy/video.mp4", Filename: "video.mp4", MD5: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	}

	tests := []struct {
		name        string
		target      string
		files       baidupcs.FileDirectoryList
		want        string
		wantErrText string
	}{
		{
			name:   "ordinary path is unchanged",
			target: "course",
			files:  files,
			want:   "course",
		},
		{
			name:   "explicit relative path can address literal prefix",
			target: "./$101",
			files:  files,
			want:   "./$101",
		},
		{
			name:   "resolve directory by fs id",
			target: "$101",
			files:  files,
			want:   "/apps/bypy/course",
		},
		{
			name:   "resolve file by fs id",
			target: "$102",
			files:  files,
			want:   "/apps/bypy/video.mp4",
		},
		{
			name:   "resolve directory md5 case insensitively",
			target: "%0123456789ABCDEF0123456789ABCDEF",
			files:  files,
			want:   "/apps/bypy/course",
		},
		{
			name:   "resolve file by md5",
			target: "%aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			files:  files,
			want:   "/apps/bypy/video.mp4",
		},
		{
			name:   "build path when response path is empty",
			target: "$103",
			files: baidupcs.FileDirectoryList{
				{FsID: 103, Filename: "fallback", Isdir: true},
			},
			want: "/apps/bypy/fallback",
		},
		{
			name:        "reject invalid fs id",
			target:      "$not-a-number",
			files:       files,
			wantErrText: "无效的 fs_id",
		},
		{
			name:        "reject invalid md5",
			target:      "%1234",
			files:       files,
			wantErrText: "32 位十六进制字符串",
		},
		{
			name:        "identifier not found",
			target:      "$999",
			files:       files,
			wantErrText: "当前目录 /apps/bypy 中未找到 fs_id=999",
		},
		{
			name:   "duplicate md5 is ambiguous",
			target: "%0123456789abcdef0123456789abcdef",
			files: append(files,
				&baidupcs.FileDirectory{FsID: 104, Path: "/apps/bypy/duplicate", Filename: "duplicate", Isdir: true, MD5: directoryMD5},
			),
			wantErrText: "有 2 个项目匹配 MD5",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := resolvePathIdentifier(test.target, currentDirectory, test.files)
			if test.wantErrText != "" {
				if err == nil || !strings.Contains(err.Error(), test.wantErrText) {
					t.Fatalf("error = %v, want containing %q", err, test.wantErrText)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolve path identifier: %v", err)
			}
			if got != test.want {
				t.Fatalf("target = %q, want %q", got, test.want)
			}
		})
	}
}
