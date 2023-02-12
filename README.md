# traefik-jwt-internal
Generate JWT token for internal services

## Configuration

| Key | Type | Required | Default | Description |
| :-- | :-- | :-: | :-- | :-- |
| `alg` | string | | `HS256` | Hash algorithm, currently supported `HS256` and `HS512` |
| `claims` | string | :white_check_mark: | | Token claims JSON (parsed as Go template) |
| `header.name` | string | | `Authorization` | Header in which to write token |
| `header.prefix` | string | | `Bearer` | Header value prefix |
| `secret` | string | :white_check_mark: | | Signing secret, at least 32 characters long |
| `ttl` | int64 | | `120` | Token expiry time in seconds |

## Claims examples

### From header

```
{
    "sub": "{{ .Header.Get "X-User-Id" }}"
}
```

### From JSON header

```
{
    {{ $header := .Header.Get "X-User" }}
    {{ $user := unmarshalJson $header }}
    "sub": "{{ $user.id }}"
}
```
