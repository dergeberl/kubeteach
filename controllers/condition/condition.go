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

// Package condition is used in kubeteach to run condition checks against kubernetes api
package condition

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
	"github.com/tidwall/gjson"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	teachv1alpha1 "github.com/dergeberl/kubeteach/api/v1alpha1"
)

// Checks is used for configuration of the condition checks
type Checks struct {
	Client client.Client
	Log    logr.Logger
}

// ApplyChecks apply all TaskConditions and returns true if all conditions are successful
func (c *Checks) ApplyChecks(
	ctx context.Context,
	taskConditions []teachv1alpha1.TaskCondition,
) (bool, error) {
	if len(taskConditions) < 1 {
		return false, errors.New("no checks to apply")
	}
	for _, taskCondition := range taskConditions {
		success, err := c.runTaskCondition(ctx, taskCondition)
		if err != nil {
			return false, err
		}
		if !success {
			return false, nil
		}
	}
	return true, nil
}

// runTaskCondition runs once per TaskCondition to check if contentions are successful
func (c *Checks) runTaskCondition(
	ctx context.Context,
	taskCondition teachv1alpha1.TaskCondition,
) (bool, error) {
	u, err := c.getConditionObject(ctx, taskCondition)
	if taskCondition.NotExists {
		if err != nil && client.IgnoreNotFound(err) == nil {
			return true, nil
		}
		return false, nil
	}
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			return false, nil
		}
		return false, err
	}

	success, err := c.runResourceConditions(taskCondition.ResourceCondition, *u)
	if err != nil {
		return false, err
	}

	if success {
		return true, nil
	}

	return false, nil
}

// runResourceConditions runs all ResourceConditions to the given object
// and returns true if all conditions are successful
func (c *Checks) runResourceConditions(
	resourceConditions []teachv1alpha1.ResourceCondition,
	item unstructured.Unstructured,
) (bool, error) {
	if len(resourceConditions) == 0 {
		return true, nil
	}
	parsed, _ := json.Marshal(item.Object)
	for _, resourceCondition := range resourceConditions {
		success, err := c.runResourceCondition(resourceCondition, string(parsed))
		if err != nil {
			return false, err
		}
		if !success {
			return false, nil
		}
	}
	return true, nil
}

// runResourceCondition run one condition to a json object and return true if condition is successful
func (c *Checks) runResourceCondition(
	resourceCondition teachv1alpha1.ResourceCondition,
	json string,
) (bool, error) {
	value := gjson.Get(json, resourceCondition.Field)

	switch resourceCondition.Operator {
	case "gt":
		checkValue, err := strconv.ParseInt(resourceCondition.Value, 10, 0)
		if err != nil {
			return false, err
		}
		if value.Int() > checkValue {
			return true, nil
		}
	case "nil":
		if !value.Exists() {
			return true, nil
		}
	case "notnil":
		if value.Exists() {
			return true, nil
		}
	case "lt":
		checkValue, err := strconv.ParseInt(resourceCondition.Value, 10, 0)
		if err != nil {
			return false, err
		}
		if value.Int() < checkValue {
			return true, nil
		}
	case "eq":
		if value.String() == resourceCondition.Value {
			return true, nil
		}
	case "neq":
		if value.String() != resourceCondition.Value {
			return true, nil
		}
	case "contains":
		if strings.Contains(value.String(), resourceCondition.Value) {
			return true, nil
		}
	default:
		return false, errors.New("invalid operator")
	}
	return false, nil
}

// getConditionObject returns the object for a TaskCondition as a unstructured object
func (c *Checks) getConditionObject(
	ctx context.Context,
	taskCondition teachv1alpha1.TaskCondition,
) (*unstructured.Unstructured, error) {
	u := unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   taskCondition.APIGroup,
		Version: taskCondition.APIVersion,
		Kind:    taskCondition.Kind,
	})

	err := c.Client.Get(ctx,
		client.ObjectKey{Name: taskCondition.Name, Namespace: taskCondition.Namespace},
		&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
