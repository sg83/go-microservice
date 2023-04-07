package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/sg83/go-microservice/article-api/models"
	"go.uber.org/zap"
)

type ArticlesDb struct {
	postgres *sql.DB
	l        *zap.Logger
}

// a struct to hold all the db connection information
type connection struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewDB(l *zap.Logger) *ArticlesDb {
	artdb := &ArticlesDb{nil, l}
	err := artdb.init()
	if err != nil {
		l.Fatal("Could not initialize database", zap.String(" error: ", err.Error()))
		return nil
	}
	return artdb
}

// InitDB initializes the database connection
func (db *ArticlesDb) init() error {
	err := godotenv.Load("config/.env")

	if err != nil {
		fmt.Printf("Error loading .env file: %s\n", err.Error())
		return err
	}

	connInfo := connection{
		Host:     os.Getenv("POSTGRES_URL"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
	}

	// Replace the connection string with your PostgreSQL connection details
	db.postgres, err = sql.Open("postgres", connToString(connInfo))
	if err != nil {
		db.l.Fatal(err.Error())
		return err
	}

	// Ping the database to ensure a connection is established
	err = db.postgres.Ping()
	if err != nil {
		db.l.Fatal(err.Error())
		return err
	}

	db.l.Info("Connected to the database")
	return nil
}

// Take our connection struct and convert to a string for our db connection info
func connToString(info connection) string {

	//result := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	result := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		info.User, info.Password, info.Host, info.Port, info.DBName)
	fmt.Println(result)

	return result

}

func (db *ArticlesDb) GetArticleByID(id int) (*models.Article, error) {
	var a models.Article
	db.l.Info("Get article ", zap.Int("id :", id))

	err := db.postgres.QueryRow("SELECT * FROM articles WHERE id = $1", id).Scan(&a.ID, &a.Title, &a.Body, &a.Date, pq.Array(&a.Tags))

	if err != nil {
		db.l.Error(err.Error())
		return nil, err
	}

	db.l.Info("Get article success")
	return &a, nil
}

// AddProduct adds a new product to the database
func (db *ArticlesDb) AddArticle(ar models.Article) error {
	db.l.Info("Add new article ", zap.String("title :", ar.Title))
	// get the next id in sequence
	query := `insert into articles(title, content) values($1, $2, $3, $4, $5);`

	_, err := db.postgres.Exec(query, ar.Title, ar.Body, ar.Date, ar.Tags)

	if err != nil {
		return err
	}
	return nil
}

func (db *ArticlesDb) Close() {
	db.postgres.Close()
}
