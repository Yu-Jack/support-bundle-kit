package utils

import (
	"bytes"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/cmd/exec"
)

type PodExec struct {
	RestConfig *rest.Config
	*kubernetes.Clientset
	Namespace     string
	PodName       string
	ContainerName string
}

func NewPodExec(config rest.Config, clientset *kubernetes.Clientset, namespace, podname, containername string) *PodExec {
	config.APIPath = "/api"
	config.GroupVersion = &schema.GroupVersion{Version: "v1"}
	config.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}
	return &PodExec{
		RestConfig:    &config,
		Clientset:     clientset,
		Namespace:     namespace,
		PodName:       podname,
		ContainerName: containername,
	}
}

// ExecCmd
// EXAMPLES
// // Execute ls -l /tmp on container test-container on pod test in namespace default
// // and print the resulting output
// in, out, errOut, err := podExec.ExecCmd([]string{"ls", "-l", "/tmp"}, "default", "test", "test-container")
//
//	    if err != nil {
//			fmt.Printf("%v\n", err)
//		}
//		fmt.Println("out:")
//		fmt.Printf("%s", out.String()) // will execute ls -l /tmp in the pod and output the result
func (p *PodExec) ExecCmd(command []string) (*bytes.Buffer, *bytes.Buffer, *bytes.Buffer, error) {
	ioStreams, in, out, errOut := genericclioptions.NewTestIOStreams()
	options := &exec.ExecOptions{
		StreamOptions: exec.StreamOptions{
			Namespace:       p.Namespace,
			PodName:         p.PodName,
			ContainerName:   p.ContainerName,
			Stdin:           true,
			TTY:             false,
			Quiet:           false,
			InterruptParent: nil,
			IOStreams:       ioStreams,
		},
		Command:       command,
		Executor:      &exec.DefaultRemoteExecutor{},
		PodClient:     p.Clientset.CoreV1(),
		GetPodTimeout: 0,
		Config:        p.RestConfig,
	}

	err := options.Run()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Could not run exec operation: %v", err)
	}

	return in, out, errOut, nil
}
