# KSecretManager

Kustomize plugin to load secrets from [Secret Manager](https://cloud.google.com/secret-manager)

Install the plugin where kustomize can find it, see the [documentation](https://kubectl.docs.kubernetes.io/guides/extending_kustomize/#placement).

Add a the generator to your kustomization.yaml:

```yaml
resources:
  - ...
generators:
  - my-secret.yaml
```

Create `my-secret.yaml`:

```yaml
apiVersion: github.com/julienp
kind: SecretManager
metadata:
  name: the-secret
secrets:
  - name: some_database_password # The name as specified in Google Secrets manager
  - name: another_secret
    key: ANOTHER_SECRET # optional, the key to use in the k8s Secret
```

Provide `PROJECT_ID` as env variable when running kustomize:

```bash
make build
PROJECT_ID=<my project id> KUSTOMIZE_PLUGIN_HOME=`pwd`/bin kustomize build --enable_alpha_plugins ./example
```

This will generate the following secret:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: the-secret
type: Opaque
data:
  some_database_password: <value>
  ANOTHER_SECRET: <value>
```
