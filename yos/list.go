package yos

import (
	"os"
	"path/filepath"
	"strings"
)

// A FilePathInfo describes a file's path and stat.
type FilePathInfo struct {
	Path string
	Info os.FileInfo
}

// listCondEntries returns a list of conditional directory entries.
func listCondEntries(root string, cond func(os.FileInfo) (bool, error)) (entries []*FilePathInfo, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if root == path {
			return nil
		}
		var ok bool
		if ok, err = cond(info); ok {
			entries = append(entries, &FilePathInfo{
				Path: path,
				Info: info,
			})
		}
		return err
	})
	return
}

// ListAll returns a list of all directory entries in the given directory in lexical order.
// It searches recursively, but symbolic links will be not be followed.
func ListAll(root string) (entries []*FilePathInfo, err error) {
	return listCondEntries(root, func(info os.FileInfo) (bool, error) { return true, nil })
}

// ListFile returns a list of file directory entries in the given directory in lexical order.
// It searches recursively, but symbolic links will be not be followed.
func ListFile(root string) (entries []*FilePathInfo, err error) {
	return listCondEntries(root, func(info os.FileInfo) (bool, error) { return !info.IsDir(), nil })
}

// ListDir returns a list of nested directory entries in the given directory in lexical order.
// It searches recursively, but symbolic links will be not be followed.
func ListDir(root string) (entries []*FilePathInfo, err error) {
	return listCondEntries(root, func(info os.FileInfo) (bool, error) { return info.IsDir(), nil })
}

// The flags are used by the ListMatch methods.
const (
	// ListRecursive indicates ListMatch to recursively list directory entries encountered.
	ListRecursive int = 1 << iota
	// ListRecursive indicates ListMatch to convert file name to lower case before the pattern matching.
	ListToLower
	// ListRecursive indicates ListMatch to include matched files in the returned list.
	ListIncludeFile
	// ListRecursive indicates ListMatch to include matched directories in the returned list.
	ListIncludeDir
)

// ListMatch returns a list of directory entries that matches any given pattern in the directory in lexical order.
// ListMatch requires the pattern to match all of the filename, not just a substring.
// Symbolic links will be not be followed. ErrBadPattern is returned if any pattern is malformed.
func ListMatch(root string, flag int, patterns ...string) (entries []*FilePathInfo, err error) {
	return listCondEntries(root, func(info os.FileInfo) (ok bool, err error) {
		fileName := info.Name()
		if flag&ListToLower != 0 {
			fileName = strings.ToLower(fileName)
		}
		isDir := info.IsDir()
		if (isDir && (flag&ListIncludeDir != 0)) || (!isDir && (flag&ListIncludeFile != 0)) {
			for _, pattern := range patterns {
				ok, err = filepath.Match(pattern, fileName)
				if ok || err != nil {
					break
				}
			}
		}
		if err == nil && isDir && (flag&ListRecursive == 0) {
			err = filepath.SkipDir
		}
		return
	})
}
