# How to check after updating LB based on nginx

## processes

* is nginx running?

  * In this box, worker is only `1`

```sh
test `ps aux | grep 'nginx: master process' | grep -v grep | wc -l` -eq 1
test `ps aux | grep 'nginx: worker process' | grep -v grep | wc -l` -eq 1
```

## configs

* does nginx configs exist?

  * TODO: Check other configs

```sh
ls /usr/local/etc/nginx/nginx.conf
```

## logs

* do access/error logs exist?

```sh
ls /usr/local/var/log/nginx/access.log
ls /usr/local/var/log/nginx/error.log
```

