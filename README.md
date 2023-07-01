# Introduction

A simple go http server that uses [libdns](https://github.com/libdns/libdns)' [leaseweb implementation](https://github.com/libdns/leaseweb) to update a record based on the `RemoteAddr` of the request.

It is very basic. It was written for use with the Dynamic DNS feature of a ZYXEL VMG8825-T50; whose are provided by T-Mobile Netherlands to ADLS users. Which seems to be some version of [ez-ipupdate](https://linux.die.net/man/8/ez-ipupdate).

# How

Listens on port 80 and will attempt to update the A record for the given host in the `hostname` query parameter to a value equal to: (in order) `X-Real-IP` header, `X-Forwarder-For` header and the requests `.RemoteAddr`.

The request is expected to provide a basic auth header equal to whaterver `DYNDNS_USERNAME` and `DYNDNS_PASSWORD` are set to.

To achieve this a Leaseweb API key is needed. Also, cause I'm lazy and don't want to parse the domain in any way the `DYNDNS_ZONE` must be set and the `hostname` query param must match. For for example when `DYNDNS_ZONE=example.com` the `hostname` query param could be `my-sub.exmaple.com`.

# Deployment

One can run it with docker compose.

Example compose deployment:

```env
DYNDNS_USERNAME=
DYNDNS_PASSWORD=
DYNDNS_ZONE=
DYNDNS_LEASEWEB_API_KEY=
```

```yml
version: '3.3'

services:
  app:
    image: justinvdk/dyndns:latest
    restart: unless-stopped
    ports:
      - '8090:80/tcp'
    environment:
      DYNDNS_USERNAME: "${DYNDNS_USERNAME}"
      DYNDNS_PASSWORD: "${DYNDNS_PASSWORD}"
      DYNDNS_ZONE: "${DYNDNS_ZONE}"
      DYNDNS_LEASEWEB_API_KEY: "${DYNDNS_LEASEWEB_API_KEY}"
```

# Client usage

Using cURL could be something like this:

```bash
DYNDNS_USERNAME=<your username>
DYNDNS_PASSWORD=<your password>

curl -u "${DYNDNS_USERNAME}:${DYNDNS_PASSWORD}" http://server.example.com/path?hostname=home-ip.example.com
```
