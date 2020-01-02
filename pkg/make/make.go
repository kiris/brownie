package make

import (
	"os/exec"
    "strings"
)


func ExecMake(dir string, _ string, targets ... string) (string, error)  {
	cmd := exec.Command("make", "-C", dir, strings.Join(targets, " "))
	out, err := cmd.CombinedOutput()

	return string(out), err
}
