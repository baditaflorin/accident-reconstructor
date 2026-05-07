// Package jobs tracks reconstruction case state for the API process.
package jobs

import (
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/baditaflorin/accident-reconstructor/pkg/reconstruct"
)

// Case stores runtime case state and generated artifacts.
type Case struct {
	Summary  reconstruct.CaseSummary
	WorkDir  string
	Artifact *reconstruct.Artifact
	Report   string
}

// Store is an in-memory case registry for the current API process.
type Store struct {
	mu    sync.RWMutex
	cases map[string]Case
}

// NewStore creates an empty case store.
func NewStore() *Store {
	return &Store{cases: make(map[string]Case)}
}

// Create inserts a new queued case.
func (s *Store) Create(name string, workDir string, uploads []reconstruct.UploadInfo) Case {
	now := time.Now().UTC()
	id := uuid.NewString()
	item := Case{
		WorkDir: workDir,
		Summary: reconstruct.CaseSummary{
			ID:        id,
			Name:      name,
			Status:    reconstruct.StatusQueued,
			Message:   "Queued for reconstruction",
			Progress:  5,
			CreatedAt: now,
			UpdatedAt: now,
			Uploads:   uploads,
		},
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cases[id] = item
	return item
}

// SetProcessing marks a case as processing with progress.
func (s *Store) SetProcessing(id string, message string, progress int) bool {
	return s.update(id, func(item *Case) {
		item.Summary.Status = reconstruct.StatusProcessing
		item.Summary.Message = message
		item.Summary.Progress = progress
	})
}

// SetWorkDir records the artifact directory for a case.
func (s *Store) SetWorkDir(id string, workDir string) bool {
	return s.update(id, func(item *Case) {
		item.WorkDir = workDir
	})
}

// Complete stores the final artifact and report for a case.
func (s *Store) Complete(id string, artifact *reconstruct.Artifact, report string) bool {
	return s.update(id, func(item *Case) {
		item.Artifact = artifact
		item.Report = report
		item.Summary.Status = reconstruct.StatusComplete
		item.Summary.Message = "Reconstruction complete"
		item.Summary.Progress = 100
		item.Summary.ArtifactURL = "/api/v1/cases/" + id + "/artifact"
		item.Summary.ReportURL = "/api/v1/cases/" + id + "/report"
		item.Summary.Uploads = artifact.Uploads
	})
}

// Fail marks a case as failed with a structured API error.
func (s *Store) Fail(id string, err reconstruct.ErrorResponse) bool {
	return s.update(id, func(item *Case) {
		item.Summary.Status = reconstruct.StatusFailed
		item.Summary.Message = err.Message
		item.Summary.Progress = 100
		item.Summary.Error = &err
	})
}

// Get returns a case by ID.
func (s *Store) Get(id string) (Case, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.cases[id]
	return item, ok
}

// List returns all case summaries known by this process.
func (s *Store) List() []reconstruct.CaseSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]reconstruct.CaseSummary, 0, len(s.cases))
	for _, item := range s.cases {
		items = append(items, item.Summary)
	}
	return items
}

func (s *Store) update(id string, fn func(item *Case)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.cases[id]
	if !ok {
		return false
	}
	fn(&item)
	item.Summary.UpdatedAt = time.Now().UTC()
	s.cases[id] = item
	return true
}
