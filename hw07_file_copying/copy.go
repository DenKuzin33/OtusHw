package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

var (
	ErrSamePath              = errors.New("fromPath equals toPath")
	ErrFromPath              = errors.New("fromPath doesn't exist")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type progressBar struct {
	TotalBytesCount   int
	CurrentBytesCount int
}

func (pb *progressBar) Write(p []byte) (n int, err error) {
	pLen := len(p)
	pb.CurrentBytesCount += pLen
	fmt.Printf("%d%%\n", (pb.CurrentBytesCount*100)/pb.TotalBytesCount)
	return pLen, nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == toPath {
		return ErrSamePath
	}

	source, err := os.OpenFile(fromPath, os.O_RDONLY, 0444)
	if err != nil {
		if err == os.ErrNotExist {
			return ErrFromPath
		}
		return err
	}
	defer source.Close()

	sourceInfo, err := source.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}
	sourceLen := sourceInfo.Size()
	bytesToCopy := sourceLen - offset
	if limit != 0 && limit < bytesToCopy {
		bytesToCopy = limit
	}

	if offset > sourceLen {
		return ErrOffsetExceedsFileSize
	}

	target, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer target.Close()

	// Подбираем подходящий под количество байт размер буфера
	var bufferSize int
	if bytesToCopy < 100 {
		bufferSize = 4
	} else {
		bufferSize = int(math.Pow(2, math.Ceil(math.Log(float64(bytesToCopy/100))/math.Log(2))))
	}
	buffer := make([]byte, bufferSize)
	pb := progressBar{TotalBytesCount: int(bytesToCopy)}
	multiWriter := io.MultiWriter(target, &pb)

	if offset > 0 {
		_, err := source.Seek(offset, 0)
		if err != nil {
			return err
		}
	}
	for pb.CurrentBytesCount < int(bytesToCopy) {
		_, err := source.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		remained := bytesToCopy - int64(pb.CurrentBytesCount)
		if remained < int64(bufferSize) {
			multiWriter.Write(buffer[:remained])
		} else {
			multiWriter.Write(buffer)
		}
	}
	return nil
}
