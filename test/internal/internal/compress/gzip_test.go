package compress

// $Id$

import (
	"fmt"
	"testing"
)

func TestGzip(t *testing.T) {
	b := []byte("Hello there, how are you today? I am fine, thank-you very much.")
	t.Log(fmt.Sprintf("buffer length: %d", len(b)))
	gzipped, err := GzipBytes(b)
	if err != nil {
		t.Logf("Got error compressing content %v", err)
	}

	t.Log(fmt.Sprintf("gzipped length: %d", len(gzipped)))

	ungzipped, err := GunzipBytes(gzipped)
	if err != nil {
		t.Logf("Got error compressing %v", err)
	}
	t.Log(fmt.Sprintf("gzipped %v (should be false)", IsGzipped([]byte("hello"))))
	t.Log(fmt.Sprintf("gzipped %v (should be true)", IsGzipped(gzipped)))
	t.Log(fmt.Sprintf("gzipped %v (should be false)", IsGzipped(ungzipped)))
	t.Log("Gunzipped content is: '" + string(ungzipped) + "'")
	t.Log(fmt.Sprintf("Pre and post gzip string length the same: %v", (len(b) == len(ungzipped))))
}
