<div align="center">
<img src="./img/logo.svg" alt="kubeteach logo" title="kubeteach" width="350px" align=""/>

# kubeteach

[![Go Report Card](https://goreportcard.com/badge/github.com/dergeberl/kubeteach)](https://goreportcard.com/report/github.com/dergeberl/kubeteach)
[![Licence](https://img.shields.io/github/license/dergeberl/kubeteach)](https://github.com/dergeberl/kubeteach/blob/main/LICENSE)
[![Latest release](https://img.shields.io/github/v/release/dergeberl/kubeteach?include_prereleases)](https://github.com/dergeberl/kubeteach/releases)
[![Coverage Status](https://coveralls.io/repos/github/dergeberl/kubeteach/badge.svg?branch=main)](https://coveralls.io/github/dergeberl/kubeteach?branch=main)
[![Test status](https://img.shields.io/github/workflow/status/dergeberl/kubeteach/tests/main?label=test)](https://github.com/dergeberl/kubeteach/actions?query=branch%3Amain++workflow%3Atests++)
[![Build status](https://img.shields.io/github/workflow/status/dergeberl/kubeteach/build/main?label=build)](https://github.com/dergeberl/kubeteach/actions?query=branch%3Amain++workflow%3Abuild++)

</div>


Kubeteach is an operator build with [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) to learn kubernetes in kubernetes.

Kubeteach was created because I wanted to improve my golang and kubernetes operator knowledge. I came across kubebuilder and had the idea to learn kubernetes in kubernetes.

The idea is to get tasks as a kubernetes objects (custom resource) to learn how to interact with kubernetes/kubectl while solving tasks in the kubernetes cluster.
Kubeteach checks whether a task has been completed successfully based on defined conditions.

:warning: **Note:** kubeteach is not designed to deploy it to a production cluster. I recommend to use [kind](https://kind.sigs.k8s.io/) to use kubeteach and learn kubernetes. There is **no** deletion job for created objects from exercises.

:construction: **Kubeteach is still in a very early stage of development, which means it may not always be stable and major api changes are necessary.**

## Learn kubernetes with kubeteach

:warning: Unfortunately, only a few tasks are currently available. More tasks will be added soon.

### Preparation

To install kubeteach you need a kubernetes cluster. I recommend [kind](https://kind.sigs.k8s.io/) as a local environment, checkout the [kind quick start](https://kind.sigs.k8s.io/docs/user/quick-start/).

You need also `kubectl` to interact with your cluster and `helm` to install kubeteach to your cluster. 
- Checkout the [install kubectl guide](https://kubernetes.io/de/docs/tasks/tools/install-kubectl/).
- Checkout the [install helm guide](https://helm.sh/docs/intro/install/).


### Installation

#### Add kubeteach helm repo

To be able to deploy kubeteach you need to add the kubeteach helm repo to your local helm installation. 

```bash
helm repo add kubeteach https://dergeberl.github.io/kubeteach-charts
```


#### Install kubeteach with ExerciseSet

To deploy kubeteach with an ExerciseSet you can select one of this [list](#list-of-exercisesets).

With the following command you can install kubeteach with an ExerciseSet to your cluster. (Change `<helm-chart` to your selected helm chart. For example `kubeteach/kubeteach-exerciseset1`) 
```bash
helm install exerciseset1 <helm-chart> --namespace exerciseset --create-namespace
```

:warning: Don't use the helm flag `--wait`, because some deployments won't get ready and the helm install command will fail.


#### Enable kubeteach dashboard

:warning: The dashboard is an experimental feature. DO NOT MAKE IT AVAILABLE VIA INTERNET! :warning:

To enable the dashboard you need to add  2 settings for the helm install command (see above):
```bash
--set kubeteach.dashboard.enabled=true --set kubeteach.webterminal.enabled=true
```

Example:
```bash
helm install exerciseset1 <helm-chart> --namespace exerciseset --set kubeteach.dashboard.enabled=true --set kubeteach.webterminal.enabled=true
...
You can use it with the following command (to forward a local port):
kubectl port-forward -n exerciseset service/kubeteach-core-dashboard 8080:80

Now you can access the dashboard via http://localhost:8080
Username: kubeteach
Password: <yourpassword>
```

The command will prompt a command (`kubectl port-forward`) and the credentials which are needed to log in into the dashboard.

### Update kubeteach

To update kubeteach you can run the following commands.
```bash
helm repo update
helm upgrade exerciseset1 <helm-chart> --namespace exerciseset
```

:warning: Don't use the helm flag `--wait`, because some deployments won't get ready and the helm install command will fail.


### Usage

You can get the tasks that should be performed with `kubectl get tasks -n exerciseset`

```bash
kubectl get tasks
NAME    TITLE                           DESCRIPTION                                                                                                          STATUS
task01   Create namespace                Create a new namespace with the name kubeteach                                                                           active
task02   Create pod                      Create a pod in namespace kubeteach, name it pod1 and use nginx:latest as image                                         pending
...
```

To get more information of one task you can use `kubectl describe task -n exerciseset <taskname>`

In some task you can find a `helpURL` and/or a `longDescription` with more information about this task.

```bash
kubectl describe task task01   
Name:         task01
Namespace:    default
Labels:       <none>
Annotations:  <none>
API Version:  kubeteach.geberl.io/v1alpha1
Kind:         Task
Metadata:
  Creation Timestamp:  2021-03-14T18:35:49Z
  Generation:          1
  Owner References:
    API Version:     kubeteach.geberl.io/v1alpha1
    Kind:            TaskDefinition
    Name:            task1
    UID:             21b8853d-11d4-4930-bdfd-ea3c945ae536
  Resource Version:  633
  UID:               d392614d-6a42-4d97-8500-ee29f3121674
Spec:
  Description:  Create a new namespace with the name kubeteach
  Title:        Create namespace
Status:
  State:  active
Events:
  Type    Reason  Age   From  Message
  ----    ------  ----  ----  -------
  Normal  Active  12m   Task  Task has no pre required task, task is now active
```

Now you can solve tasks by doing what's described in the task.

For example `task01`: `Create a new namespace with the name kubeteach`

```bash
kubectl create namespace kubeteach
```

A few seconds later the task state is changed to `successful`.

```bash
kubectl get task task01            
NAME    TITLE              DESCRIPTION                                  STATUS
task01   Create namespace   Create a Namespace with the name kubeteach   successful
```

The task state `pending` shows that another task must be successfully done before.

If you need help you can take a look into the solution folder of the exercise set you use (for example [dergeberl/kubeteach-charts/solutions/exerciseset1](https://github.com/dergeberl/kubeteach-charts/tree/main/solutions/exerciseset1))

**An update to a new status can take up to 5 seconds**

## List of ExerciseSets

| name | description | link | helm |
| --- | --- | --- | --- |
| kubeteach-exerciseset1 | Example ExerciseSet to try out kubeteach with basic tasks for first steps in kubernetes| [dergeberl/kubeteach-charts/charts/exerciseset1](https://github.com/dergeberl/kubeteach-charts/tree/main/charts/exerciseset1) | `kubeteach/kubeteach-exerciseset1` |

## How it works / How to write own exercises

### ExerciseSet (optional)

This is optional you can create directly `TaskDefinition`.

An `ExerciseSet` contains one or multiple `TaskDefinitions` to group them together and get some metadata from this `TaskDefinition`'s.

Each `spec.taskDefinitions` consists of a `name` (name of the `TaskDefinition` object) and a `taskDefinitionSpec` (spec of the `TaskDefinition`, see below).

#### Example

```yaml
apiVersion: kubeteach.geberl.io/v1alpha1
kind: ExerciseSet
metadata:
  name: set1
spec:
  taskDefinitions:
    - name: task01
      taskDefinitionSpec:
        taskSpec:
          title: "Create namespace"
          description: "Create a new namespace with the name kubeteach"
        taskCondition:
          - apiVersion: v1
            kind: Namespace
            name: "kubeteach"
        points: 5
    - name: task02
      taskDefinitionSpec:
        taskSpec:
          title: "Create namespace"
          description: "Create a new namespace with the name kubeteach2"
        taskCondition:
          - apiVersion: v1
            kind: Namespace
            name: "kubeteach2"
        points: 5
```

#### Status

The `ExerciseSet` status contains some metadata information of the tasks.

```yaml
...
status:
  numberOfActiveTasks: 2
  numberOfPendingTasks: 11
  numberOfSuccessfulTasks: 0
  numberOfTasks: 13
  numberOfTasksWithoutPoints: 0
  numberOfUnknownTasks: 0
  pointsAchieved: 0
  pointsTotal: 65
...
```

### TaskDefinition

A `TaskDefinition` describes a `Task` and conditions to check if the task is successful.

#### taskSpec

The `taskSpec` will be copied to the `Task` and is the object which is used for solving tasks. It should contain all information which are needed to solve the `Task`.

The following fields are available:
- `title` - title of the task
- `description` - description of the task which is shown by `kubectl get tasks`
- `longDescription` (optional) - longer description of the task which is shown by `kubectl describe tasks`
- `helpURL` (optional) - an url to more information about the topic in the task

#### points

`points` is an optional field which is only used if the `TaskDefinition` is created by an `ExerciseSet` to sum all points inside the `ExerciseSet`-status.

#### taskCondition

To check if the task is successful there is a list of `taskCondition`. 
Each `taskCondition` describes an object (apiVersion, kind and name) and contains a list of `resourceCondition` (see below).
If there is no `resourceCondition` the `taskCondition` is successful if the object exists.

To check if an object doesn't exist you can use `spec.taskConditions.notExists` and set it to true. In this case all `resourceCondition` are ignored for this `taskCondition` and this `taskCondition` is successful if the kubernetes object does not exist.

To depend on another task you can link a task as required with `spac.requiredTaskName`. This task will be in pending until the required task is successful. Be careful there is no check if the tasks can ever become active or are stuck in pending forever.

#### resourceCondition

Each `resourceCondition` contains a `field` which should be checked, an `operator` (see table below) and a `value`.

The `field` is a json path to find the field witch should be checked. The json path is based on [tidwall/gjson](https://github.com/tidwall/gjson). Check out the [gjson](https://github.com/tidwall/gjson) repository for the syntax. There is also an [online playground](https://gjson.dev/) for test and evaluate a json path.

| operators | description |
| --- | --- |
| eq | equal |
| neq | not equal |
| gt | greater than, value and field must be a number |
| lt | less than, value and field must be a number |
| nil | field is not set, value will be ignored |
| notnil | field is set, value will be ignored |
| contains | string is contained in field |

If there are multiple `taskCondition` and `resourceCondition` then **all** must be successful to complete the task.

#### Example

A simple example to check if a namespace is created:

```yaml
apiVersion: kubeteach.geberl.io/v1alpha1
kind: TaskDefinition
metadata:
  name: task1
spec:
  taskSpec:
    title: "Create namespace"
    description: "Create a new namespace with the name kubeteach"
  taskConditions:
    - apiVersion: v1
      kind: Namespace
      name: "kubeteach"
      resourceCondition:
        - field: "metadata.name"
          operator: "eq"
          value: "kubeteach"
```

Example, delete the created kubeteach namespace from task1 (`notExists` and `requiredTaskName`):

```yaml
apiVersion: kubeteach.geberl.io/v1alpha1
kind: TaskDefinition
metadata:
  name: task2
spec:
  taskSpec:
    title: "Delete namespace"
    description: "Delete the namespace kubeteach that is created in task1"
  requiredTaskName: task1
  taskConditions:
    - apiVersion: v1
      kind: Namespace
      name: "kubeteach"
      notExists: true
```

## Contribution

Feedback and PRs are very welcome. For major changes, please open an issue beforehand to clarify if this is in line with the project and to avoid unnecessary work.

New exercises or/and exercise sets are highly welcome, if you have ideas feel free to open a PR or issue.