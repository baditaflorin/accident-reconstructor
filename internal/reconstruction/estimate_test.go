package reconstruction

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/baditaflorin/accident-reconstructor/pkg/reconstruct"
)

func TestEstimateSceneProducesTrackAndSpeed(t *testing.T) {
	artifact := EstimateScene(EstimateInput{
		CaseID:      "case-1",
		CaseName:    "test case",
		Version:     "0.1.0",
		Commit:      "test",
		ScaleMeters: 10,
		Uploads: []reconstruct.UploadInfo{
			{FileName: "dashcam.mp4", DurationSeconds: 5, Width: 640, Height: 360},
			{FileName: "phone.mp4", DurationSeconds: 6, Width: 640, Height: 360},
		},
		Toolchain: []reconstruct.ToolStatus{
			{Name: "colmap", Status: "available"},
			{Name: "gdalinfo", Status: "available"},
		},
	})

	require.Equal(t, "case-1", artifact.CaseID)
	require.Len(t, artifact.VehicleTrack, 24)
	require.NotEmpty(t, artifact.Points)
	require.Greater(t, artifact.Speed.MeanKPH, 0.0)
	require.GreaterOrEqual(t, artifact.Speed.Confidence, 0.7)
	require.Contains(t, artifact.ReportMarkdown, "Accident Reconstruction Report")
}
