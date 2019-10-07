package dd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const barChar = string('\U000025A0')

type DataCopier struct {
}

func (d *DataCopier) Copy(from string, to string, offset int64, limit int64) error {
	src, err := os.OpenFile(from, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer src.Close()

	info, err := src.Stat()
	if err != nil {
		return err
	}

	totalLength := info.Size()

	dst, err := os.OpenFile(to, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer dst.Close()

	if offset > 0 {
		if offset > totalLength {
			return fmt.Errorf("the offset value [%d] should not exceed the source file size [%d]", offset, totalLength)
		}
		offset, err = src.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	if limit <= 0 {
		limit = info.Size()
		if offset > 0 {
			limit -= offset
		}
	}

	if (offset + limit) > totalLength {
		return fmt.Errorf("the sum of offset value [%d] and limit value [%d] should not exceed the source file size [%d]", offset, limit, totalLength)
	}

	bufferSize := int64(1024)
	if limit < bufferSize {
		bufferSize = limit
	}
	buffer := make([]byte, bufferSize)

	reader := bufio.NewReader(src)
	writer := bufio.NewWriter(dst)

	length := int64(0)
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		read := int64(n)
		if (length + read) > limit {
			read = limit - length
		}
		length += read

		_, err = writer.Write(buffer[:read])
		if err != nil {
			return err
		}

		percents := float32(100*length) / float32(limit)
		fmt.Printf("\rCopying data: %s [%.2f%%]", strings.Repeat(barChar, int(percents)), percents)

		if length == limit {
			break
		}
	}
	writer.Flush()

	if length > 0 {
		fmt.Println()
	}

	return nil
}
