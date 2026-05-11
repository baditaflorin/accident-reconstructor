package reconstruction

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

// errNoMP4Duration is returned when the input is not an MP4/MOV ISO BMFF file
// or does not contain a usable mvhd atom. Callers should treat this as
// "metadata unavailable" rather than a hard failure.
var errNoMP4Duration = errors.New("no mp4 duration available")

// readMP4Duration extracts duration from the ISO Base Media File Format
// (MP4, MOV, M4V) mvhd atom inside the moov box. Returns seconds.
//
// This is the pure-Go fallback used when ffprobe is not installed. It is
// intentionally minimal: it walks top-level boxes, descends into moov, finds
// mvhd, and reads timescale/duration. No external tools, no decoding.
func readMP4Duration(path string) (float64, error) {
	f, err := os.Open(path) // #nosec G304 -- caller passes a path from our upload dir.
	if err != nil {
		return 0, fmt.Errorf("open: %w", err)
	}
	defer func() { _ = f.Close() }()

	stat, err := f.Stat()
	if err != nil {
		return 0, fmt.Errorf("stat: %w", err)
	}
	fileSize := stat.Size()

	moovOffset, moovSize, err := findBox(f, 0, fileSize, "moov")
	if err != nil {
		return 0, err
	}
	mvhdOffset, mvhdSize, err := findBox(f, moovOffset, moovOffset+moovSize, "mvhd")
	if err != nil {
		return 0, err
	}
	return readMvhdDuration(f, mvhdOffset, mvhdSize)
}

// findBox scans box headers starting at `start` and returns the first child
// whose four-cc type matches `wantType`. `end` is the exclusive byte limit.
func findBox(r io.ReaderAt, start, end int64, wantType string) (int64, int64, error) {
	pos := start
	for pos+8 <= end {
		var header [8]byte
		if _, err := r.ReadAt(header[:], pos); err != nil {
			return 0, 0, fmt.Errorf("read box header: %w", err)
		}
		size := int64(binary.BigEndian.Uint32(header[0:4]))
		boxType := string(header[4:8])
		headerSize := int64(8)
		switch size {
		case 0:
			// Box extends to EOF.
			size = end - pos
		case 1:
			// Extended size in next 8 bytes.
			var ext [8]byte
			if _, err := r.ReadAt(ext[:], pos+8); err != nil {
				return 0, 0, fmt.Errorf("read extended size: %w", err)
			}
			raw := binary.BigEndian.Uint64(ext[:])
			// Reject 64-bit sizes that won't fit in int64 — they can't be valid
			// container offsets anyway, and the conversion would otherwise wrap.
			if raw > uint64(1<<62) {
				return 0, 0, fmt.Errorf("%w: extended size %d out of range", errNoMP4Duration, raw)
			}
			size = int64(raw)
			headerSize = 16
		}
		if size < headerSize || pos+size > end {
			return 0, 0, fmt.Errorf("box %q has invalid size %d at offset %d", boxType, size, pos)
		}
		if boxType == wantType {
			return pos + headerSize, size - headerSize, nil
		}
		pos += size
	}
	return 0, 0, fmt.Errorf("%w: %q not found", errNoMP4Duration, wantType)
}

func readMvhdDuration(r io.ReaderAt, offset, size int64) (float64, error) {
	if size < 4 {
		return 0, fmt.Errorf("%w: mvhd too small", errNoMP4Duration)
	}
	var versionAndFlags [4]byte
	if _, err := r.ReadAt(versionAndFlags[:], offset); err != nil {
		return 0, fmt.Errorf("read mvhd version: %w", err)
	}
	version := versionAndFlags[0]

	// mvhd layout (skipping version/flags):
	//   v0: creation(4) modification(4) timescale(4) duration(4) ...
	//   v1: creation(8) modification(8) timescale(4) duration(8) ...
	switch version {
	case 0:
		if size < 20 {
			return 0, fmt.Errorf("%w: mvhd v0 truncated", errNoMP4Duration)
		}
		var buf [8]byte
		if _, err := r.ReadAt(buf[:], offset+12); err != nil {
			return 0, fmt.Errorf("read mvhd v0 body: %w", err)
		}
		timescale := binary.BigEndian.Uint32(buf[0:4])
		duration := binary.BigEndian.Uint32(buf[4:8])
		if timescale == 0 {
			return 0, fmt.Errorf("%w: timescale is zero", errNoMP4Duration)
		}
		return float64(duration) / float64(timescale), nil
	case 1:
		if size < 32 {
			return 0, fmt.Errorf("%w: mvhd v1 truncated", errNoMP4Duration)
		}
		var buf [12]byte
		if _, err := r.ReadAt(buf[:], offset+20); err != nil {
			return 0, fmt.Errorf("read mvhd v1 body: %w", err)
		}
		timescale := binary.BigEndian.Uint32(buf[0:4])
		duration := binary.BigEndian.Uint64(buf[4:12])
		if timescale == 0 {
			return 0, fmt.Errorf("%w: timescale is zero", errNoMP4Duration)
		}
		return float64(duration) / float64(timescale), nil
	default:
		return 0, fmt.Errorf("%w: unknown mvhd version %d", errNoMP4Duration, version)
	}
}
