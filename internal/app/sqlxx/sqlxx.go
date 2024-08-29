package sqlxx

import (
	"fmt"
	"time"

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
	config.Port = "5431"
	config.Dbname = "postgres"
	config.SslMode = "disable"
	config.User = "postgres"
	config.Pass = "postgres"
	config.ConnMaxLifetime = 5 * time.Second
	config.MaxOpenConns = 10
	config.MaxIdleConns = 5
	config.DriverName = "pgx"

	return &config
}

// DB is a wrapper around sqlx.DB
type DB struct {
	*sqlx.DB

	cfg *Config
}

// New - create new connect to DB
func New(cfg *Config) (db *DB, myerr error) {

	// Сформировать строку подключения
	cfg.ConnectString = fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s user=%s password=%s ", cfg.Host, cfg.Port, cfg.Dbname, cfg.SslMode, cfg.User, cfg.Pass)

	// Создаем новый сервис
	db = &DB{
		cfg: cfg,
	}
	sqlxDb, err := sqlx.Connect(cfg.DriverName, cfg.ConnectString)

	db.DB = sqlxDb

	if err != nil {
		return nil, err
	}
	return db, nil
}
