# How to talk with your Kubernetes cluster without kubectl

The things would be much simpler, but unfortunately, the current state of `client-go` versioning is that you need to go get a version of client-go that matches the version you got the example from.
And you need to do that manually, by hand, without using `go mod`.
The latest `semver` tags on the repo resolve to old versions (e.g `v11.0.0+incompatible`) lacking `go.mod` info, which makes `go fetch` master of all transitive dependencies, which are not compatible with the old version of client go.
As a result, go modules do not let the maintainers of `client-go` tag further major versions without renaming the module.
They are investigating doing that, but that is a larger change that has not been settled on yet.

> Whole story can be found [here](https://github.com/kubernetes/client-go/issues/757)

So, for the time being, here is how you can use it:

1. Find the version of the API Server your Kubernetes cluster is running

```zsh
kubectl version --short | awk -F "Server Version: " '{ print $2 }' | tail -n 1 
# In my case the output is: v1.16.2
```

2. Find the nearest tag of `client-go` to your API Server version:

```zsh
go list -m -versions k8s.io/client-go | tr " " "\n" | grep <$VERSION>

v0.16.4         # <--- We will pick this one, since there is no 16.2
v0.16.5-beta.0
v0.16.5-beta.1
v0.16.5
v0.16.6-beta.0
v0.16.6
v0.16.7-beta.0
v0.16.7
v0.16.8-beta.0
```

```zsh
export tag="v0.16.4"
```

3. Fetch the `main.go` from the branch/tag.

```zsh
wget "https://raw.githubusercontent.com/kubernetes/client-go/$tag/examples/out-of-cluster-client-configuration/main.go"
```

4. Initialize the modules

```zsh
go mod init # if you are inside the $GOPATH
go mod init github.com/<$USERNAME>/<$REPO> # if you are outside of the $GOPATH
```

5. Configure the modules

```zsh
go mod tidy
```

6. Replace the tags of k8s.io/{api,apimachinery,client-go} with the <$tag>

```zsh
for dep in apimachinery api client-go; do sed -i "s/$dep .*/$dep $tag/" go.mod; done; rm go.mod-e
```

7. Build and run

```zsh
go build -o app
./app
```
