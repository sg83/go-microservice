package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/lib/pq"

	"github.com/sg83/go-microservice/article-api/models"
	"go.uber.org/zap"
)

type ArticlesDb struct {
	postgres *sql.DB
	l        *zap.Logger
}

type ArticlesData interface {
	GetArticleByID(id int) (*models.Article, error)
	AddArticle(ar models.Article) error
	GetArticlesForTagAndDate(tag string, d time.Time) ([]int, error)
	GetRelatedTagsForTag(tag string, articles []int) ([]string, error)
	init() error
	Close()
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
		l.Fatal("Could not initialize database", zap.Error(err))
		return nil
	}
	return artdb
}

// InitDB initializes the database connection
func (db *ArticlesDb) init() error {
	err := godotenv.Load("config/.env")

	if err != nil {
		db.l.Error("Error loading .env file", zap.Error(err))
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
		db.l.Error("Error connecting to database", zap.Error(err))
		return err
	}

	// Ping the database to ensure a connection is established
	err = db.postgres.Ping()
	if err != nil {
		db.l.Error("Could not Ping database", zap.Error(err))
		return err
	}

	db.l.Info("Connected to the database")
	return nil
}

// Take our connection struct and convert to a string for our db connection info
func connToString(info connection) string {
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		info.User, info.Password, info.Host, info.Port, info.DBName)
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

	query := `insert into articles(id, title, date, body, tags) values(nextval('articles_id_seq'), $1, $2, $3, $4);`

	_, err := db.postgres.Exec(query, ar.Title, ar.Date, ar.Body, pq.Array(ar.Tags))

	if err != nil {
		db.l.Error("DB Query failed ", zap.Error(err))
		return err
	}
	return nil
}

func (db *ArticlesDb) Close() {
	db.postgres.Close()
}

func (db *ArticlesDb) GetArticlesForTagAndDate(tag string, d time.Time) ([]int, error) {
	db.l.Info("GetArticlesForTagAndDate", zap.String("date: ", d.String()))
	rows, err := db.postgres.Query("SELECT id FROM articles WHERE $1 = ANY(tags) AND date = $2", tag, d.Format("2006-01-02"))
	if err != nil {
		db.l.Error("sql query failed", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			db.l.Error("error scanning row", zap.Error(err))
			return nil, err
		}
		db.l.Info("Article id", zap.Int("id ", id))

		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		db.l.Error("error scanning:", zap.Error(err))
		return nil, err
	}

	return ids, nil
}

func (db *ArticlesDb) GetRelatedTagsForTag(tag string, articles []int) ([]string, error) {
	var tags []string
	for _, id := range articles {
		rows, err := db.postgres.Query("SELECT tags FROM articles WHERE id = $1 AND $2 = ANY(tags)", id, tag)
		if err != nil {
			db.l.Error("sql query failed", zap.Error(err))
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var tagArr []string
			err := rows.Scan(pq.Array(&tagArr))
			if err != nil {
				db.l.Error("row scan failed", zap.Error(err))
				return nil, err
			}
			for _, t := range tagArr {
				if t != tag && !contains(tags, t) {
					tags = append(tags, t)
				}
			}
		}
		if err := rows.Err(); err != nil {
			db.l.Error("Errors scanning rows", zap.Error(err))
			return nil, err
		}
	}

	return tags, nil
}

// Function to check if given string is present in the slice of strings
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
