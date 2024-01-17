# Virtualhost

## Name

*virtualhost* - This plugin allows to resolve docker containers hostnames.

## Description

This plugin use the docker api to inspect docker conntainers and fetch each containers environment variable <code>VIRTAUL_HOST</code> if the variable exists.
The plugin doesn't handle port numbers, so a reverse proxy such as [nginx-proxy](https://github.com/nginx-proxy/nginx-proxy) should act as frontend.


## Syntax

```
virtualhost 192.168.0.100
```

## Metrics

If monitoring is enabled (via the *prometheus* plugin) then the following metrics are exported:
* `coredns_virtualhost_hostname_count{hostname} - Counter of hostname responses`

## Examples

The IP address(es) must point to the host running Docker.

With IPv4
```
example.com {
    virtualhost 192.168.0.100
}
```

With IPv6
```
example.com {
    virtualhost fe80::8770:87a8:3d30:84da
}
```

With IPv4 and IPv6
```
example.com {
    virtualhost 192.168.0.100 fe80::8770:87a8:3d30:84da
}
```

