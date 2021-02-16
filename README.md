# Rabbitmq testapp
Testapp to testing rabbitmq with different configs

Install rabbitmq with docker

```docker run -d --hostname my-rabbit --name some-rabbit-mng rabbitmq:3-management```

Sign in rabbitmq-management

```user 'guest' password 'guest'```

Manage rabbitmq. Select rabbitmq container

```docker ps```

Manage

```docker exec -it <container name> /bin/bash```

Checkout ```main``` or ```another``` branch



