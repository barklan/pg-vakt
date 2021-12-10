package pkg

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/foomo/htpasswd"
)

// TODO project should be replaced with something identifiable to be able to start multiple nginx containers
// and so they can reuse docker volume

func buildNginx(
	data *Data,
	basicAuthUsername,
	basicAuthPassword string,
) error {
	builderPath := data.CreateMediaDirIfNotExists("tempnginx")

	dockerFilePath := builderPath + "/Dockerfile"
	dockerFileStr := `FROM nginx:1.19.0-alpine
ARG BUILDKIT_INLINE_CACHE=1
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d
COPY htpasswd /etc/nginx/htpasswd
`

	if err := os.WriteFile(dockerFilePath, []byte(dockerFileStr), 0777); err != nil {
		return err
	}

	nginxFilePath := builderPath + "/nginx.conf"
	nginxFileStr := fmt.Sprintf(`server {

	auth_basic           "Restricted Area";
	auth_basic_user_file /etc/nginx/htpasswd;

    listen 9090;
    client_max_body_size 5M;

    location / {
        alias /home/app/media/%s/;
        autoindex on;
    }
}
`, "project")

	if err := os.WriteFile(nginxFilePath, []byte(nginxFileStr), 0777); err != nil {
		return err
	}

	htpasswdPath := builderPath + "/htpasswd"
	file, err := os.Create(htpasswdPath)
	if err != nil {
		log.Panic(err)
	}
	file.Close()

	err = htpasswd.SetPassword(htpasswdPath, basicAuthUsername, basicAuthPassword, htpasswd.HashBCrypt)
	if err != nil {
		log.Panic(err)
	}

	cmd := []string{"docker", "build", "-t", fmt.Sprintf("nginx:%s", "project"), builderPath}
	_, err = ExecNoShell(cmd)
	return err
}

func TemporaryNginx(
	data *Data,
	addressChan chan string,
	minutes int,
	basicAuthUsername,
	basicAuthPassword string,
) {
	data.CreateMediaDirIfNotExists("project")

	containerName := "tempnginx"
	alreadyRunning, err := CheckIfNamedContainerIsRunning(containerName)
	if err != nil {
		data.Send("Failed to check if container is running. Will not proceed.")
		return
	}
	if alreadyRunning {
		data.Send("Only one temporary nginx container can be running (project agnostic).")
		return
	}

	if err := buildNginx(data, basicAuthUsername, basicAuthPassword); err != nil {
		log.Println(err)
		data.Send("Failed to build temporary nginx server.")
		return
	}

	// TODO randomize port and name to allow for multiple nginx servers running simultaneously
	port := "9090"

	cmd := []string{
		"docker",
		"run",
		"--name",
		containerName,
		"--rm",
		"-d",
		"-v", "pg-vakt-media:/home/app/media",
		"-p", fmt.Sprintf("%s:%s", port, port),
		fmt.Sprintf("nginx:%s", "project"),
	}

	_, err = ExecNoShell(cmd)
	if err != nil {
		log.Println(err)
		data.Send("Failed to start temporary nginx server.")
		return
	}

	// TODO this should not be hardcoded.
	hostname := os.Getenv("HOST_HOSTNAME")
	addressChan <- fmt.Sprintf("%s:%s", hostname, port)

	time.Sleep(time.Duration(minutes) * time.Minute)

	_, err = ExecNoShell(
		[]string{"docker", "stop", containerName},
	)
	if err != nil {
		data.Send("Failed to stop temporary nginx server. Please intervene.")
		return
	}

	data.Send("Temporary nginx server stopped as planned.")
}
