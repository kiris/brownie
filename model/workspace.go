package model

import (
	"path/filepath"
)

type Workspace struct {
	RootDir string
}

func NewWorkspace(rootDir string) *Workspace {
	return &Workspace {
		RootDir: rootDir,
	}
}

func (w *Workspace) GetProject(name string) *Project {
	path := w.getProjectPath(name)
	return GetProject(path)
}

func (w *Workspace) getProjectPath(project string) string {
	return filepath.Join(w.RootDir, project)
}