package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
	"k8s.io/kubernetes/pkg/kubelet/util"
)

const (
	dockerep       = "unix:///var/run/dockershim.sock"
	containerdep   = "unix:///run/containerd/containerd.sock"
	crioep         = "unix:///run/crio/crio.sock"
	defaultTimeout = 2 * time.Second
)

func main() {
	var inputEp string
	var ep string
	fmt.Printf("Please input the runtime 'docker'/'containerd'/cri': ")
	fmt.Scanln(&inputEp)
	switch inputEp {
	case "docker":
		ep = dockerep
	case "containerd":
		ep = containerdep
	case "crioep":
		ep = crioep
	}
	addr, dialer, err := util.GetAddressAndDialer(ep)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to get GRPC addr/dialer due to %s", err))
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(defaultTimeout), grpc.WithContextDialer(dialer))
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to the GRPC endpoint %s", err))
	}
	runtimeClient := pb.NewRuntimeServiceClient(conn)
	response, err := runtimeClient.ListPodSandbox(context.Background(), &pb.ListPodSandboxRequest{})
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to get the pod list"))
	}
	for _, res := range response.GetItems() {
		log.Info(fmt.Sprintf("%s", res))
	}
}
