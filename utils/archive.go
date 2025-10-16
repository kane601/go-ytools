package tools

import (
	"archive/zip"
	"context"
	"os"
	"path/filepath"
	"strings"
)

// CompressPath 压缩
func CompressPath(ctx context.Context, srcFile string, destZip string) error {
	return CompressPathFun(ctx, srcFile, destZip, nil)
}

// DeCompress 解压
func DeCompress(ctx context.Context, zipFile, dest string) error {
	return DeCompressFun(ctx, zipFile, dest, nil)
}

// DecompressSize 解压大小
func DecompressSize(zipFile string) int64 {
	zipFile = AbsPath(zipFile)
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return 0
	}
	defer reader.Close()

	decompressSize := int64(0)
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		decompressSize += int64(file.UncompressedSize64)

	}
	return decompressSize
}

// DeCompressFun 解压
func DeCompressFun(ctx context.Context, zipFile, dest string, progHand ProgressHand) error {
	zipFile = AbsPath(zipFile)
	dest = AbsPath(dest)
	mutilCopy := NewMutilCopyHander(ctx, DecompressSize(zipFile), progHand)

	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := AbsJoinPath(dest, file.Name)

		if file.FileInfo().IsDir() {
			err := CreateDirs(filename)
			if err != nil {
				return err
			}
			continue
		}

		w, err := CreateFile(filename, file.Mode())
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = mutilCopy.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

// CompressPathFun 压缩
func CompressPathFun(ctx context.Context, srcFile, destZip string, progHand ProgressHand) error {
	srcFile = AbsPath(srcFile)
	destZip = AbsPath(destZip)
	mutilCopy := NewMutilCopyHander(ctx, PathSize(srcFile), progHand)

	zipfile, err := CreateFile(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile))
		header.Name = strings.TrimPrefix(header.Name, "/")
		header.Name = strings.TrimPrefix(header.Name, "\\")
		// header.Name = path
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := OpenReadFile(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = mutilCopy.Copy(writer, file)
		}
		return err
	})
	return err
}
