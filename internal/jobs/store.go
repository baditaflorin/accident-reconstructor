package jobs

import (
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/baditaflorin/accident-reconstructor/pkg/reconstruct"
)

type Case struct {
	Summary  reconstruct.CaseSummary
	WorkDir  string
	Artifact *reconstruct.Artifact
	Report   string
}

type Store struct {
	mu    sync.RWMutex
	cases map[string]Case
}

func NewStore() *Store {
	return &Store{cases: make(map[string]Case)}
}

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

func (s *Store) SetProcessing(id string, message string, progress int) bool {
	return s.update(id, func(item *Case) {
		item.Summary.Status = reconstruct.StatusProcessing
		item.Summary.Message = message
		item.Summary.Progress = progress
	})
}

func (s *Store) SetWorkDir(id string, workDir string) bool {
	return s.update(id, func(item *Case) {
		item.WorkDir = workDir
	})
}

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

func (s *Store) Fail(id string, err reconstruct.ErrorResponse) bool {
	return s.update(id, func(item *Case) {
		item.Summary.Status = reconstruct.StatusFailed
		item.Summary.Message = err.Message
		item.Summary.Progress = 100
		item.Summary.Error = &err
	})
}

func (s *Store) Get(id string) (Case, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.cases[id]
	return item, ok
}

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
