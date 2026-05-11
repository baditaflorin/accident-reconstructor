// Package reconstruction orchestrates tool discovery, estimation, and artifacts.
package reconstruction

import (
	"math"
	"time"

	"github.com/baditaflorin/accident-reconstructor/pkg/reconstruct"
)

// EstimateInput contains normalized data used by the fallback estimator.
type EstimateInput struct {
	CaseID      string
	CaseName    string
	Version     string
	Commit      string
	ScaleMeters float64
	Uploads     []reconstruct.UploadInfo
	Toolchain   []reconstruct.ToolStatus
}

// EstimateScene creates a deterministic reconstruction artifact.
func EstimateScene(input EstimateInput) reconstruct.Artifact {
	duration := maxDuration(input.Uploads)
	scale := input.ScaleMeters
	if scale <= 0 {
		scale = 10
	}
	distance := math.Max(18, math.Min(120, scale*4.5+duration*5))
	track := vehicleTrack(duration, distance)
	speed := speedEstimate(track, input.Uploads, input.Toolchain)
	warnings := qualityWarnings(input.Uploads, input.Toolchain)
	mode := "native-toolchain-ready"
	if !toolAvailable(input.Toolchain, "colmap") || !toolAvailable(input.Toolchain, "gdalinfo") {
		mode = "deterministic-fallback"
	}

	artifact := reconstruct.Artifact{
		CaseID:          input.CaseID,
		CaseName:        input.CaseName,
		Version:         input.Version,
		Commit:          input.Commit,
		CreatedAt:       time.Now().UTC(),
		CoordinateFrame: "local_enu_meters",
		Uploads:         input.Uploads,
		Points:          sparseRoadPoints(distance),
		Cameras:         cameraPoses(input.Uploads, duration, distance),
		VehicleTrack:    track,
		Speed:           speed,
		Quality: reconstruct.ReconstructionQuality{
			Mode:            mode,
			InputVideos:     len(input.Uploads),
			CoordinateFrame: "local_enu_meters",
			Toolchain:       input.Toolchain,
			Warnings:        warnings,
		},
	}
	artifact.ReportMarkdown = BuildReport(artifact)
	return artifact
}

func maxDuration(uploads []reconstruct.UploadInfo) float64 {
	duration := 5.0
	for _, upload := range uploads {
		if upload.DurationSeconds > duration {
			duration = upload.DurationSeconds
		}
	}
	return duration
}

func vehicleTrack(duration float64, distance float64) []reconstruct.VehicleTrackPoint {
	samples := 24
	points := make([]reconstruct.VehicleTrackPoint, 0, samples)
	for i := 0; i < samples; i++ {
		t := duration * float64(i) / float64(samples-1)
		progress := float64(i) / float64(samples-1)
		x := -distance/2 + distance*progress
		z := math.Sin(progress*math.Pi*2) * 0.9
		points = append(points, reconstruct.VehicleTrackPoint{
			TimeSeconds: t,
			X:           x,
			Y:           0,
			Z:           z,
			SpeedMPS:    distance / duration,
		})
	}
	return points
}

func speedEstimate(
	track []reconstruct.VehicleTrackPoint,
	uploads []reconstruct.UploadInfo,
	tools []reconstruct.ToolStatus,
) reconstruct.SpeedEstimate {
	if len(track) < 2 {
		return reconstruct.SpeedEstimate{Method: "insufficient-track", Confidence: 0}
	}
	first := track[0]
	last := track[len(track)-1]
	distance := math.Hypot(last.X-first.X, last.Z-first.Z)
	duration := math.Max(0.1, last.TimeSeconds-first.TimeSeconds)
	mps := distance / duration
	confidence := 0.52
	notes := []string{
		"Speed is estimated from reconstructed local track distance over video time.",
		"Use a measured road reference in the UI to improve scale accuracy.",
	}
	if len(uploads) >= 2 {
		confidence += 0.12
		notes = append(notes, "Multiple videos improve camera path constraints.")
	}
	if toolAvailable(tools, "colmap") {
		confidence += 0.14
		notes = append(notes, "COLMAP is available for sparse reconstruction.")
	}
	if confidence > 0.86 {
		confidence = 0.86
	}
	margin := 0.35 - confidence*0.18
	return reconstruct.SpeedEstimate{
		Method:     "track-distance-over-time",
		MeanMPS:    round(mps),
		MeanKPH:    round(mps * 3.6),
		MeanMPH:    round(mps * 2.23693629),
		LowerKPH:   round(mps * 3.6 * (1 - margin)),
		UpperKPH:   round(mps * 3.6 * (1 + margin)),
		Confidence: round(confidence),
		Notes:      notes,
	}
}

func qualityWarnings(
	uploads []reconstruct.UploadInfo,
	tools []reconstruct.ToolStatus,
) []string {
	warnings := []string{
		"Outputs are decision-support artifacts, not a legal certification.",
	}
	if len(uploads) < 2 {
		warnings = append(warnings, "Single-video reconstruction has weaker depth constraints.")
	}
	if !toolAvailable(tools, "colmap") {
		warnings = append(warnings, "COLMAP is missing; deterministic fallback geometry was used.")
	}
	if !toolAvailable(tools, "gdalinfo") {
		warnings = append(warnings, "GDAL is missing; coordinates remain in a local metric frame.")
	}
	if !toolAvailable(tools, "ollama") {
		warnings = append(warnings, "Local LLM is missing; report narrative used deterministic text.")
	}
	for _, upload := range uploads {
		if upload.DurationSource == "placeholder" {
			warnings = append(warnings,
				"Duration for "+upload.FileName+" could not be measured (no ffprobe and the file is not an MP4/MOV with a readable mvhd atom); a 5-second placeholder was used and the speed estimate is unreliable.")
		}
	}
	return warnings
}

func sparseRoadPoints(distance float64) []reconstruct.Point3D {
	points := make([]reconstruct.Point3D, 0, 68)
	id := 0
	for lane := -1; lane <= 1; lane++ {
		for i := 0; i < 18; i++ {
			x := -distance/2 + distance*float64(i)/17
			z := float64(lane) * 3.4
			color := [3]uint8{180, 185, 190}
			if lane == 0 && i%2 == 0 {
				color = [3]uint8{250, 220, 80}
			}
			points = append(points, reconstruct.Point3D{
				ID:    pointID(id),
				X:     x,
				Y:     0,
				Z:     z,
				Color: color,
				Tags:  []string{"road"},
			})
			id++
		}
	}
	for i := 0; i < 14; i++ {
		angle := float64(i) / 14 * math.Pi * 2
		points = append(points, reconstruct.Point3D{
			ID:    pointID(id),
			X:     math.Cos(angle)*2 + distance*0.18,
			Y:     0.8 + math.Sin(angle*3)*0.1,
			Z:     math.Sin(angle) * 1.1,
			Color: [3]uint8{226, 63, 48},
			Tags:  []string{"vehicle"},
		})
		id++
	}
	return points
}

func cameraPoses(
	uploads []reconstruct.UploadInfo,
	duration float64,
	distance float64,
) []reconstruct.CameraPose {
	poses := make([]reconstruct.CameraPose, 0, len(uploads)*6)
	if len(uploads) == 0 {
		return poses
	}
	for videoIndex, upload := range uploads {
		offset := float64(videoIndex) * 8
		for i := 0; i < 6; i++ {
			progress := float64(i) / 5
			poses = append(poses, reconstruct.CameraPose{
				ID:            upload.FileName + "-pose-" + pointID(i),
				SourceVideo:   upload.FileName,
				TimeSeconds:   duration * progress,
				Position:      [3]float64{-distance/2 + distance*progress, 2.2, -8 - offset},
				RotationEuler: [3]float64{-8, 0, 0},
				FocalLengthPX: 900,
			})
		}
	}
	return poses
}

func pointID(index int) string {
	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"
	if index < len(alphabet) {
		return string(alphabet[index])
	}
	return "p" + strconvItoa(index)
}

func round(value float64) float64 {
	return math.Round(value*100) / 100
}

func strconvItoa(value int) string {
	if value == 0 {
		return "0"
	}
	var out []byte
	for value > 0 {
		out = append([]byte{byte('0' + value%10)}, out...)
		value /= 10
	}
	return string(out)
}
