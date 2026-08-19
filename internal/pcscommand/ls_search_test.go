package pcscommand

import (
	"reflect"
	"testing"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
)

func TestSortFileDirectoryListByExtension(t *testing.T) {
	newList := func() baidupcs.FileDirectoryList {
		return baidupcs.FileDirectoryList{
			{Filename: "z-dir", Isdir: true},
			{Filename: "b.txt"},
			{Filename: "README"},
			{Filename: "a-dir", Isdir: true},
			{Filename: "a.go"},
			{Filename: "Z.TXT"},
		}
	}
	filenames := func(files baidupcs.FileDirectoryList) []string {
		result := make([]string, 0, len(files))
		for _, file := range files {
			result = append(result, file.Filename)
		}
		return result
	}

	tests := []struct {
		name  string
		order baidupcs.Order
		want  []string
	}{
		{
			name:  "ascending",
			order: baidupcs.OrderAsc,
			want:  []string{"a-dir", "z-dir", "README", "a.go", "b.txt", "Z.TXT"},
		},
		{
			name:  "descending",
			order: baidupcs.OrderDesc,
			want:  []string{"z-dir", "a-dir", "Z.TXT", "b.txt", "a.go", "README"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			files := newList()
			sortFileDirectoryListByExtension(files, test.order)
			if got := filenames(files); !reflect.DeepEqual(got, test.want) {
				t.Fatalf("filenames = %v, want %v", got, test.want)
			}
		})
	}
}
