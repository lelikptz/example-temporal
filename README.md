# Example of distributed workflow with temporal.io

#### Run temporal

```bash 
brew install temporal 
temporal server start-dev
```

#### Run distributed temporal workflow

```bash
go run cmd/order-management-system/*.go
```

#### Run microservice with "cancel-order-service" activity

```bash
go run cmd/cancel-order-service/*.go

```
#### Run microservice with "create-order-service" activity

```bash
go run cmd/create-order-service/*.go
```

#### Run microservice with "notification-service" activity

```bash
go run cmd/notification-service/*.go
```

#### Run microservice with "send-order-service" activity

```bash
go run cmd/send-order-service/*.go
```

#### Run microservice with task producer

```bash
go run cmd/task-producer/*.go
```

### Enjoy!
