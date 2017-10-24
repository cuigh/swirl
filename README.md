# SWIRL

[![Swirl](https://goreportcard.com/badge/cuigh/swirl)](https://goreportcard.com/report/cuigh/swirl)

Swirl is a web management tool for Docker, focused on swarm cluster.

## Features

* Swarm components management
* Image and container management
* Compose management with deployment support
* LDAP authentication support
* Full permission control based on RBAC model
* Scale out as you want
* Multiple language support
* And more...

## Snapshots

### Dashboard

![Dashboard](docs/images/dashboard.png)

### Service list

![Service list](docs/images/service-list.png)

### Compose list

![Compose list](docs/images/compose-list.png)

### Role editing

![Role editing](docs/images/role-edit.png)

### Settings

![Setting](docs/images/setting.png)

## Configuration

### With config file

All options can be set with `config/app.conf`.

```xml
<config>
    <app>
        <add key="name" value="swirl"/>
    </app>
    <web>
        <add key="address" value=":8001"/>
        <!-- default authorize mode, valid options: *(everyone)/?(login user)/!(authorized explicitly) -->
        <add key="authorize_mode" value="?"/>
    </web>
    <swirl>
        <!-- optional -->
        <add key="docker_endpoint" value="tcp://docker-proxy:2375"/>
        <!-- optional, valid options: mongo -->
        <add key="db_type" value="mongo"/>
        <!-- required, database connection string, must match with db.type option -->
        <add key="db_address" value="localhost:27017/swirl"/>
    </swirl>
</config>
```

### With environment variables

Only three main options can be set by environment variables for now.

| Name            | Value                                           |
| --------------- | ------------------------------------------------|
| DB_TYPE         | mongo                                           |
| DB_ADDRESS      | localhost:27017/swirl                           |
| DOCKER_ENDPOINT | tcp://docker-proxy:2375                         |

### With swarm config

Docker support mounting configuration file through swarm from v17.06, so you can store your config in swarm and mount it to your program.

## Deployment

### Stand alone

Just copy the swirl binary and config/assets/views directories to the host, and run it.

```bash
./swirl
```

### Docker

```bash
docker run -d -p 8001:8001 \
    --mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
    -e DB_TYPE=mongo \
    -e DB_ADDRESS=localhost:27017/swirl \
    --name=swirl \
    cuigh/swirl
```

### Docker swarm

```bash
docker service create \
  --name=swirl \
  --publish=8001:8001/tcp \
  --env DB_ADDRESS=localhost:27017/swirl \
  --constraint=node.role==manager \
  --mount=type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
  cuigh/swirl
```

## Build

**Swirl** use `dep` as dependency management tool(coming soon).

## License

This product is licensed to you under the MIT License. You may not use this product except in compliance with the License. See LICENSE and NOTICE for more information.