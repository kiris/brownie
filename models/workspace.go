package models

import (
	"github.com/kiris/brownie/lib/file"
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

func (w *Workspace) Repositories() ([]*Repository, error) {
	fileInfos, err := ioutil.ReadDir(w.RootDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get repositories: rootDir = %s", w.RootDir)
	}

	var repositories []*Repository
	for _, f := range fileInfos {
		repository := w.Repository(f.Name())
		if repository != nil {
			repositories = append(repositories, repository)
		}
	}
	return repositories, nil
}

func (w *Workspace) Repository(name string) *Repository {
	path := w.RepositoryPath(name)

	if !file.IsExistsDir(path) {
		return nil
	}

	return NewRepository(path)
}

func (w *Workspace) CreateRepository(url string) (*Repository, error) {
	return NewRepositoryByGitClone(w.RootDir, url)
}

func (w *Workspace) RepositoryPath(name string) string {
	return filepath.Join(w.RootDir, name)
}