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

// Package api to get metadata for dashboard
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	kubeteachv1alpha1 "github.com/dergeberl/kubeteach/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-chi/chi/v5"
)

// Config values for api
type Config struct {
	client     client.Client
	listenAddr string
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
func New(client client.Client, listenAddr string) Config {
	return Config{
		client:     client,
		listenAddr: listenAddr,
	}
}

// Run api webserver
func (c *Config) Run() error {
	return http.ListenAndServe(c.listenAddr, c.configureChi())
}

func (c *Config) configureChi() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api", func(r chi.Router) {
		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", c.taskList)
		})
		r.Route("/taskstatus", func(r chi.Router) {
			r.Get("/{uid}", c.taskStatus)
		})
	})
	return r
}

func (c *Config) taskList(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	tasklist := &kubeteachv1alpha1.TaskDefinitionList{}
	err := c.client.List(ctx, tasklist)
	if err != nil {
		// todo handle error
		return
	}
	var tasksAPI tasks
	for _, t := range tasklist.Items {
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
		// todo handle error
		return
	}
	_, _ = fmt.Fprint(w, string(output))
}

func (c *Config) taskStatus(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	uid := chi.URLParam(r, "uid")
	taskList := &kubeteachv1alpha1.TaskDefinitionList{}
	err := c.client.List(ctx, taskList)
	if err != nil {
		// todo handle error
		return
	}
	for _, t := range taskList.Items {
		if string(t.UID) == uid {
			output, err := json.Marshal(taskStatus{Status: *t.Status.State})
			if err != nil {
				// todo handle error
				return
			}
			_, _ = fmt.Fprint(w, string(output))
			return
		}
	}
	// todo handle error
}
