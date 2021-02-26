package condition

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-logr/logr"
	"github.com/tidwall/gjson"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubeteachv1 "kubeteach/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"strings"
)

type ConditionChecks struct {
	Client client.Client
	Log    logr.Logger
}

func (c *ConditionChecks) ApplyChecks(ctx context.Context, taskConditions []kubeteachv1.TaskCondition) (bool, error) {
	if len(taskConditions) < 1 {
		return false, errors.New("no checks to apply")
	}
	for _, taskCondition := range taskConditions {
		success, err := c.runCheckItem(ctx, taskCondition)
		if err != nil {
			return false, err
		}
		if !success {
			return false, nil
		}
	}
	return true, nil
}

func (c *ConditionChecks) runCheckItem(ctx context.Context, taskCondition kubeteachv1.TaskCondition) (bool, error) {
	u, err := c.getItemList(ctx, taskCondition)
	if err != nil {
		return false, err
	}

	var successfulItem int
	for _, item := range u.Items {
		success, err := c.runChecks(taskCondition.ResourceCondition, item)
		if err != nil {
			return false, err
		}
		if success && !taskCondition.MatchAll {
			return true, nil
		}
		if success {
			successfulItem++
		}
	}
	if taskCondition.MatchAll && successfulItem == len(u.Items) {
		return true, nil
	}
	return false, nil
}

func (c *ConditionChecks) runChecks(resourceConditions []kubeteachv1.ResourceCondition, item unstructured.Unstructured) (bool, error) {
	parsed, _ := json.Marshal(item.Object)
	for _, resourceCondition := range resourceConditions {
		success, err := c.runCheck(resourceCondition, string(parsed))
		if err != nil {
			return false, err
		}
		if !success {
			return false, nil
		}
	}
	return true, nil
}

func (c *ConditionChecks) runCheck(resourceCondition kubeteachv1.ResourceCondition, json string) (bool, error) {
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

func (c *ConditionChecks) getItemList(ctx context.Context, taskCondition kubeteachv1.TaskCondition) (*unstructured.UnstructuredList, error) {
	u := unstructured.UnstructuredList{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   taskCondition.ApiGroup,
		Version: taskCondition.ApiVersion,
		Kind:    taskCondition.Kind,
	})

	err := c.Client.List(ctx, &u, &client.ListOptions{})
	if err != nil && client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	return &u, nil
}
