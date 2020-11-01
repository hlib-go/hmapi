module hapi

go 1.14

require (
	github.com/go-redis/redis/v8 v8.2.3
	github.com/sirupsen/logrus v1.7.0
	hgenid v0.0.0
	hhttp v0.0.0
	hredis v0.0.0
)

replace (
	hgenid v0.0.0 => github.com/hlib-go/hgenid v0.0.0-20201010154643-5bb5cd815855
	hhttp v0.0.0 => github.com/hlib-go/hhttp v0.0.0-20201010155030-030f71f1769e
	hredis v0.0.0 => github.com/hlib-go/hredis v0.0.0-20201009152956-4faca4c92a1b
)
