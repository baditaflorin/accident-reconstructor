package reconstruction

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/baditaflorin/accident-reconstructor/pkg/reconstruct"
)

type ffprobeOutput struct {
	Streams []struct {
		Width        int    `json:"width"`
		Height       int    `json:"height"`
		AvgFrameRate string `json:"avg_frame_rate"`
	} `json:"streams"`
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

// ProbeVideo extracts stable upload metadata using ffprobe when available.
func ProbeVideo(ctx context.Context, path string, fileName string) (reconstruct.UploadInfo, error) {
	info := reconstruct.UploadInfo{FileName: fileName}
	stat, err := os.Stat(path)
	if err != nil {
		return info, fmt.Errorf("stat upload: %w", err)
	}
	info.SizeBytes = stat.Size()
	sum, err := fileSHA256(path)
	if err != nil {
		return info, err
	}
	info.SHA256 = sum

	if _, err := exec.LookPath("ffprobe"); err != nil {
		applyDurationFallback(&info, path)
		return info, nil
	}

	probeCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	// #nosec G204 -- ffprobe is a fixed binary and path is a saved upload.
	out, err := exec.CommandContext(
		probeCtx,
		"ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height,avg_frame_rate",
		"-show_entries", "format=duration",
		"-of", "json",
		path,
	).Output()
	if err != nil {
		applyDurationFallback(&info, path)
		return info, nil
	}

	var parsed ffprobeOutput
	if err := json.Unmarshal(out, &parsed); err != nil {
		return info, fmt.Errorf("parse ffprobe output: %w", err)
	}
	if parsed.Format.Duration != "" {
		info.DurationSeconds, _ = strconv.ParseFloat(parsed.Format.Duration, 64)
		info.DurationSource = "ffprobe"
	}
	if len(parsed.Streams) > 0 {
		stream := parsed.Streams[0]
		info.Width = stream.Width
		info.Height = stream.Height
		info.FrameRate = parseFrameRate(stream.AvgFrameRate)
	}
	if info.DurationSeconds == 0 {
		applyDurationFallback(&info, path)
	}
	return info, nil
}

// applyDurationFallback sets info.DurationSeconds and DurationSource using the
// pure-Go MP4/MOV atom parser when possible, and a clearly-labelled 5-second
// placeholder when not. Speed estimates downstream check DurationSource to
// know whether the input is trustworthy.
func applyDurationFallback(info *reconstruct.UploadInfo, path string) {
	if duration, err := readMP4Duration(path); err == nil && duration > 0 {
		info.DurationSeconds = duration
		info.DurationSource = "mp4_atom"
		return
	}
	info.DurationSeconds = 5
	info.DurationSource = "placeholder"
}

func fileSHA256(path string) (string, error) {
	// #nosec G304 -- path is produced by the upload storage layer.
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open upload for checksum: %w", err)
	}
	defer func() { _ = file.Close() }()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("hash upload: %w", err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func parseFrameRate(value string) float64 {
	if value == "" || value == "0/0" {
		return 0
	}
	parts := strings.Split(value, "/")
	if len(parts) == 1 {
		rate, _ := strconv.ParseFloat(parts[0], 64)
		return rate
	}
	num, _ := strconv.ParseFloat(parts[0], 64)
	den, _ := strconv.ParseFloat(parts[1], 64)
	if den == 0 {
		return 0
	}
	return num / den
}
