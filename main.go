/*
Copyright 2021 Maximilian Geberl.

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
	"flag"
	"os"
	"time"

	_ "go.uber.org/automaxprocs"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	kubeteachv1alpha1 "github.com/dergeberl/kubeteach/api/v1alpha1"
	"github.com/dergeberl/kubeteach/controllers"
	kubeteachdashboard "github.com/dergeberl/kubeteach/pkg/dashboard"
	kubeteachmetrics "github.com/dergeberl/kubeteach/pkg/metrics"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(kubeteachv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var debugMode bool
	var requeueTimeTaskDefinition int
	var requeueTimeExerciseSet int
	var enableDashboard bool
	var dashboardListenAddr string
	var dashboardContent string
	var dashboardBasicAuthUser string
	var dashboardBasicAuthPassword string
	var dashboardWebterminalEnable bool
	var dashboardWebterminalHost string
	var dashboardWebterminalPort string
	var dashboardWebterminalCredentials string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&debugMode, "debug", false, "Enables debug logging mode")
	flag.IntVar(&requeueTimeTaskDefinition, "requeue-time-taskdefinition", 5, //nolint: gomnd
		"sets the requeue time in seconds for active and pending tasks")
	flag.IntVar(&requeueTimeExerciseSet, "requeue-time-exerciseset", 60, //nolint: gomnd
		"sets the requeue time in seconds for exercisesets")
	flag.BoolVar(&enableDashboard, "dashboard", false,
		"Enable dashboard for kubeteach.")
	flag.StringVar(&dashboardListenAddr, "dashboard-bind-address", ":8090",
		"Address that dashboard endpoint binds to.")
	flag.StringVar(&dashboardContent, "dashboard-content", "/dashboard",
		"The folder that contains the static files for the dashboard.")
	flag.StringVar(&dashboardBasicAuthUser, "dashboard-basic-auth-user", "",
		"Username for a basic auth. Can be also set via ENV: "+kubeteachdashboard.EnvDashboardBasicAuthUser)
	flag.StringVar(&dashboardBasicAuthPassword, "dashboard-basic-auth-password", "",
		"password for a basic auth. Can be also set via ENV: "+kubeteachdashboard.EnvDashboardBasicAuthPassword)
	flag.BoolVar(&dashboardWebterminalEnable, "dashboard-webterminal", false,
		"Enable webterminal forwarding in kubeteach dashboard.")
	flag.StringVar(&dashboardWebterminalHost, "dashboard-webterminal-host", "kubeteach-core-dashboard-webterminal",
		"Host for the webterminal container.")
	flag.StringVar(&dashboardWebterminalPort, "dashboard-webterminal-port", "8080",
		"Port for the webterminal Container.")
	flag.StringVar(&dashboardWebterminalCredentials, "dashboard-webterminal-credentials", "",
		"Basic auth for the connection to webterminal container (format user:password). "+
			"Can be also set via ENV: "+kubeteachdashboard.EnvWebterminalCredentials)

	opts := zap.Options{
		Development: debugMode,
	}
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443, // nolint: gomnd
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "06237eb5.geberl.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.TaskDefinitionReconciler{
		Client:      mgr.GetClient(),
		Log:         ctrl.Log.WithName("controllers").WithName("TaskDefinition"),
		Scheme:      mgr.GetScheme(),
		Recorder:    mgr.GetEventRecorderFor("Task"),
		RequeueTime: time.Duration(requeueTimeTaskDefinition) * time.Second,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "TaskDefinition")
		os.Exit(1)
	}
	if err = (&controllers.ExerciseSetReconciler{
		Client:      mgr.GetClient(),
		Log:         ctrl.Log.WithName("controllers").WithName("ExerciseSet"),
		Scheme:      mgr.GetScheme(),
		RequeueTime: time.Duration(requeueTimeExerciseSet) * time.Second,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "ExerciseSet")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	// register metrics endpoint
	metrics.Registry.MustRegister(
		kubeteachmetrics.New(
			mgr.GetClient(),
			ctrl.Log.WithName("metrics"),
		),
	)

	// start api if enabled
	if enableDashboard {
		setupLog.Info("starting api")
		dashboardConfig := kubeteachdashboard.New(mgr.GetClient(),
			dashboardListenAddr,
			dashboardContent,
			dashboardBasicAuthUser,
			dashboardBasicAuthPassword,
			dashboardWebterminalEnable,
			dashboardWebterminalHost,
			dashboardWebterminalPort,
			dashboardWebterminalCredentials)
		go func() {
			if err := dashboardConfig.Run(); err != nil {
				setupLog.Error(err, "problem running api")
				os.Exit(1)
			}
		}()
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
