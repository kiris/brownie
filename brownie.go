package brownie

import (
	log "github.com/sirupsen/logrus"

	"github.com/kiris/brownie/pkg/make"
)
type Brownie struct {
	WorkSpace string
}

func (b *Brownie) ExecMake(project string, targets []string) {
	cmd := make.Make {
		Dir: b.WorkSpace + "/" + project,
		Targets: targets,
		DryRun: false,
	}
	out, err := cmd.Exec()

	log.Info(string(out))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Failed to exec make.")
	}

}