package models

import (
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
}


type RepositoryConfig struct {
}

type Makefile struct {
	targets []string
}


type ExecMakeSetting struct {
}

type ExecMakeResult struct {
	Repository *Repository
	Branch     string
	Targets    []string
	Exec       string
	Output     string
	Success    bool
	Error      error
}

func NewRepository(path string) *Repository {
	return &Repository{
		Name: filepath.Base(path),
		Path: path,
	}
}

func NewRepositoryByGitClone(rootPath string, url string) (*Repository, error) {
	name := extractRepositoryName(url)
	_, err := git.PlainClone(rootPath + "/" + name, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})

	return nil, err
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

func (p *Repository) ExecMake(targets []string) *ExecMakeResult {
	cmd := make.Make {
		Dir    : p.Path,
		Targets: targets,
		DryRun : false,
	}
	exec, output, err := cmd.Exec()
	if err != nil {
		log.WithField("cause", err).Info("Failed to ExecCommand handleMakeCommand.")
	}

	return &ExecMakeResult{
		Repository: p,
		Branch :    "master",
		Targets:    targets,
		Exec   :    exec,
		Output :    output,
		Success:    err == nil,
		Error  :    err,
	}
}

func (p *Repository) CollectMakeTargets() ([]string, error) {
	cmd := make.Make {
		Dir               : p.Path,
		PrintDataBase     : true,
		NoBuiltinRules    : true,
		NoBuiltinVariables: true,
	}

	_, output, err := cmd.Exec()
	if err != nil {
		log.WithField("cause", err).Info("Failed to ExecCommand handleMakeCommand.")
		return nil, err
	}

	return make.ParseDataBase(output), nil
}

//func getConfig(_ string) (*RepositoryConfig, Error) {
//	return nil, nil
//}
//
//func getMakefile