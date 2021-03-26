package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// http://stackoverflow.com/questions/13660864/reading-specific-number-of-bytes-from-a-buffered-reader-in-golang

// IsGzipped tests to see if a byte array is gzipped.
func IsGzipped(test []byte) bool {
	// Guard against empty input
	if len(test) == 0 {
		return false
	}
	isGzipped := false
	if test[0] == 31 && test[1] == 139 {
		isGzipped = true
	} else {
		isGzipped = false
	}
	return isGzipped
}

// GzipBytes gzip compress bytes if they are not yet compressed.
func GzipBytes(input []byte) ([]byte, error) {
	if IsGzipped(input) {
		return input, nil
	}
	buf := new(bytes.Buffer)
	gz, err := gzip.NewWriterLevel(buf, gzip.BestCompression)
	if err != nil {
		fmt.Println("got error getting gzip writer", err)
	}
	defer gz.Close()
	if _, err = gz.Write(input); err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	gz.Flush()
	retVal := buf.Bytes()

	return retVal, nil
}

// GunzipBytes gzip byte input. Don't gzip already gzipped input.
func GunzipBytes(input []byte) ([]byte, error) {
	if !IsGzipped(input) {
		return input, nil
	}

	buf := new(bytes.Buffer)
	buf.Write(input)

	gzReader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()
	payload, err := io.ReadAll(gzReader)
	// ReadAll reads until EOF. Reacing EOF is not an error to be worried
	// about according to documentation.
	if err != nil {
		//		panic(err)
	}

	return payload, nil
}
