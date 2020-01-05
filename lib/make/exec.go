package make

import (
	"os/exec"
)

type Make struct {
	Dir     string
	Targets []string
	Args    map[string]string
	DryRun  bool

	PrintDataBase bool
	NoBuiltinRules bool
	NoBuiltinVariables bool

}

func (m *Make) Exec() (string, string, error)  {
	cmd := exec.Command("make")
	cmd.Args = append(cmd.Args, m.options()...)
	cmd.Args = append(cmd.Args, m.args()...)
	cmd.Args = append(cmd.Args, m.Targets...)

	out, err := cmd.CombinedOutput()
	return cmd.String(), string(out), err
}

func (m *Make) options() []string {
	var options []string
	if m.Dir != "" {
		options = append(options, "-C", m.Dir)
	}
	if m.DryRun {
		options = append(options, "-n")
	}
	if m.PrintDataBase {
		options = append(options, "-p")
	}
	if m.NoBuiltinRules {
		options = append(options, "-r")
	}
	if m.NoBuiltinVariables {
		options = append(options, "-R")
	}
	return options
}

func (m *Make) args() []string {
	args := make([]string, len(m.Args))
	for key, value := range m.Args {
		args = append(args, key + "=" + value)
	}

	return args
}


