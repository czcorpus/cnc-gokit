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
	"os"
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
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	finfo, err := f.Stat()
	if err != nil {
		return false, err
	}
	return finfo.Mode().IsDir(), nil
}
