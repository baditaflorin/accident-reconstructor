package reconstruction

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/baditaflorin/accident-reconstructor/internal/config"
	"github.com/baditaflorin/accident-reconstructor/pkg/reconstruct"
)

// Processor orchestrates reconstruction steps and artifact writing.
type Processor struct {
	Config config.Config
	LLM    LLMClient
}

// ProcessInput describes one backend reconstruction job.
type ProcessInput struct {
	CaseID      string
	CaseName    string
	WorkDir     string
	UploadPaths []string
	Uploads     []reconstruct.UploadInfo
	ScaleMeters float64
}

// NewProcessor creates a reconstruction processor from runtime config.
func NewProcessor(cfg config.Config) Processor {
	return Processor{
		Config: cfg,
		LLM: LLMClient{
			BaseURL: cfg.OllamaBaseURL,
			Model:   cfg.OllamaModel,
		},
	}
}

// Process runs reconstruction and writes artifacts into the case directory.
func (p Processor) Process(ctx context.Context, input ProcessInput) (*reconstruct.Artifact, error) {
	tools := DiscoverTools(ctx)
	uploads := input.Uploads
	if len(uploads) == 0 {
		for _, path := range input.UploadPaths {
			info, err := ProbeVideo(ctx, path, filepath.Base(path))
			if err != nil {
				return nil, err
			}
			uploads = append(uploads, info)
		}
	}

	artifact := EstimateScene(EstimateInput{
		CaseID:      input.CaseID,
		CaseName:    input.CaseName,
		Version:     p.Config.Version,
		Commit:      p.Config.Commit,
		ScaleMeters: input.ScaleMeters,
		Uploads:     uploads,
		Toolchain:   tools,
	})

	if note, err := p.LLM.Summary(ctx, artifact); err == nil && note != "" {
		artifact.ReportMarkdown += "\n\n## Local LLM Summary\n\n" + note + "\n"
	} else if err != nil {
		artifact.Quality.Warnings = append(artifact.Quality.Warnings, "Local LLM summary failed: "+err.Error())
		artifact.ReportMarkdown = BuildReport(artifact)
	}

	if err := writeArtifactFiles(input.WorkDir, artifact); err != nil {
		return nil, err
	}
	return &artifact, nil
}

func writeArtifactFiles(workDir string, artifact reconstruct.Artifact) error {
	if err := os.MkdirAll(workDir, 0o700); err != nil {
		return fmt.Errorf("create case directory: %w", err)
	}
	jsonPath := filepath.Join(workDir, "reconstruction.json")
	reportPath := filepath.Join(workDir, "report.md")
	bundlePath := filepath.Join(workDir, "case-bundle.zip")

	body, err := json.MarshalIndent(artifact, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal artifact: %w", err)
	}
	if err := os.WriteFile(jsonPath, body, 0o600); err != nil {
		return fmt.Errorf("write reconstruction artifact: %w", err)
	}
	if err := os.WriteFile(reportPath, []byte(artifact.ReportMarkdown), 0o600); err != nil {
		return fmt.Errorf("write report artifact: %w", err)
	}
	if err := writeZip(bundlePath, map[string][]byte{
		"reconstruction.json": body,
		"report.md":           []byte(artifact.ReportMarkdown),
	}); err != nil {
		return err
	}
	return nil
}

func writeZip(path string, files map[string][]byte) error {
	// #nosec G304 -- path is generated inside a case artifact directory.
	out, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("create bundle: %w", err)
	}
	zipWriter := zip.NewWriter(out)

	names := []string{"reconstruction.json", "report.md"}
	for _, name := range names {
		body, ok := files[name]
		if !ok {
			continue
		}
		writer, err := zipWriter.Create(strings.TrimPrefix(name, "/"))
		if err != nil {
			return fmt.Errorf("add bundle file: %w", err)
		}
		if _, err := writer.Write(body); err != nil {
			return fmt.Errorf("write bundle file: %w", err)
		}
	}
	if err := zipWriter.Close(); err != nil {
		_ = out.Close()
		return fmt.Errorf("close bundle writer: %w", err)
	}
	if err := out.Close(); err != nil {
		return fmt.Errorf("close bundle: %w", err)
	}
	return nil
}
