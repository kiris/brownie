package models

import (
	"github.com/pkg/errors"
	"io/ioutil"
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

func (w *Workspace) GetRepositories() ([]*Repository, error) {
	fileInfos, err := ioutil.ReadDir(w.RootDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get repositories: rootDir = %s", w.RootDir)
	}

	var repositories []*Repository
	for _, f := range fileInfos {
		repository := w.GetRepository(f.Name())
		if repository != nil {
			repositories = append(repositories, repository)
		}
	}
	return repositories, nil
}

func (w *Workspace) GetRepository(name string) *Repository {
	path := w.getRepositoryPath(name)
	return GetProject(path)
}


func (w *Workspace) getRepositoryPath(name string) string {
	return filepath.Join(w.RootDir, name)
}