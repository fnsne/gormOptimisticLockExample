package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
	"os"
	"sync"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("cannot load .env file")
		panic(err)
	}
	gormDB, err := gorm.Open(postgres.Open(
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei",
			os.Getenv("POSTGRES_IP"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			"test_database",
			os.Getenv("POSTGRES_PORT")),
	),
		&gorm.Config{})
	if err != nil {
		slog.Error("cannot connect to postgres")
		panic(err)
	}
	slog.Info("connection pool is ready")

	ctx := context.Background()
	err = gormDB.WithContext(ctx).Migrator().AutoMigrate(&Book{})
	if err != nil {
		slog.Error("cannot migrate")
		panic(err)
	}
	slog.Info("prepare BookRepo")
	bookRepo := NewBookRepo(gormDB)

	//AddBook
	b := NewBook(uuid.NewString(), "The Art of Computer Programming", "Donald Knuth", time.Now())
	err = bookRepo.Add(ctx, b)
	if err != nil {
		slog.Error("cannot add book")
		panic(err)
	}

	//GetBook
	var wg sync.WaitGroup
	//UpdateBook
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			increaseCount(bookRepo, ctx, b)
			wg.Done()
		}()
	}
	wg.Wait()
}

func increaseCount(bookRepo *BookRepo, ctx context.Context, b *Book) {
	for {
		currentB, err := bookRepo.Get(ctx, b.ID)
		if err != nil {
			slog.Error("cannot get book")
			panic(err)
		}
		currentB.Count = currentB.Count + 1
		err = bookRepo.Update(ctx, currentB)
		if err != nil {
			if errors.Is(err, ErrNoUpdated) {
				continue
			}
			if err != nil {
				slog.Error("cannot update book")
				panic(err)
			}
		}
		break
	}
}

type Book struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	Title       string
	Author      string
	Count       int
	PublishedAt time.Time
	Version     optimisticlock.Version
}
type BookRepo struct {
	gormDB *gorm.DB
}

func (r *BookRepo) Add(ctx context.Context, b *Book) error {
	err := r.gormDB.WithContext(ctx).Create(b).Error
	return err
}

func (r *BookRepo) Get(ctx context.Context, id string) (Book, error) {
	var b Book
	err := r.gormDB.WithContext(ctx).First(&b, "id = ?", id).Error
	return b, err
}

var ErrNoUpdated = fmt.Errorf("no row updated caused there is another transaction update before this transaction")

func (r *BookRepo) Update(ctx context.Context, book Book) error {
	tx := r.gormDB.WithContext(ctx).Updates(&book)
	if tx.RowsAffected != 1 {
		return ErrNoUpdated
	}
	return tx.Error
}

func NewBookRepo(gormDB *gorm.DB) *BookRepo {
	return &BookRepo{gormDB: gormDB.Debug()}
}

func NewBook(ID string, title string, author string, publishedAt time.Time) *Book {
	return &Book{ID: ID, Title: title, Author: author, PublishedAt: publishedAt}
}
