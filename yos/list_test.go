package yos

import (
	"os"
	"strings"
	"testing"
)

var TestCaseRootList string

func init() {
	TestCaseRootList = JoinPath(os.Getenv("TESTRSSDIR"), "yos", "list")
	//TestCaseRootList = "/Users/vej/go/src/github.com/1set/gut/local/test_resource/yos/list"
}

func TestListAll(t *testing.T) {
	content := []string{
		"yos/list",
		"yos/list/File0.txt",
		"yos/list/File4.txt",
		"yos/list/broken_symlink.wtf",
		"yos/list/deep_folder",
		"yos/list/deep_folder/deep",
		"yos/list/deep_folder/deep/deeper",
		"yos/list/deep_folder/deep/deeper/deepest",
		"yos/list/deep_folder/deep/deeper/deepest/text_file.txt",
		"yos/list/empty_folder",
		"yos/list/file1.txt",
		"yos/list/file2.txt",
		"yos/list/file3.txt",
		"yos/list/folder_like_file.txt",
		"yos/list/nested_empty",
		"yos/list/nested_empty/empty1",
		"yos/list/nested_empty/empty1/empty2",
		"yos/list/nested_empty/empty1/empty2/empty3",
		"yos/list/nested_empty/empty1/empty2/empty3/empty4",
		"yos/list/nested_empty/empty1/empty2/empty3/empty4/empty5",
		"yos/list/no_ext_name_file",
		"yos/list/simple_folder",
		"yos/list/simple_folder/file1.txt",
		"yos/list/simple_folder/file2.txt",
		"yos/list/simple_folder/file3.txt",
		"yos/list/symlink_to_dir.txt",
		"yos/list/symlink_to_file.txt",
		"yos/list/white space",
		"yos/list/white space/only one file",
		"yos/list/white space.txt",
		"yos/list/测试文件.md",
		"yos/list/🤙🏝️.md",
	}

	items, err := ListAll(TestCaseRootList)
	if err != nil {
		t.Errorf("ListAll() got error = %v, wantErr %v", err, false)
		return
	}
	if len(items) != len(content) {
		t.Errorf("ListAll() got length = %v, want = %v", len(items), len(content))
		return
	}

	for idx, item := range items {
		suffix := content[idx]
		if !strings.HasSuffix(item.Path, suffix) {
			t.Errorf("ListAll() got #%d path = %q, want suffix = %q", idx, item.Path, suffix)
			return
		}
		fileName := (*item.Info).Name()
		if !strings.HasSuffix(suffix, fileName) {
			t.Errorf("ListAll() got #%d suffix = %q, want name = %q", idx, suffix, fileName)
			return
		}
	}
}
