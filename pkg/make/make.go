package make

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

type Make struct {
	Dir     string
	Branch  string
	Targets []string
	Args    map[string]string
	DryRun  bool
}

func (this *Make) Exec() (string, error)  {
	cmd := exec.Command(
		"make",
		"-C",
		this.Dir,
	)
	cmd.Args = append(cmd.Args, this.options()...)
	cmd.Args = append(cmd.Args, this.args()...)
	cmd.Args = append(cmd.Args, this.Targets...)
	log.Info("exec: " + cmd.String())

	out, err := cmd.CombinedOutput()
	return string(out), err
}

func (this *Make) options() []string {
	var options []string
	if this.DryRun {
		options = append(options, "-n")
	}
	return options
}

func (this *Make) args() []string {
	args := make([]string, len(this.Args))
	for key, value := range this.Args {
		args = append(args, key + "=" + value)
	}

	return args
}


