package redisMem

import (
	"context"
	"github.com/go-redis/redis"
	"go.opentelemetry.io/otel"
	"log"
	"strconv"
	"sync"
	"therealbroker/Internal/modules"
	broker2 "therealbroker/package/broker"
	"time"
)

const name = "redis"

type handler struct {
	pipeline  redis.Pipeliner
	waitermng modules.WaiterManager

	locker sync.Mutex
}

var clt *redis.Client

var pl handler

func InitRedisClient() {
	clt = redis.NewClient(&redis.Options{Addr: "redis:6379", Password: "", DB: 0})

	_, err := clt.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	if ok, _ := clt.Exists("ID").Result(); ok == 0 {
		clt.Set("ID", 0, 0)
	}

	pl = handler{
		pipeline:  clt.Pipeline(),
		waitermng: *modules.NewWaiterManager(),
		locker:    sync.Mutex{},
	}

	//go Scheduler(executePipeline, 15*time.Millisecond)
}

func StoreData(ctx context.Context, id int, msg broker2.Message, subject string, creation time.Time) {
	_, span := otel.Tracer(name).Start(ctx, "redisDB")
	defer span.End()

	storeMessage(id, msg)
	storeSubject(id, subject)
	storeCreation(id, creation)

	pl.locker.Lock()
	ch := pl.waitermng.RegisterWaiter()
	pl.locker.Unlock()
	<-ch
}

func executePipeline() {
	if _, err := pl.pipeline.Exec(); err != nil {
		log.Fatal("failed to write to redis")
	}

	pl.locker.Lock()
	pl.waitermng.FreeWaiters()
	pl.locker.Unlock()
}

func storeMessage(id int, msg broker2.Message) {
	prc := strconv.Itoa(int(msg.Expiration)) + "&" + msg.Body
	identifier := "m:" + strconv.Itoa(id)
	pl.pipeline.Set(identifier, prc, 0)
}

func storeSubject(id int, subject string) {
	identifier := "h:" + strconv.Itoa(id)
	pl.pipeline.Set(identifier, subject, 0)
}

func storeCreation(id int, creation time.Time) {
	identifier := "c:" + strconv.Itoa(id)
	pl.pipeline.Set(identifier, creation.Format(time.RFC1123), 0)
}

func GetNewId() int {
	id, err := clt.Incr("ID").Result()
	if err != nil {
		log.Fatal("Failed to get new Id with error: ", err)
	}
	return int(id)
}
