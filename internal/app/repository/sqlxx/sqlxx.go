package sqlxx

import (
	"fmt"
	"log"
	"time"

	"go-dex/internal/pkg/token"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config конфигурационные настройки БД
type Config struct {
	ConnectString   string        `yaml:"-" json:"-"`                                 // строка подключения к БД
	Host            string        `yaml:"host" json:"host"`                           // host БД
	Port            string        `yaml:"port" json:"port"`                           // порт БД
	Dbname          string        `yaml:"dbname" json:"dbname"`                       // имя БД
	SslMode         string        `yaml:"ssl_mode" json:"ssl_mode"`                   // режим SSL
	User            string        `yaml:"user" json:"user"`                           // пользователь БД
	Pass            string        `yaml:"pass" json:"pass"`                           // пароль пользователя БД
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" json:"conn_max_lifetime"` // время жизни подключения в миллисекундах
	MaxOpenConns    int           `yaml:"max_open_conns" json:"max_open_conns"`       // максимальное количество открытых подключений
	MaxIdleConns    int           `yaml:"max_idle_conns" json:"max_idle_conns"`       // максимальное количество простаивающих подключений
	DriverName      string        `yaml:"driver_name" json:"driver_name"`             // имя драйвера "postgres" | "pgx" | "godror"
}

func LoadConfig() *Config {
	var config Config

	config.Host = "localhost"
	config.Port = "5432"
	config.Dbname = "go_dex"
	config.SslMode = "disable"
	config.User = "postgres"
	config.Pass = "postgres"
	config.ConnMaxLifetime = 5 * time.Second
	config.MaxOpenConns = 10
	config.MaxIdleConns = 5
	config.DriverName = "pgx"

	return &config
}

// DBrepo is a wrapper around sqlx.DBrepo
type DBrepo struct {
	*sqlx.DB

	cfg *Config
}

// New - create new connect to DB
func New(cfg *Config) (db *DBrepo, myerr error) {

	// Сформировать строку подключения
	cfg.ConnectString = fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s user=%s password=%s ", cfg.Host, cfg.Port, cfg.Dbname, cfg.SslMode, cfg.User, cfg.Pass)

	// Создаем новый сервис
	db = &DBrepo{
		cfg: cfg,
	}
	sqlxDb, err := sqlx.Connect(cfg.DriverName, cfg.ConnectString)

	db.DB = sqlxDb

	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DBrepo) GetTokens(tokens *[]token.Token) error {
	err := db.Select(&tokens, "SELECT symbol, name, address FROM token")

	if err != nil {
		fmt.Printf("failed to get tokens: %v", err)
	}

	return nil
}

func (db *DBrepo) AddUser(address string, inviterId int) error {
	result, err := db.Exec(
		`INSERT INTO "user" (address, points) VALUES ($1, $2)`,
		address, 100)

	if err != nil {
		log.Fatalf("Error inserting user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Error getting affected rows: %v", err)
	}

	fmt.Printf("Number of rows affected: %d\n", rowsAffected)
	return nil
}
