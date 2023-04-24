// Copyright 2023 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2023 Martin Zimandl <martin.zimandl@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

// IsFile tests whether the file with a specified
// path is a regular file.
func IsFile(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, nil
	}
	finfo, err := f.Stat()
	if err != nil {
		return false, err
	}
	return finfo.Mode().IsRegular(), nil
}

// DeleteFile deletes a regular file. If not found
// or any other occurs, error is returned.
func DeleteFile(path string) error {
	isFile, err := IsFile(path)
	if err != nil {
		return fmt.Errorf("failed to delete file %s: %w", path, err)
	}
	if !isFile {
		return fmt.Errorf("failed to delete file %s: path is not a file", path)
	}
	return os.Remove(path)
}

// GetFileMtime returns a file modification time
func GetFileMtime(filePath string) (time.Time, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return time.Time{}, err
	}
	finfo, err := f.Stat()
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(finfo.ModTime().Unix(), 0), nil
}

// IsDir tests whether a provided path represents
// a directory. If not or in case of an IO error,
// false is returned along with the error
func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

// PathExists tests whether  the provided path exists no matter what it
// is (file, dir, ...)
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// FileSize returns size of a provided file.
// In case of an error the function returns -1 and logs the error.
func FileSize(path string) (int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	finfo, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return finfo.Size(), nil
}

// FileList is an abstraction for list of files along with their
// modification time information. It supports sorting.
type FileList struct {
	files []os.FileInfo
}

func (f *FileList) Len() int {
	return len(f.files)
}

func (f *FileList) Less(i, j int) bool {
	return f.files[i].ModTime().After(f.files[j].ModTime())
}

func (f *FileList) Swap(i, j int) {
	f.files[i], f.files[j] = f.files[j], f.files[i]
}

// First returns an item with the latest modification time.
func (f *FileList) First() os.FileInfo {
	return f.files[0]
}

func (f *FileList) ForEach(fn func(info os.FileInfo, idx int) bool) {
	for i, v := range f.files {
		if !fn(v, i) {
			break
		}
	}
}

// ListFilesInDir lists files according to their modification time
// (newest first).
func ListFilesInDir(path string, newestFirst bool) (FileList, error) {
	var ans FileList
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return ans, err
	}
	ans.files = make([]os.FileInfo, len(files))
	copy(ans.files, files)
	if newestFirst {
		sort.Sort(&ans)
	}
	return ans, nil
}

// ListDirsInDir lists directories in a directory (without recursion).
func ListDirsInDir(path string, newestFirst bool) (FileList, error) {
	var ans FileList
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return ans, err
	}
	ans.files = make([]os.FileInfo, 0, 200)
	for _, v := range files {
		if v.IsDir() {
			ans.files = append(ans.files, v)
		}
	}
	return ans, nil
}
