# kubectl-mo
Executes `kubectl get ...` command for multiple namespaces to get any resource you want.

## Feature 1: Execute `kubectl get ...` in current and last namespace
### Pre-requisites
- [kubens](https://github.com/ahmetb/kubectx#kubectx--kubens-power-tools-for-kubectl) (`kubectl-mo` uses information stored by `kubens` to figure out the last used namespace)

For example, let's say your current namespace is `kube-system`:
```
$ kubectl config view --minify -o jsonpath='{..namespace}'
kube-system
```
You change the namespace to `default` using [`kubens`](https://github.com/ahmetb/kubectx#kubectx--kubens-power-tools-for-kubectl):
```
$ kubens default
Context "k3d-dynamic-corefile" modified.
Active namespace is "default".
```
Now, if you want to get *any resource* from the current and the last namespace, you can use `kubectl-mo` like this:
```
$ kubectl mo get configmaps
executing command /home/suraj/bin/kubectl -n default get configmaps
NAME               DATA   AGE
kube-root-ca.crt   1      33d

executing command /home/suraj/bin/kubectl -n kube-system get configmaps
NAME                                 DATA   AGE
extension-apiserver-authentication   6      33d
cluster-dns                          2      33d
local-path-config                    4      33d
chart-content-traefik                0      33d
chart-values-traefik                 1      33d
chart-content-traefik-crd            0      33d
chart-values-traefik-crd             0      33d
kube-root-ca.crt                     1      33d
coredns                              2      33d
```
```
$ kubectl mo get pods
executing command /home/suraj/bin/kubectl -n default get po
NAME    READY   STATUS    RESTARTS       AGE
nginx   1/1     Running   14 (53m ago)   11d

executing command /home/suraj/bin/kubectl -n kube-system get po
NAME                                      READY   STATUS      RESTARTS        AGE
helm-install-traefik-crd-vplbc            0/1     Completed   0               33d
helm-install-traefik-q26gv                0/1     Completed   1               33d
local-path-provisioner-7b7dc8d6f5-q8k6x   1/1     Running     51 (53m ago)    33d
coredns-b96499967-tmswr                   1/1     Running     51 (53m ago)    33d
svclb-traefik-98312b71-66crq              2/2     Running     121 (53m ago)   33d
traefik-7cd4fcff68-dt82t                  1/1     Running     51 (53m ago)    33d
metrics-server-668d979685-tbblc           1/1     Running     51 (53m ago)    33d
```

In general, you can use:
```
$ kubectl mo get <any-resource> <flags>
```
To execute any `kubectl get ...` command with flags in your current and last namespace.

## Feature 2: Execute `kubectl get ...` in multiple namespaces
### Pre-requisites
- none (you can use this feature without `kubens`) 

Instead of running `kubectl get ...` in the current and the last namespace, you want to run `kubectl get ...` in arbitrary multiple namespaces, you can do:
```
$ kubectl mo get configmaps  -ns=kube-public,kube-node-lease
executing command /home/suraj/bin/kubectl -n kube-public get configmaps
NAME               DATA   AGE
kube-root-ca.crt   1      33d

executing command /home/suraj/bin/kubectl -n kube-node-lease get configmaps
NAME               DATA   AGE
kube-root-ca.crt   1      33d
```

In general,
```
$ kubectl mo get <any-resource> <flags> -ns=ns1,ns2...
```
If you add `-ns` flag `kubectl-mo` doesn't execute `kubectl get ...` in the current and the last namespace. 

## Good to know
- `kubectl-mo` uses the `kubectl` binary available on the user's machine. This is intentional. We want `kubectl-mo` to be compatible with different versions of `kubectl`.
## What problem is `kubectl-mo` trying to solve?
- When you are working with kubernetes, you often want to see resources in multiple namespaces but `kubectl` doesn't give you that functionality. To be specific, you often want one of the following:
  1. Look for resources in your current and the last used namespace to avoid constantly switching namespaces
  2. Look for resources in your current and last N namespaces
  3. Look for resources in arbitrary multiple namespaces
- There have been issues filed around this but there doesn't seem to be any progress on solving this problem in the upstream `kubectl` code:
  - https://github.com/kubernetes/kubectl/issues/763
  - https://github.com/kubernetes/kubernetes/issues/52326
- There have been attempts to solve this problem using a plugin:
  - https://github.com/kubernetes/kubernetes/issues/52326#issuecomment-1157022134
  - Above implementation solves 3 but doesn't solve 1 and 2
  - `kubectl-mo` solves 1 and 3 but doesn't solve 2 yet. 

## Cons
- If you have too many resources in your namespaces, using `kubectl-mo get` might take time to execute because the underlying `kubectl get ...` takes time to execute. 

## Feedback
We'd love to know what you think about the plugin. Feel free to create an issue to give us feedback, suggestions, file bugs or ask for new features. 

## Future scope
- It would be nice to have a way to let `kubectl-mo` figure out the last 5 namespaces that the user used. This includes the namespaces user switched to (using say `kubens`) and also the namespaces user used in their `kubectl` command using `--namespace` or `-n` flag. `kubectl-mo` can then automatically query for resources in those last 5 namespaces.
- It would be nice if the user can specify namespaces once and `kubectl-mo` uses the namespaces in subsequent commands. For example:
```
$ kubectl mo get po -ns=default,kube-system
# kubectl mo saves `default` and `kube-system` namespaces
$ kubectl-mo get po # no need to specify `-ns` again

$ kubectl mo get po -ns=dev,prod
# kubectl-mo overwrites last saved namespaces i.e., `default` and `kube-system` with `dev` and `prod`
$ kubectl mo get po # no need to specify `-ns` again
``` 

**Note: Future scope will be worked upon based on user interest**