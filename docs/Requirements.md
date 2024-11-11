# Requirements

## Test task requirements

### Description

Go k8 cache service is an application to store JSON-objects 'in-memory' key-based storage with an ability to restore data from the disk.
It is required to implement a microservice that should be deployed in k8s.

There should be a funtionality to set TTL (time-to-live) for the object.
JSON-object will set to Cache using REST requests.

---


### REST API

#### Store data

```
PUT /objects/:objectId
Body(required): JSON-object
Header(optional): "expires_at"
```

_objectId_ - an object id in cache
_expires_at_ - an optional header to set object TTL (format: 'YYYY-MM-DD hh:mm:ss')

#### Get data

```
GET /objects/:objectId
```

_objectId_ - an object id in cache

Response example:

```json
{
    "server": {
        "id": "12345678-90ab-cdef-1234-567890abcdef",
        "name": "example-vm",
        "status": "ACTIVE",
        "created": "2023-08-15T10:00:00Z",
        "updated": "2023-08-15T12:00:00Z",
        "flavor": {
            "id": "flavor-id",
            "name": "m1.small",
            "vcpus": 1,
            "ram": 2048,
            "disk": 20
        },
        "image": {
            "id": "image-id",
            "name": "ubuntu-20.04"
        },
        "hostId": "host-id",
        "tenant_id": "tenant-id",
        "user_id": "user-id",
    }
}
```

---

### Service health and readiness checks

Go Cache service supports health and readiness checking

#### Liveness check

```
GET /probes/liveness
```

Respond with `200 Ok` in case of application liveness

#### Readiness check

```
GET /probes/readiness
```

Respond with:

- `200 OK` In case of application ready to receive traffic, and data has been uploded from storage to memory cache
- `503 Service Unavailable` in case of application NOT ready to receive traffic, that means data transfer to memory cache is in progress

---

### Metrics

Go Cache service supports prometheus metrics exposal.

#### Metrics

```
GET /metrics
```

Service respond with standart prometheus metrics (e.g. goroutines stats, http requests count filtered by codes and response times) and custom metrics:

```
# HELP memory_cache_items_count The current objects count in memory cache
# TYPE memory_cache_items_count gauge
memory_cache_items_count 1660
```

This metric shows a current object count stores in Memory Cache
