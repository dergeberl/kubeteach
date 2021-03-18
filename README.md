# kubeteach

Kubeteach is an operator build with [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) to learn kubernetes in kubernetes.

**Note:** kubeteach is not designed to deploy it to a production cluster. I recommend to use [kind](https://kind.sigs.k8s.io/) to use kubeteach and learn kubernetes.

## Installation

To install kubeteach you first need a kubernetes cluster. Checkout the [kind quick start](https://kind.sigs.k8s.io/docs/user/quick-start/)

At first, you deploy the operator itself by applying the deployment file in the `deployment` folder.
```bash
git clone dergeberl/kubeteach
cd kubeteach
kubectl apply -f deployment/
```

Now you can  deploy a set of exercises to your cluster.

```bash
kubectl apply -f exercises/set1/
```

## Usage

You can get the tasks that should be performed with `kubectl`

```bash
kubectl get tasks
NAME    TITLE                           DESCRIPTION                                                                                                          STATUS
task1   Create namespace                Create a new namespace with the name kubeteach                                                                           active
task2   Create pod                      Create a pod in namespace kubeteach, name it pod1 and use nginx:latest as image                                         pending
...
```

To get more information of one task you can use `kubectl describe task <taskname>`


```bash
kubectl describe task task1   
Name:         task1
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

For example `task1`: Create a new namespace with the name kubeteach

```bash
kubectl create namespace kubeteach
```

A few seconds later the task state is changed to `successful`

```bash
kubectl get task task1            
NAME    TITLE              DESCRIPTION                                  STATUS
task1   Create namespace   Create a Namespace with the name kubeteach   successful
```

If a task is in state `pending` you have to solve another task before.


## how it works

In the `exercise` folder is a set of `taskdefinitions` which describe a `task` and conditions to check if the task is successful.

To check if the task is successful there is a list of `taskCondition`. 
Each `taskCondition` describes an object (apiVersion, kind and name) and contains a list of `resourceCondition`. 
Each `resourceCondition` contains a field which should be checked, an operator (see table below) and a value. 
The field is a json path to find the field witch should be checked (it is based on [tidwall/gjson](https://github.com/tidwall/gjson)).


| operators | description |
| --- | --- |
| eq | equal |
| neq | not equal |
| gt | greater than, value and field must be a number |
| lt | less than, value and field must be a number |
| nil | field is not set, value will be ignored |
| notnil | field is set, value will be ignored |
| contains | string is contained in field |

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

To check if an object doesn't exist you can use `spec.taskdefiniton.notExists` and set it to true. To depend on another task you can link a task as required with `spac.requiredTaskName`. This task will be in pending until the required task is successful. 

Example, delete the created kubeteach namespace from task1:

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

Feel free to open an issue or PR if you miss a feature or found a bug.

New exercises or exercise sets are highly welcome.