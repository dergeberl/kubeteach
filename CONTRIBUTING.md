# Contributing

Feedback and PRs are very welcome. For major changes, please open an issue beforehand to clarify if this is in line with the project and to avoid unnecessary work.

## Test setup

### `kubeteach` controller (go)

1. Export your `KUBECONFIG`-Environment variable
2. Install the CRDs with `make install`
3. Run the controller with `make run` 
4. You can use `.github/kind-test-taskdefinitions.yaml` as example `ExerciseSet`

## `kubeteach` dashboard 

### Test setup

1. Export your `KUBECONFIG`-Environment variable
2. Install the CRDs with `make install`
3. Run the controller with `make run`
4. Open in a second terminal the `dashboard` folder
5. Run `npm start`

To see also the Terminal window you can run an example pod:
```shell
kubectl run --image ghcr.io/dergeberl/kubeteach-webterminal --env GOTTY_PATH="/shell/" --env GOTTY_WS_ORIGIN=".*" kubeteach-webterminal
kubectl port-forward pod/kubeteach-webterminal 8091:8080
```

> This Terminal does not work, due to missing RBAC Roles.

## Open a PR

1. Make sure you have the correct Go version installed
2. Fork this repository
3. Run `make test` to check if everything works
4. Open a PR and make sure the CI pipelines succeeds
5. Wait until a maintainer will review your code
6. After approval your PR will get merged
7. Thank you! :heart:
