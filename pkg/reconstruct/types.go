// Package reconstruct contains public API data shapes for reconstruction artifacts.
package reconstruct

import "time"

// CaseStatus describes the lifecycle state of a reconstruction case.
type CaseStatus string

// Case status values returned by the runtime API.
const (
	StatusQueued     CaseStatus = "queued"
	StatusProcessing CaseStatus = "processing"
	StatusComplete   CaseStatus = "complete"
	StatusFailed     CaseStatus = "failed"
)

// ErrorResponse is the structured JSON error returned by the API.
type ErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// UploadInfo records normalized metadata for an uploaded video.
//
// DurationSource records how the duration was obtained — "ffprobe" when the
// native tool was available, "mp4_atom" when the pure-Go MP4/MOV parser
// resolved it from the container header, or "placeholder" when neither
// worked and a 5-second default was applied. Downstream estimators inspect
// this to flag low-confidence cases.
type UploadInfo struct {
	FileName        string  `json:"fileName"`
	SizeBytes       int64   `json:"sizeBytes"`
	DurationSeconds float64 `json:"durationSeconds"`
	DurationSource  string  `json:"durationSource,omitempty"`
	Width           int     `json:"width"`
	Height          int     `json:"height"`
	FrameRate       float64 `json:"frameRate"`
	SHA256          string  `json:"sha256"`
}

// ToolStatus describes one external reconstruction tool on the backend.
type ToolStatus struct {
	Name    string `json:"name"`
	Path    string `json:"path,omitempty"`
	Version string `json:"version,omitempty"`
	Status  string `json:"status"`
}

// Point3D is a sparse reconstructed scene point in local metric space.
type Point3D struct {
	ID    string   `json:"id"`
	X     float64  `json:"x"`
	Y     float64  `json:"y"`
	Z     float64  `json:"z"`
	Color [3]uint8 `json:"color"`
	Tags  []string `json:"tags,omitempty"`
}

// CameraPose is an estimated camera position for a source video frame.
type CameraPose struct {
	ID            string     `json:"id"`
	SourceVideo   string     `json:"sourceVideo"`
	TimeSeconds   float64    `json:"timeSeconds"`
	Position      [3]float64 `json:"position"`
	RotationEuler [3]float64 `json:"rotationEuler"`
	FocalLengthPX float64    `json:"focalLengthPx"`
}

// VehicleTrackPoint is one time-synchronized vehicle position estimate.
type VehicleTrackPoint struct {
	TimeSeconds float64 `json:"timeSeconds"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Z           float64 `json:"z"`
	SpeedMPS    float64 `json:"speedMps"`
}

// SpeedEstimate summarizes speed and uncertainty for the reconstructed track.
type SpeedEstimate struct {
	Method     string   `json:"method"`
	MeanMPS    float64  `json:"meanMps"`
	MeanKPH    float64  `json:"meanKph"`
	MeanMPH    float64  `json:"meanMph"`
	LowerKPH   float64  `json:"lowerKph"`
	UpperKPH   float64  `json:"upperKph"`
	Confidence float64  `json:"confidence"`
	Notes      []string `json:"notes"`
}

// ReconstructionQuality explains the pipeline mode and known limitations.
type ReconstructionQuality struct {
	Mode             string       `json:"mode"`
	InputVideos      int          `json:"inputVideos"`
	CoordinateFrame  string       `json:"coordinateFrame"`
	Toolchain        []ToolStatus `json:"toolchain"`
	Warnings         []string     `json:"warnings"`
	ReprojectionRMSE float64      `json:"reprojectionRmse,omitempty"`
}

// Artifact is the full reconstruction result consumed by the frontend.
type Artifact struct {
	CaseID          string                `json:"caseId"`
	CaseName        string                `json:"caseName"`
	Version         string                `json:"version"`
	Commit          string                `json:"commit"`
	CreatedAt       time.Time             `json:"createdAt"`
	CoordinateFrame string                `json:"coordinateFrame"`
	Uploads         []UploadInfo          `json:"uploads"`
	Points          []Point3D             `json:"points"`
	Cameras         []CameraPose          `json:"cameras"`
	VehicleTrack    []VehicleTrackPoint   `json:"vehicleTrack"`
	Speed           SpeedEstimate         `json:"speed"`
	Quality         ReconstructionQuality `json:"quality"`
	ReportMarkdown  string                `json:"reportMarkdown"`
}

// CaseSummary is the polling-friendly state returned for a case.
type CaseSummary struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Status      CaseStatus     `json:"status"`
	Message     string         `json:"message"`
	Progress    int            `json:"progress"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Uploads     []UploadInfo   `json:"uploads"`
	ArtifactURL string         `json:"artifactUrl,omitempty"`
	ReportURL   string         `json:"reportUrl,omitempty"`
	Error       *ErrorResponse `json:"error,omitempty"`
}

// CreateCaseResponse wraps the accepted case returned after upload.
type CreateCaseResponse struct {
	Case CaseSummary `json:"case"`
}
