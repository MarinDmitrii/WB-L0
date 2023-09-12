package builder

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/MarinDmitrii/WB-L0/internal/order/adapters"
	"github.com/MarinDmitrii/WB-L0/internal/order/usecase"
	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Application struct {
	SaveOrder     *usecase.SaveOrderUseCase
	RestoreCache  *usecase.RestoreCacheUseCase
	GetOrderByUID *usecase.GetOrderByUIDUseCase
	StanConnect   stan.Conn
}

type PostgresConfig struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

func NewPostgresConfig() *PostgresConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	return &PostgresConfig{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
	}
}

type NatsConfig struct {
	clusterId string
	clientID  string
}

func NewNatsConfig() *NatsConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clusterID := os.Getenv("NATS_CLUSTER_ID")
	clientID := os.Getenv("NATS_CLIENT_ID")

	return &NatsConfig{
		clusterId: clusterID,
		clientID:  clientID,
	}
}

func NewApplication(ctx context.Context) (*Application, func()) {
	PostgresConfig := NewPostgresConfig()
	postgresConnect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		PostgresConfig.host,
		PostgresConfig.port,
		PostgresConfig.user,
		PostgresConfig.password,
		PostgresConfig.dbname,
	)

	db, err := sqlx.ConnectContext(ctx, "postgres", postgresConnect)
	if err != nil {
		panic(err)
	}

	NatsConfig := NewNatsConfig()
	stanConnect, err := stan.Connect(NatsConfig.clusterId, NatsConfig.clientID)
	if err != nil {
		log.Fatalf("Can't connect to the NATS Streaming: %v\n", err)
	}

	orderRepository := adapters.NewPostgresOrderRepository(db)
	cacheRepository := adapters.NewCacheOrderRepository(200)

	restoreCache := usecase.NewRestoreCacheUseCase(orderRepository, cacheRepository)
	err = restoreCache.Execute(ctx)
	if err != nil {
		log.Fatalf("Can't restore cache: %v\n", err)
	}

	return &Application{
			SaveOrder:     usecase.NewSaveOrderUseCase(orderRepository, cacheRepository),
			GetOrderByUID: usecase.NewGetOrderByUIDUseCase(orderRepository, cacheRepository),
			RestoreCache:  usecase.NewRestoreCacheUseCase(orderRepository, cacheRepository),
			StanConnect:   stanConnect,
		}, func() {
			db.Close()
			stanConnect.Close()
		}
}
