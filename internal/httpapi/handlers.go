package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/baditaflorin/accident-reconstructor/internal/reconstruction"
	"github.com/baditaflorin/accident-reconstructor/pkg/reconstruct"
)

func (a App) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a App) ready(w http.ResponseWriter, _ *http.Request) {
	if err := os.MkdirAll(a.Config.StorageDir, 0o755); err != nil {
		writeError(w, http.StatusServiceUnavailable, "storage_unavailable", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (a App) tools(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, reconstruction.DiscoverTools(r.Context()))
}

func (a App) listCases(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"cases": a.Store.List()})
}

func (a App) createCase(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, a.Config.MaxUploadBytes)
	if err := r.ParseMultipartForm(a.Config.MaxUploadBytes); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_multipart", "Upload must be multipart/form-data.")
		return
	}

	name := r.FormValue("case_name")
	if name == "" {
		name = "Untitled accident reconstruction"
	}
	scaleMeters, _ := strconv.ParseFloat(r.FormValue("scale_meters"), 64)
	if scaleMeters <= 0 {
		scaleMeters = 10
	}

	files := r.MultipartForm.File["videos"]
	if len(files) == 0 {
		writeError(w, http.StatusBadRequest, "videos_required", "Upload at least one video.")
		return
	}

	workingRoot := filepath.Join(a.Config.StorageDir, time.Now().UTC().Format("20060102"))
	if err := os.MkdirAll(workingRoot, 0o755); err != nil {
		writeError(w, http.StatusInternalServerError, "storage_error", err.Error())
		return
	}
	caseItem := a.Store.Create(name, "", nil)
	workDir := filepath.Join(workingRoot, caseItem.Summary.ID)
	a.Store.SetWorkDir(caseItem.Summary.ID, workDir)
	uploadsDir := filepath.Join(workDir, "uploads")
	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		writeError(w, http.StatusInternalServerError, "storage_error", err.Error())
		return
	}

	var paths []string
	var uploads []reconstruct.UploadInfo
	for _, header := range files {
		path, info, err := a.saveUpload(r.Context(), uploadsDir, header.Filename, header.Size, header.Open)
		if err != nil {
			writeError(w, http.StatusBadRequest, "upload_failed", err.Error())
			return
		}
		paths = append(paths, path)
		uploads = append(uploads, info)
		a.Metrics.UploadedBytes.Add(float64(info.SizeBytes))
	}

	a.Store.SetProcessing(caseItem.Summary.ID, "Video metadata extracted", 25)
	go a.runCase(r.Context(), caseItem.Summary.ID, name, workDir, paths, uploads, scaleMeters)

	updated, _ := a.Store.Get(caseItem.Summary.ID)
	writeJSON(w, http.StatusAccepted, reconstruct.CreateCaseResponse{Case: updated.Summary})
}

func (a App) saveUpload(
	ctx context.Context,
	dir string,
	name string,
	size int64,
	open func() (multipart.File, error),
) (string, reconstruct.UploadInfo, error) {
	if size == 0 {
		return "", reconstruct.UploadInfo{}, errors.New("empty uploads are not supported")
	}
	src, err := open()
	if err != nil {
		return "", reconstruct.UploadInfo{}, fmt.Errorf("open upload: %w", err)
	}
	defer src.Close()

	cleanName := filepath.Base(name)
	path := filepath.Join(dir, cleanName)
	dst, err := os.Create(path)
	if err != nil {
		return "", reconstruct.UploadInfo{}, fmt.Errorf("create upload: %w", err)
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		return "", reconstruct.UploadInfo{}, fmt.Errorf("save upload: %w", err)
	}
	if err := dst.Close(); err != nil {
		return "", reconstruct.UploadInfo{}, fmt.Errorf("close upload: %w", err)
	}
	info, err := reconstruction.ProbeVideo(ctx, path, cleanName)
	if err != nil {
		return "", reconstruct.UploadInfo{}, err
	}
	return path, info, nil
}

func (a App) runCase(
	parent context.Context,
	id string,
	name string,
	workDir string,
	paths []string,
	uploads []reconstruct.UploadInfo,
	scaleMeters float64,
) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.WithoutCancel(parent), 30*time.Minute)
	defer cancel()

	a.Store.SetProcessing(id, "Running reconstruction pipeline", 60)
	artifact, err := a.Processor.Process(ctx, reconstruction.ProcessInput{
		CaseID:      id,
		CaseName:    name,
		WorkDir:     workDir,
		UploadPaths: paths,
		Uploads:     uploads,
		ScaleMeters: scaleMeters,
	})
	a.Metrics.PipelineTime.Observe(time.Since(start).Seconds())
	if err != nil {
		a.Metrics.Jobs.WithLabelValues("failed").Inc()
		a.Store.Fail(id, reconstruct.ErrorResponse{
			Code:    "pipeline_failed",
			Message: err.Error(),
		})
		return
	}
	a.Metrics.Jobs.WithLabelValues("complete").Inc()
	a.Store.Complete(id, artifact, artifact.ReportMarkdown)
}

func (a App) getCase(w http.ResponseWriter, r *http.Request) {
	item, ok := a.Store.Get(chi.URLParam(r, "caseID"))
	if !ok {
		writeError(w, http.StatusNotFound, "case_not_found", "No case exists with that ID.")
		return
	}
	writeJSON(w, http.StatusOK, item.Summary)
}

func (a App) getArtifact(w http.ResponseWriter, r *http.Request) {
	item, ok := a.Store.Get(chi.URLParam(r, "caseID"))
	if !ok || item.Artifact == nil {
		writeError(w, http.StatusNotFound, "artifact_not_found", "The reconstruction artifact is not ready.")
		return
	}
	writeJSON(w, http.StatusOK, item.Artifact)
}

func (a App) getReport(w http.ResponseWriter, r *http.Request) {
	item, ok := a.Store.Get(chi.URLParam(r, "caseID"))
	if !ok || item.Report == "" {
		writeError(w, http.StatusNotFound, "report_not_found", "The report is not ready.")
		return
	}
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(item.Report))
}

func (a App) getBundle(w http.ResponseWriter, r *http.Request) {
	item, ok := a.Store.Get(chi.URLParam(r, "caseID"))
	if !ok {
		writeError(w, http.StatusNotFound, "case_not_found", "No case exists with that ID.")
		return
	}
	path := filepath.Join(item.WorkDir, "case-bundle.zip")
	http.ServeFile(w, r, path)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	writeJSON(w, status, reconstruct.ErrorResponse{Code: code, Message: message})
}
