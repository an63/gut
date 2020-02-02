package yos

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"

	"github.com/1set/gut/ystring"
)

var (
	errInvalidPath    = errors.New("invalid path")
	errSameFile       = errors.New("files are identical")
	errDiffFileSize   = errors.New("different file size")
	errShortRead      = errors.New("short read")
	errIsDirectory    = errors.New("is a directory")
	errNotDirectory   = errors.New("not a directory")
	errNotRegularFile = errors.New("not a regular file")
	errNotSymlink     = errors.New("not a symbolic link")
)

// operation names for the Op field of os.PathError.
var (
	opnCompare = "compare"
	opnCopy    = "copy"
	opnMove    = "move"
	opnList    = "list"
)

// underlyingError returns the underlying error for known os error types. forked from: os/error.go
func underlyingError(err error) error {
	switch err := err.(type) {
	case *os.LinkError:
		return err.Err
	case *os.PathError:
		return err.Err
	case *os.SyscallError:
		return err.Err
	}
	return err
}

// opError returns error struct with given details.
func opError(op, path string, err error) *os.PathError {
	return &os.PathError{
		Op:   op,
		Path: path,
		Err:  underlyingError(err),
	}
}

type (
	funcStatFileInfo  func(name string) (os.FileInfo, error)
	funcCheckFileInfo func(fi *os.FileInfo) bool
	funcRemoveEntry   func(path string) error
	funcCopyEntry     func(src, dest string) error
)

// isFileFi indicates whether the FileInfo is for a regular file.
func isFileFi(fi *os.FileInfo) bool {
	return fi != nil && (*fi).Mode().IsRegular()
}

// isDirFi indicates whether the FileInfo is for a directory.
func isDirFi(fi *os.FileInfo) bool {
	return fi != nil && (*fi).Mode().IsDir()
}

// isSymlinkFi indicates whether the FileInfo is for a symbolic link.
func isSymlinkFi(fi *os.FileInfo) bool {
	return fi != nil && ((*fi).Mode()&os.ModeType == os.ModeSymlink)
}

func isLinkErrorCrossDevice(err error) bool {
	lerr, ok := err.(*os.LinkError)
	return ok && lerr.Err == syscall.EXDEV
}

func isLinkErrorNotDirectory(err error) bool {
	lerr, ok := err.(*os.LinkError)
	return ok && lerr.Err == syscall.ENOTDIR
}

// refineOpPaths validates, cleans up and adjusts the source and destination paths for operations like copy or move.
func refineOpPaths(opName, srcRaw, destRaw string, followLink bool) (src, dest string, err error) {
	if ystring.IsBlank(srcRaw) {
		err = opError(opName, srcRaw, errInvalidPath)
		return
	}
	if ystring.IsBlank(destRaw) {
		err = opError(opName, destRaw, errInvalidPath)
		return
	}

	// clean up paths
	srcRaw, destRaw = filepath.Clean(srcRaw), filepath.Clean(destRaw)

	// use os.Stat to follow symbolic links
	statFunc := os.Lstat
	if followLink {
		statFunc = os.Stat
	}

	// check if source exists
	var srcInfo, destInfo os.FileInfo
	if srcInfo, err = statFunc(srcRaw); err != nil {
		return
	}

	// check if destination exists
	if destInfo, err = statFunc(destRaw); err != nil {
		// check existence of parent of the missing destination
		if os.IsNotExist(err) {
			_, err = os.Stat(filepath.Dir(destRaw))
		}
	} else {
		if os.SameFile(srcInfo, destInfo) {
			err = opError(opName, destRaw, errSameFile)
		} else if destInfo.IsDir() {
			// append file name of source to path of the existing destination
			destRaw = JoinPath(destRaw, srcInfo.Name())
		}
	}

	if err == nil {
		src, dest = srcRaw, destRaw
	}
	return
}

// refineComparePaths validates, cleans up path for file comparison.
func refineComparePaths(pathRaw1, pathRaw2 string, check funcCheckFileInfo, errMode error) (path1, path2 string, err error) {
	if ystring.IsBlank(pathRaw1) {
		err = opError(opnCompare, pathRaw1, errInvalidPath)
		return
	}
	if ystring.IsBlank(pathRaw2) {
		err = opError(opnCompare, pathRaw2, errInvalidPath)
		return
	}

	// clean up paths
	path1, path2 = filepath.Clean(pathRaw1), filepath.Clean(pathRaw2)

	// quit if got no further check
	if check == nil {
		return
	}

	// check file mode of path1
	var fi1, fi2 os.FileInfo
	if fi1, err = os.Stat(path1); err != nil {
		return
	} else if !check(&fi1) {
		err = opError(opnCompare, path1, errMode)
		return
	}

	// check file mode of path2
	if fi2, err = os.Stat(path2); err != nil {
		return
	} else if !check(&fi2) {
		err = opError(opnCompare, path2, errMode)
		return
	}

	// check if it's the identical file and file size
	if os.SameFile(fi1, fi2) {
		err = errSameFile
	} else if isFileFi(&fi1) && fi1.Size() != fi2.Size() {
		err = errDiffFileSize
	}
	return
}
