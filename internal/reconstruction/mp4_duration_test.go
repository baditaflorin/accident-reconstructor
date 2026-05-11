package reconstruction

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// buildSyntheticMP4 writes a tiny ISO BMFF stream containing only the boxes
// readMP4Duration needs: a placeholder ftyp, a moov, and an mvhd. There is no
// actual A/V data — just the metadata atoms.
func buildSyntheticMP4(t *testing.T, timescale uint32, duration uint32) []byte {
	t.Helper()

	mvhd := make([]byte, 108)
	// version=0, flags=0
	mvhd[0], mvhd[1], mvhd[2], mvhd[3] = 0, 0, 0, 0
	// creation_time(4), modification_time(4)
	binary.BigEndian.PutUint32(mvhd[4:8], 0)
	binary.BigEndian.PutUint32(mvhd[8:12], 0)
	// timescale(4), duration(4)
	binary.BigEndian.PutUint32(mvhd[12:16], timescale)
	binary.BigEndian.PutUint32(mvhd[16:20], duration)
	// remaining fields (rate, volume, matrix, etc.) are zero — readMvhdDuration
	// only reads through offset 20, so the rest is just padding.

	mvhdBox := encodeBox("mvhd", mvhd)
	moovBox := encodeBox("moov", mvhdBox)
	ftypBody := []byte("isom\x00\x00\x02\x00mp41avc1")
	ftypBox := encodeBox("ftyp", ftypBody)

	out := append([]byte{}, ftypBox...)
	out = append(out, moovBox...)
	return out
}

func encodeBox(boxType string, body []byte) []byte {
	total := 8 + len(body)
	if total < 0 || total > int(^uint32(0)) {
		panic("encodeBox: synthetic box body too large for test")
	}
	header := make([]byte, 8)
	binary.BigEndian.PutUint32(header[0:4], uint32(total)) // #nosec G115 -- bounded above.
	copy(header[4:8], boxType)
	return append(header, body...)
}

func TestReadMP4Duration_VanillaV0(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "synthetic.mp4")
	// timescale = 1000 (ms), duration = 7500 ticks → 7.5 seconds
	if err := os.WriteFile(tmp, buildSyntheticMP4(t, 1000, 7500), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	got, err := readMP4Duration(tmp)
	if err != nil {
		t.Fatalf("readMP4Duration: %v", err)
	}
	if got < 7.49 || got > 7.51 {
		t.Fatalf("duration = %.3f, want ~7.5", got)
	}
}

func TestReadMP4Duration_NotAnMP4(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "not_mp4.bin")
	if err := os.WriteFile(tmp, []byte("this is not a video"), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	_, err := readMP4Duration(tmp)
	if err == nil {
		t.Fatal("expected error for non-MP4 input")
	}
	if !errors.Is(err, errNoMP4Duration) {
		// Garbage input collapses to either "moov not found" (errNoMP4Duration)
		// or a header-read error; either way, the call must not succeed.
		t.Logf("non-MP4 returned %v (acceptable as long as it's an error)", err)
	}
}
