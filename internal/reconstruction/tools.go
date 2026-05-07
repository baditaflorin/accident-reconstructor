package reconstruction

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"github.com/baditaflorin/accident-reconstructor/pkg/reconstruct"
)

// DiscoverTools reports external reconstruction tool availability.
func DiscoverTools(ctx context.Context) []reconstruct.ToolStatus {
	names := []string{"colmap", "ffmpeg", "ffprobe", "gdalinfo", "python3", "ollama"}
	result := make([]reconstruct.ToolStatus, 0, len(names))
	for _, name := range names {
		path, err := exec.LookPath(name)
		if err != nil {
			result = append(result, reconstruct.ToolStatus{Name: name, Status: "missing"})
			continue
		}
		result = append(result, reconstruct.ToolStatus{
			Name:    name,
			Path:    path,
			Version: toolVersion(ctx, path),
			Status:  "available",
		})
	}
	return result
}

func toolAvailable(tools []reconstruct.ToolStatus, name string) bool {
	for _, tool := range tools {
		if tool.Name == name && tool.Status == "available" {
			return true
		}
	}
	return false
}

func toolVersion(ctx context.Context, path string) string {
	versionCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	out, err := exec.CommandContext(versionCtx, path, "--version").CombinedOutput()
	if err != nil {
		return ""
	}
	line := strings.TrimSpace(strings.Split(string(out), "\n")[0])
	if len(line) > 120 {
		return line[:120]
	}
	return line
}
