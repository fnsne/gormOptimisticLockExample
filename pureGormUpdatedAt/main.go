package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	//UpdateBook's count concurrently 20 times
	var wg sync.WaitGroup
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
	err := bookRepo.UpdateBookCount(ctx, b.ID, 1)
	if err != nil {
		slog.Error("cannot update book")
		panic(err)
	}
}

type Book struct {
	ID          string `gorm:"type:uuid;primaryKey" json:"ID,omitempty"`
	Title       string
	Author      string
	Count       int
	PublishedAt time.Time
	UpdatedAt   time.Time `gorm:"type:timestamp(3)" json:"updatedAt"`
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

func (r *BookRepo) UpdateBookCount(ctx context.Context, id string, count int) error {
	for {
		var b Book
		err := r.gormDB.WithContext(ctx).First(&b, "id = ?", id).Error
		if err != nil {
			return err
		}
		b.Count += count
		tx := r.gormDB.WithContext(ctx).Model(&Book{}).Where("updated_at = ?", b.UpdatedAt).Updates(&b)
		if tx.Error != nil {
			return tx.Error
		}
		if tx.RowsAffected == 1 {
			break
		}
	}
	return nil
}

func NewBookRepo(gormDB *gorm.DB) *BookRepo {
	return &BookRepo{gormDB: gormDB.Debug()}
}

func NewBook(ID string, title string, author string, publishedAt time.Time) *Book {
	return &Book{ID: ID, Title: title, Author: author, PublishedAt: publishedAt}
}
