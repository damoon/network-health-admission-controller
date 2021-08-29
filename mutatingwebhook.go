/*
Copyright 2018 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"encoding/json"
	"net/http"

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
	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}

	pod.Annotations["network-health-sidecar/enabled"] = "true"
}

func container(pod *v1.Pod) *v1.Container {
	return &v1.Container{
		Name:  "network-health-sidecar",
		Image: "ghcr.io/damoon/network-health-sidecar:latest",
		Ports: []v1.ContainerPort{
			{
				Name:          "network-health",
				ContainerPort: 8080,
				Protocol:      v1.ProtocolTCP,
			},
		},
		ReadinessProbe: &v1.Probe{
			Handler: v1.Handler{
				HTTPGet: &v1.HTTPGetAction{
					Path: "/healthz",
					Port: intstr.IntOrString{
						IntVal: 8080,
					},
				},
			},
		},
	}
}
