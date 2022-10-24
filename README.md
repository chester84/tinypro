# tinypro
base on beego v2.0

# `sudo vim /etc/hosts`

## develop host
```
127.0.0.1       mysql.db.rds
127.0.0.1       cache.redis
127.0.0.1       storage.redis
127.0.0.1       es.dev
```

### develop env

```shel
vim ~/.bashrc

export GO111MODULE=on
export GOPROXY="https://goproxy.io"
```

## local develop

1. download `GeoLite2-City.mmdb` from `https://dev.maxmind.com/geoip/geoip2/geolite2/`, place it under api/data/
1. set up the `hosts` above
1. reference `docs/db/init.sql` to init mysql dev account
1. execute `docs/db/db-schema.sql`, init db

1. run api

```shell
cd api
go run main.go
```

2. docker for local

```shell
docker run --rm --name api-8915 -p 8915:8915 -e DEV_HOST_IP=devIP -d tinypro:TAG
```

# beego docs

https://beego.me/docs/intro/

