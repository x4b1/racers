module github.com/xabi93/racers

go 1.15

replace golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/DATA-DOG/go-txdb v0.1.3
	github.com/caarlos0/env/v6 v6.3.0
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/cockroachdb/errors v1.7.5
	github.com/containerd/continuity v0.0.0-20200928162600-f2cc35102c2a // indirect
	github.com/facebook/ent v0.4.4-0.20201006132909-48362e79cdd9
	github.com/facebookincubator/ent-contrib v0.0.0-20201006141524-9617dddb71dc
	github.com/go-kit/kit v0.8.0
	github.com/golang-migrate/migrate/v4 v4.13.0
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.7.4
	github.com/hashicorp/go-multierror v1.1.0
	github.com/kataras/golog v0.0.9
	github.com/kr/pretty v0.1.0
	github.com/lib/pq v1.8.0
	github.com/matryer/is v1.4.0
	github.com/ory/dockertest/v3 v3.6.0
	github.com/prometheus/client_golang v0.9.3
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/vektah/gqlparser/v2 v2.1.0
	github.com/vmihailenco/msgpack/v5 v5.0.0-beta.1
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20201009032441-dbdefad45b89 // indirect
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	golang.org/x/sys v0.0.0-20201009025420-dfb3f7c4e634 // indirect
	gorm.io/driver/postgres v1.0.2
	gorm.io/gorm v1.20.2
)
