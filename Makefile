build:
	mkdir -p bin/caffeine.lu/secretmanager/
	go build -o bin/caffeine.lu/secretmanager/SecretManager cmd/main.go

example:
	KUSTOMIZE_PLUGIN_HOME=`pwd`/bin kustomize build --enable_alpha_plugins ./example

.PHONY: example
