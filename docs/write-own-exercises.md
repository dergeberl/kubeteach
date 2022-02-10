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