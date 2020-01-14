package models

import (
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/kiris/brownie/lib/make"
	"gopkg.in/src-d/go-git.v4"
)

type Repository struct {
	Name string
	Path string
	gitRepo *git.Repository
}


type RepositoryConfig struct {
}

type Makefile struct {
	targets []string
}


type RunMakeResult struct {
	Repository *Repository
	Branch     string
	Targets    []string
	Exec       string
	Output     string
	Success    bool
	Error      error
}

func NewRepository(path string) (*Repository, error) {
	gitRepo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Name: filepath.Base(path),
		Path: path,
		gitRepo: gitRepo,
	}, nil
}

func NewRepositoryByGitClone(rootPath string, url string) (*Repository, error) {
	name := extractRepositoryName(url)
	path := rootPath + "/" + name
	gitRepo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return nil, err
	}

	return &Repository{
		Name: name,
		Path: path,
		gitRepo: gitRepo,
	}, nil
}


func extractRepositoryName(name string) string {
	// Strip trailing slashes.
	for len(name) > 0 && name[len(name)-1] == '/' {
		name = name[0 : len(name)-1]
	}

	// Find the last element
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}

	// Find
	if i := strings.LastIndex(name, "."); i >= 0 {
		name = name[:i]
	}

	return name

}

func (p *Repository) RunMake(targets []string) *RunMakeResult {
	cmd := make.Make {
		Dir    : p.Path,
		Targets: targets,
		DryRun : false,
	}
	exec, output, err := cmd.Exec()
	if err != nil {
		log.WithField("cause", err).Warn("Failed to ExecCommand handleMakeCommand.")
	}

	return &RunMakeResult{
		Repository: p,
		Branch :    "master",
		Targets:    targets,
		Exec   :    exec,
		Output :    output,
		Success:    err == nil,
		Error  :    err,
	}
}

func (p *Repository) Branches() ([]string, error) {
	r, err := git.PlainOpen(p.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to branches: path %s", p.Path)
	}

	branches, err := r.Branches()

	var results []string
	_ := branches.ForEach(func(reference *plumbing.Reference) error {

		reference.Name().Short()
		return nil
	})
	for b :=  branches.ForEach() {

	}
}

func (p *Repository) Targets() ([]string, error) {
	cmd := make.Make {
		Dir               : p.Path,
		PrintDataBase     : true,
		NoBuiltinRules    : true,
		NoBuiltinVariables: true,
	}

	_, output, err := cmd.Exec()
	if err != nil {
		return nil, errors.Wrap(err, "failed to Targets")
	}

	return make.ParseDataBase(output), nil
}

func (p *Repository) fetch() error {

}

func (p *Repository) fetch() error {
	return nil
}