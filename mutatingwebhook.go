package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type networkHealthSidecarInjector struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (a *networkHealthSidecarInjector) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}

func (a *networkHealthSidecarInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
	// unmarshal
	pod := &v1.Pod{}
	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// process
	prepare(pod)

	container := container(pod)
	if container != nil {
		pod.Spec.Containers = append(pod.Spec.Containers, *container)
	}

	// marshal
	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

func prepare(pod *v1.Pod) {
	if pod.Labels == nil {
		pod.Labels = map[string]string{}
	}

	_, found := pod.Labels["network-health-sidecar/enabled"]
	if !found {
		pod.Labels["network-health-sidecar/enabled"] = "true"
	}
}

func container(pod *v1.Pod) *v1.Container {
	if strings.ToLower(pod.Labels["network-health-sidecar/enabled"]) != "true" {
		return nil
	}

	var err error
	var port int

	portStr, usePort := pod.Labels["network-health-sidecar/port"]
	if usePort {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			log.Printf(
				"failed to parse label network-health-sidecar/port to int for pod %s in namespace %s: value %s",
				pod.ObjectMeta.Name,
				pod.ObjectMeta.Namespace,
				pod.Labels["network-health-sidecar/port"],
			)
			usePort = false
			port = 0
		}
	}

	container := &v1.Container{
		Name:  "network-health-sidecar",
		Image: "ghcr.io/damoon/network-health-sidecar:latest",
	}

	if usePort {
		container.Args = []string{
			fmt.Sprintf("--addr=:%d", port),
		}
		container.Ports = []v1.ContainerPort{
			{
				Name:          "network-health",
				ContainerPort: int32(port),
				Protocol:      v1.ProtocolTCP,
			},
		}
		container.ReadinessProbe = &v1.Probe{
			Handler: v1.Handler{
				HTTPGet: &v1.HTTPGetAction{
					Path: "/healthz",
					Port: intstr.IntOrString{
						IntVal: int32(port),
					},
				},
			},
		}

		return container
	}

	container.Args = []string{
		"--protocol=unix",
		"--addr=/tmp/network-health.socket",
	}
	container.ReadinessProbe = &v1.Probe{
		Handler: v1.Handler{
			Exec: &v1.ExecAction{
				Command: []string{
					"network-health-client",
					"--protocol=unix",
					"--addr=/tmp/network-health.socket",
				},
			},
		},
	}

	return container
}
