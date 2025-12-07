package pkg_test

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fazriegi/fintrack-be/internal/pkg"
	"github.com/jmoiron/sqlx"
)

func TestIntersection(t *testing.T) {
	s1 := []int{1, 2, 3, 4}
	s2 := []int{3, 4, 5}

	result := pkg.Intersection(s1, s2)
	expected := []int{3, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

type userRow struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func TestScanRowsIntoStruct(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.
		NewRows([]string{"id", "name"}).
		AddRow(1, "alice").
		AddRow(2, "bob")

	mock.ExpectQuery("SELECT .* FROM users").
		WillReturnRows(rows)

	qRows, err := sqlxDB.Queryx("SELECT * FROM users")
	if err != nil {
		t.Fatal(err)
	}

	var result []userRow
	err = pkg.ScanRowsIntoStructs(qRows, &result)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(result))
	}

	if result[0].ID != 1 || result[0].Name != "alice" {
		t.Errorf("unexpected first row: %+v", result[0])
	}

	if result[1].ID != 2 || result[1].Name != "bob" {
		t.Errorf("unexpected second row: %+v", result[1])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestScanRowsIntoStructs_InvalidDest(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"col"}))

	rows, _ := sqlxDB.Queryx("SELECT 1")

	invalid := []userRow{}
	err := pkg.ScanRowsIntoStructs(rows, invalid)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "destSlice must be a pointer to a slice" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScanRowsIntoStructs_NonStructElement(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"col"}))

	rows, _ := sqlxDB.Queryx("SELECT 1")

	var invalid []int
	err := pkg.ScanRowsIntoStructs(rows, &invalid)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "slice elements must be structs" {
		t.Fatalf("unexpected error: %v", err)
	}
}
