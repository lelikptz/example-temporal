# Example of distributed workflow with temporal.io

#### Run temporal

```bash 
brew install temporal 
temporal server start-dev
```

#### Run distributed temporal workflow

```bash
go run cmd/workflow/main.go
```

#### Run microservice with "create-order" activity

```bash
go run cmd/create-order/main.go
```

#### Run microservice with "send-order" activity

```bash
go run cmd/send-order/main.go
```

#### Run microservice with task producer

```bash
go run cmd/task-producer/main.go
```

### Enjoy!
