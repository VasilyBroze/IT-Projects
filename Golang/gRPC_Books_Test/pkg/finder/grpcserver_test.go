package finder

import (
	"database/sql"
	_ "github.com/go-sql-driver/mySQL"
	"testing"
)

//Тест на подключение к БД
func TestDBconnect(t *testing.T) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/books")
	defer db.Close()
	if err != nil {
		t.Errorf("DB connection fail, %v", err)
	}
}

//Проверка функции по книге в БД
func TestSearch(t *testing.T) {
	//Задание значений
	inputString := "Книга Капитанская дочка"
	expected := "Пушкин"

	result := Search(inputString)

	if result != expected {
		t.Errorf("Incorrect result. Expected %v, got %v", expected, result)
	}
}
