build:
	mkdir -p bin/github.com/julienp/secretmanager/
	go build -o bin/github.com/julienp/secretmanager/SecretManager cmd/main.go

example:
	KUSTOMIZE_PLUGIN_HOME=`pwd`/bin kustomize build --enable_alpha_plugins ./example

.PHONY: example
