package registry_test

import (
	"fmt"
	"testing"

	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/registry"
)

func TestAuth(t *testing.T) {
	svc, err := registry.NewService(registry.ServiceOptions{V2Only: true})

	fmt.Println(err)
	authConfig := &types.AuthConfig{}
	authConfig.Username = "yuyang"
	authConfig.Password = "whosyourdaddy"
	// authConfig.ServerAddress = "harbor.docker.i.fbank.com"
	authConfig.ServerAddress = "registry.docker.i.fbank.com"
	status, token, err := svc.Auth(context.Background(), authConfig, "docker/18.06.1 go/go1.10.3 git-commit/29fccbcb7-unsupported kernel/3.10.0-514.el7.x86_64 os/linux arch/amd64 UpstreamClient(Docker-Client/18.06.1 \x5C(linux\x5C))")
	fmt.Println(status, token, err)
}
