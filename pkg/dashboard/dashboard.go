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

// Package dashboard for kubeteach
package dashboard

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sort"

	kubeteachv1alpha1 "github.com/dergeberl/kubeteach/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Environment variables
const (
	EnvWebterminalCredentials     = "WEBTERMINAL_CREDENTIALS"
	EnvDashboardBasicAuthUser     = "DASHBOARD_BASIC_AUTH_USER"
	EnvDashboardBasicAuthPassword = "DASHBOARD_BASIC_AUTH_PASSWORD"
)

// Config values for api
type Config struct {
	client                 client.Client
	listenAddr             string
	dashboardContent       string
	basicAuthUser          string
	basicAuthPassword      string
	webterminalEnable      bool
	webterminalHost        string
	webterminalPort        string
	webterminalCredentials string
}

type task struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UID         string `json:"uid"`
}

type tasks []task

func (a tasks) Len() int           { return len(a) }
func (a tasks) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a tasks) Less(i, j int) bool { return a[i].Name < a[j].Name }

type taskStatus struct {
	Status string `json:"status"`
}

// New creates a new config for the api
func New(
	client client.Client,
	listenAddr string,
	dashboardContent string,
	basicAuthUser string,
	basicAuthPassword string,
	webterminalEnable bool,
	webterminalHost string,
	webterminalPort string,
	webterminalCredentials string,
) Config {
	if os.Getenv(EnvWebterminalCredentials) != "" {
		webterminalCredentials = os.Getenv(EnvWebterminalCredentials)
	}
	if os.Getenv(EnvDashboardBasicAuthUser) != "" {
		basicAuthUser = os.Getenv(EnvDashboardBasicAuthUser)
	}
	if os.Getenv(EnvDashboardBasicAuthPassword) != "" {
		basicAuthPassword = os.Getenv(EnvDashboardBasicAuthPassword)
	}
	return Config{
		client:                 client,
		listenAddr:             listenAddr,
		dashboardContent:       dashboardContent,
		basicAuthUser:          basicAuthUser,
		basicAuthPassword:      basicAuthPassword,
		webterminalEnable:      webterminalEnable,
		webterminalHost:        webterminalHost,
		webterminalPort:        webterminalPort,
		webterminalCredentials: webterminalCredentials,
	}
}

// Run api webserver
func (c *Config) Run() error {
	return http.ListenAndServe(c.listenAddr, c.configureChi())
}

func (c *Config) configureChi() *chi.Mux {
	r := chi.NewRouter()
	if c.basicAuthUser != "" || c.basicAuthPassword != "" {
		r.Use(middleware.BasicAuth("", map[string]string{c.basicAuthUser: c.basicAuthPassword}))
	}
	r.Route("/", func(r chi.Router) {
		fs := http.FileServer(http.Dir(c.dashboardContent))
		r.Handle("/*", fs)

		r.Route("/api", func(r chi.Router) {
			r.Route("/tasks", func(r chi.Router) {
				r.Get("/", c.taskList)
			})
			r.Route("/taskstatus", func(r chi.Router) {
				r.Get("/{uid}", c.taskStatus)
			})
		})
		if c.webterminalEnable {
			r.Route("/shell", func(r chi.Router) {
				r.HandleFunc("/*", c.webterminalForward)
			})
		}
	})
	return r
}

func (c *Config) taskList(w http.ResponseWriter, r *http.Request) {
	if c.client == nil {
		http.Error(w, "Kubernetes client not functional", http.StatusInternalServerError)
		return
	}
	ctx := context.Background()
	taskList := &kubeteachv1alpha1.TaskDefinitionList{}
	err := c.client.List(ctx, taskList)
	if err != nil {
		http.Error(w, "Kubernetes client not functional", http.StatusInternalServerError)
		return
	}
	var tasksAPI tasks
	for _, t := range taskList.Items {
		tasksAPI = append(tasksAPI, task{
			Namespace:   t.Namespace,
			Name:        t.Name,
			UID:         string(t.UID),
			Title:       t.Spec.TaskSpec.Title,
			Description: t.Spec.TaskSpec.Description,
		})
	}
	sort.Sort(tasksAPI)
	output, err := json.Marshal(tasksAPI)
	if err != nil {
		http.Error(w, "JSON could not be generated", http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprint(w, string(output))
}

func (c *Config) taskStatus(w http.ResponseWriter, r *http.Request) {
	if c.client == nil {
		http.Error(w, "Kubernetes client not functional", http.StatusInternalServerError)
		return
	}
	ctx := context.Background()
	uid := chi.URLParam(r, "uid")
	taskList := &kubeteachv1alpha1.TaskDefinitionList{}
	err := c.client.List(ctx, taskList)
	if err != nil {
		http.Error(w, "Kubernetes client not functional", http.StatusInternalServerError)
		return
	}
	for _, t := range taskList.Items {
		if string(t.UID) == uid {
			output, err := json.Marshal(taskStatus{Status: *t.Status.State})
			if err != nil {
				http.Error(w, "JSON could not be generated", http.StatusInternalServerError)
				return
			}
			_, _ = fmt.Fprint(w, string(output))
			return
		}
	}
	http.Error(w, "No task with uid found", http.StatusNotFound)
}

func (c *Config) webterminalForward(writer http.ResponseWriter, request *http.Request) {
	rev := httputil.ReverseProxy{Director: func(request *http.Request) {
		request.Header.Del("Authorization")
		shellHost, _ := url.Parse("http://" + c.webterminalHost + ":" + c.webterminalPort + "/")
		request.URL.Scheme = shellHost.Scheme
		request.URL.Host = shellHost.Host
		request.URL.RawQuery = shellHost.RawQuery + request.URL.RawQuery
		if c.webterminalCredentials != "" {
			request.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.webterminalCredentials)))
		}
		request.Host = shellHost.Host
	}}
	rev.ServeHTTP(writer, request)
}
