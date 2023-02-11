# traefik-jwt-internal
Generate JWT token for internal services

## Configuration

`alg` *(string)* [optional, default: `HS256`] - Hash algorithm, currently supported `HS256` and `HS512`

`secret` *(string)* - Signing secret, at least 32 characters long

`ttl` *(int64)* [optional, default: `120`] - Token expiry time in seconds

`header.name` *(string)* [optional, default: `Authorization`] - Header in which to write token

`header.prefix` *(string)* [optional, default: `Bearer`] - Header value prefix

`claims` *(string)* - Token claims JSON (parsed as Go template)

## Claims examples

### From header

```
{
    "sub": "{{ .Headers.Get "X-User-Id" }}"
}
```

### From JSON header

```
{
    {{ $header := .Headers.Get "X-User" }}
    {{ $user := unmarshalJson $header }}
    "sub": "{{ $user.id }}"
}
```
