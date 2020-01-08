package models

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/kiris/brownie/lib/file"
	"github.com/kiris/brownie/lib/make"
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

func GetProject(path string) *Repository {
	if !file.IsExistsDir(path) {
		return nil
	}

	return &Repository{
		Name: filepath.Base(path),
		Path: path,
	}
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