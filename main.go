package main

import (
	"net/http"
	"os"

	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func init() {
	log.SetLogger(zap.New(zap.UseDevMode(true)))
}

func main() {
	entryLog := log.Log.WithName("entrypoint")

	options := manager.Options{}

	mgr, err := manager.New(config.GetConfigOrDie(), options)
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	entryLog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	endpoint := &webhook.Admission{
		Handler: &networkHealthSidecarInjector{
			Client: mgr.GetClient(),
		},
	}
	hookServer.Register("/webhook", endpoint)

	hookServer.WebhookMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte("ok"))
		if err != nil {
			entryLog.Error(err, "unable to write health check response")
		}
	})

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}
