package check

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

type CheckController struct {
	Client     client.Client
	Log        logr.Logger
	CheckItems []kubeteachv1.TaskCondition
}

func (c *CheckController) ApplyChecks(ctx context.Context) (bool, error) {
	if len(c.CheckItems) < 1 {
		return false, errors.New("no checks to apply")
	}
	for _, checkItem := range c.CheckItems {
		success, err := c.runCheckItem(ctx, checkItem)
		if err != nil {
			return false, err
		}
		if !success {
			return false, nil
		}
	}
	return true, nil
}

func (c *CheckController) runCheckItem(ctx context.Context, checkItem kubeteachv1.TaskCondition) (bool, error) {
	u, err := c.getItemList(ctx, checkItem)
	if err != nil {
		return false, err
	}

	var successfulItem int
	for _, item := range u.Items {
		success, err := c.runChecks(checkItem.ResourceCondition, item)
		if err != nil {
			return false, err
		}
		if success && !checkItem.MatchAll {
			return true, nil
		}
		if success {
			successfulItem++
		}
	}
	if checkItem.MatchAll && successfulItem == len(u.Items) {
		return true, nil
	}
	return false, nil
}

func (c *CheckController) runChecks(checks []kubeteachv1.ResourceCondition, item unstructured.Unstructured) (bool, error) {
	parsed, _ := json.Marshal(item.Object)
	for _, check := range checks {
		success, err := c.runCheck(check, string(parsed))
		if err != nil {
			return false, err
		}
		if !success {
			return false, nil
		}
	}
	return true, nil
}

func (c *CheckController) runCheck(check kubeteachv1.ResourceCondition, json string) (bool, error) {
	value := gjson.Get(json, check.Field)

	switch check.Operator {
	case "gt":
		checkValue, err := strconv.ParseInt(check.Value, 10, 0)
		if err != nil {
			return false, err
		}
		if value.Int() > checkValue {
			return true, nil
		}
	case "lt":
		checkValue, err := strconv.ParseInt(check.Value, 10, 0)
		if err != nil {
			return false, err
		}
		if value.Int() < checkValue {
			return true, nil
		}
	case "eq":
		if value.String() == check.Value {
			return true, nil
		}
	case "neq":
		if value.String() != check.Value {
			return true, nil
		}
	case "contains":
		if strings.Contains(value.String(), check.Value) {
			return true, nil
		}
	}
	return false, nil
}

func (c *CheckController) getItemList(ctx context.Context, checkItem kubeteachv1.TaskCondition) (*unstructured.UnstructuredList, error) {
	u := unstructured.UnstructuredList{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   checkItem.ApiGroup,
		Version: checkItem.ApiVersion,
		Kind:    checkItem.Kind,
	})

	err := c.Client.List(ctx, &u, &client.ListOptions{})
	if err != nil && client.IgnoreNotFound(err) != nil {
		return nil, err
	}

	return &u, nil
}
