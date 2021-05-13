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
apiVersion: caffeine.lu
kind: SecretManager
metadata:
  name: the-secret
secrets:
  - some_database_password
  - another_secret
```

This will generate the following secret:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: the-secret
type: Opaque
data:
  some_database_password: ...
  another_secret: ...
```
