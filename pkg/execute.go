package pkg



import (
	"fmt"
	"log"
	"os/exec"
)

func ExecuteCmds(commands []string) {
	for _, cmd := range commands {
		out, err := exec.Command("bash", "-c", cmd).Output()
		output := string(out)
		fmt.Println(output)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ExecuteCmd(command string) (string, error) {
	out, err := exec.Command("bash", "-c", command).Output()
	output := string(out)
	return output, err
}

func ExecNoShell(command []string) (string, error) {
	out, err := exec.Command(command[0], command[1:]...).Output()
	output := string(out)
	return output, err
}
