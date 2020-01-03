package brownie

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/kiris/brownie/pkg/file"
	"github.com/kiris/brownie/pkg/make"
)

type Project struct {
	name string
	path string
}


type ProjectConfig struct {
}

type Makefile struct {
	targets []string
}

type MakeResult struct {
	project *Project
	branch  string
	targets []string
	exec    string
	output  string
	success bool
	error   error
}

func GetProject(path string) *Project {
	if !file.IsExistsDir(path) {
		return nil
	}

	return &Project {
		name: filepath.Base(path),
		path: path,
	}
}

func (p *Project) ExecMake(targets []string) *MakeResult {
	cmd := make.Make {
		Dir    : p.path,
		Targets: targets,
		DryRun : false,
	}
	exec, output, err := cmd.Exec()

	// TODO
	log.Info(string(output))
	if err != nil {
		log.WithField("cause", err).Info("Failed to exec make.")
	}

	return &MakeResult{
		project: p,
		branch:  "master",
		targets: targets,
		exec:    exec,
		output:  output,
		success: err == nil,
		error:   err,
	}
}

//func getConfig(_ string) (*ProjectConfig, error) {
//	return nil, nil
//}
//
//func getMakefile