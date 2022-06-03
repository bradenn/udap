// Copyright (c) 2022 Braden Nicholson

package generic

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func TestNewStore(t *testing.T) {
	conn, mock, _ := sqlmock.New()
	pg := postgres.New(postgres.Config{Conn: conn})
	db, _ := gorm.Open(pg, &gorm.Config{})

	store := NewStore[Mock](db)

	t.Run("FindAll", func(t *testing.T) {
		u1 := uuid.New()
		u2 := uuid.New()
		rows := mock.NewRows([]string{"name", "value", "id"})
		rows.AddRow("test1a", "test1b", u1)
		rows.AddRow("test2a", "test2b", u2)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "mocks"`)).WillReturnRows(rows)
		mock.ExpectBegin()
		mock.ExpectCommit()

		all, err := store.FindAll()
		if err != nil {
			t.Error(err)
		}
		if len(*all) != 2 {
			t.Error("Wrong number of rows returned")
		}
		res := *all
		if res[0].Name != "test1a" {
			t.Error("Results out of order")
		}

	})

}
