package reconstruction

import (
	"fmt"
	"strings"

	"github.com/baditaflorin/accident-reconstructor/pkg/reconstruct"
)

func BuildReport(artifact reconstruct.Artifact) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Accident Reconstruction Report\n\n")
	fmt.Fprintf(&b, "Case: %s\n\n", artifact.CaseName)
	fmt.Fprintf(&b, "Case ID: %s\n\n", artifact.CaseID)
	fmt.Fprintf(&b, "Generated: %s\n\n", artifact.CreatedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(&b, "Version: %s (%s)\n\n", artifact.Version, artifact.Commit)
	fmt.Fprintf(&b, "## Inputs\n\n")
	for _, upload := range artifact.Uploads {
		fmt.Fprintf(
			&b,
			"- %s: %.2fs, %dx%d, %.2f fps, sha256 %s\n",
			upload.FileName,
			upload.DurationSeconds,
			upload.Width,
			upload.Height,
			upload.FrameRate,
			upload.SHA256,
		)
	}
	fmt.Fprintf(&b, "\n## Speed Estimate\n\n")
	fmt.Fprintf(&b, "- Method: %s\n", artifact.Speed.Method)
	fmt.Fprintf(&b, "- Mean: %.2f km/h (%.2f mph)\n", artifact.Speed.MeanKPH, artifact.Speed.MeanMPH)
	fmt.Fprintf(&b, "- Range: %.2f-%.2f km/h\n", artifact.Speed.LowerKPH, artifact.Speed.UpperKPH)
	fmt.Fprintf(&b, "- Confidence: %.2f\n\n", artifact.Speed.Confidence)
	fmt.Fprintf(&b, "## Reconstruction Quality\n\n")
	fmt.Fprintf(&b, "- Mode: %s\n", artifact.Quality.Mode)
	fmt.Fprintf(&b, "- Coordinate frame: %s\n", artifact.Quality.CoordinateFrame)
	fmt.Fprintf(&b, "- Sparse points: %d\n", len(artifact.Points))
	fmt.Fprintf(&b, "- Camera poses: %d\n\n", len(artifact.Cameras))
	fmt.Fprintf(&b, "## Toolchain\n\n")
	for _, tool := range artifact.Quality.Toolchain {
		if tool.Path == "" {
			fmt.Fprintf(&b, "- %s: %s\n", tool.Name, tool.Status)
			continue
		}
		fmt.Fprintf(&b, "- %s: %s (%s)\n", tool.Name, tool.Status, tool.Version)
	}
	fmt.Fprintf(&b, "\n## Limitations\n\n")
	for _, warning := range artifact.Quality.Warnings {
		fmt.Fprintf(&b, "- %s\n", warning)
	}
	fmt.Fprintf(&b, "\nThis report is an evidence organization aid and does not certify legal admissibility.\n")
	return b.String()
}
