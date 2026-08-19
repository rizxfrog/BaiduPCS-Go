package pcsdownload

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qjfoidnh/BaiduPCS-Go/baidupcs"
)

func TestCheckFileMD5(t *testing.T) {
	content := []byte("BaiduPCS-Go checksum test")
	localPath := filepath.Join(t.TempDir(), "test.bin")
	if err := os.WriteFile(localPath, content, 0600); err != nil {
		t.Fatal(err)
	}
	checksum := md5.Sum(content)
	wantMD5 := hex.EncodeToString(checksum[:])

	tests := []struct {
		name       string
		fileInfo   *baidupcs.FileDirectory
		wantErr    error
		wantResult bool
	}{
		{
			name:       "matching standard md5",
			fileInfo:   &baidupcs.FileDirectory{MD5: strings.ToUpper(wantMD5), BlockListJSON: baidupcs.BlockListJSON{BlockList: []string{wantMD5}}},
			wantResult: true,
		},
		{
			name:       "mismatching md5",
			fileInfo:   &baidupcs.FileDirectory{MD5: "00000000000000000000000000000000", BlockListJSON: baidupcs.BlockListJSON{BlockList: []string{"00000000000000000000000000000000"}}},
			wantErr:    ErrDownloadChecksumFailed,
			wantResult: true,
		},
		{
			name:     "multiple blocks",
			fileInfo: &baidupcs.FileDirectory{MD5: wantMD5, BlockListJSON: baidupcs.BlockListJSON{BlockList: []string{wantMD5, wantMD5}}},
			wantErr:  ErrDownloadNotSupportChecksum,
		},
		{
			name:     "invalid remote md5",
			fileInfo: &baidupcs.FileDirectory{MD5: "not-an-md5", BlockListJSON: baidupcs.BlockListJSON{BlockList: []string{"not-an-md5"}}},
			wantErr:  ErrDownloadNotSupportChecksum,
		},
		{
			name:    "nil file info",
			wantErr: ErrDownloadFileInfoNil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CheckFileMD5(localPath, tt.fileInfo)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("CheckFileMD5() error = %v, want %v", err, tt.wantErr)
			}
			if (result != nil) != tt.wantResult {
				t.Fatalf("CheckFileMD5() result = %#v, wantResult %v", result, tt.wantResult)
			}
			if result != nil && result.LocalMD5 != wantMD5 {
				t.Errorf("local MD5 = %s, want %s", result.LocalMD5, wantMD5)
			}
		})
	}
}
