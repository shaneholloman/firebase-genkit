// Copyright 2024 Google LLC
// SPDX-License-Identifier: Apache-2.0


package core

import (
	"context"
	"os"
	"path/filepath"

	"github.com/firebase/genkit/go/internal/base"
)

// A FileFlowStateStore is a FlowStateStore that writes flowStates to files.
type FileFlowStateStore struct {
	dir string
}

// NewFileFlowStateStore creates a FileFlowStateStore that writes traces to the given
// directory. The directory is created if it does not exist.
func NewFileFlowStateStore(dir string) (*FileFlowStateStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &FileFlowStateStore{dir: dir}, nil
}

func (s *FileFlowStateStore) Save(ctx context.Context, id string, fs base.FlowStater) error {
	data, err := fs.ToJSON()
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(s.dir, base.Clean(id)), data, 0666)
}

func (s *FileFlowStateStore) Load(ctx context.Context, id string, pfs any) error {
	return base.ReadJSONFile(filepath.Join(s.dir, base.Clean(id)), pfs)
}
