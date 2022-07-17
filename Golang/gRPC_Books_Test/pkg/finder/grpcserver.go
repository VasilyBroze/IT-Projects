package finder

import (
	"context"
	"database/sql"
	"fmt"
	gRPC_Books_Test "gRPC_Books_Test/api/proto"
	"strings"
)

//Описываем структуру сервера, из сгенерированного Proto файла
// GRPCServer ...
type GRPCServer struct {
	gRPC_Books_Test.UnimplementedBooksAndAuthorsServer
}

//Функция принимает строку, и выдаёт необходимый результат, который вычисляет функция search
func (s *GRPCServer) GetAuthor(ctx context.Context, req *gRPC_Books_Test.AuthorRequest) (*gRPC_Books_Test.AuthorResponse, error) {
	return &gRPC_Books_Test.AuthorResponse{A: Search(req.A)}, nil
}

//На вход подается строка вида: Автор Пушкин / Книга Сборник Стихов, на выход получаем соответственно книги этого автора / авторы этой книги
func Search(input string) string {
	//ПОДКЛЮЧАЕМСЯ К БД
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/books")

	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected DB")

	defer db.Close()

	//ПРОВЕРЯЕМ ВВОД НА НАЛИЧИЕ КЛЮЧЕВЫХ СЛОВ И ЗАПРАШИВАЕМ С БД ДАННЫЕ
	if strings.HasPrefix(input, "Автор") {
		//Убираем ненужный префикс перед запросом в БД
		input = strings.Replace(input, "Автор ", "", 1)
		//Запрос в БД - ищем книги по нужному автору
		result, err := db.Query("SELECT book.name FROM book, author_book, author WHERE book.id_book=author_book.id_book AND author_book.id_author=author.id_author AND author.name = ?", input)

		if err != nil {
			panic(err.Error())
		}

		var bookName string
		var msg string

		//Формируем сообщение, если Книг/Авторов больше одного складываем их в одну строку
		for result.Next() {
			err = result.Scan(&bookName)
			if err != nil {
				panic(err.Error())
			}
			if msg == "" {
				msg += bookName
			} else {
				msg += ", " + bookName
			}
		}
		fmt.Println(msg)

		result.Close()

		//Если в БД нет подходящих строк - выводим "Не найдено"
		if msg == "" {
			msg = "Не найдено"
		}
		return msg
	} else {
		//ПРОВЕРЯЕМ ВВОД НА НАЛИЧИЕ КЛЮЧЕВЫХ СЛОВ И ЗАПРАШИВАЕМ С БД ДАННЫЕ
		if strings.HasPrefix(input, "Книга") {
			//Убираем ненужный префикс перед запросом в БД
			input = strings.Replace(input, "Книга ", "", 1)
			//Запрос в БД - ищем автора по нужной книге
			result, err := db.Query("SELECT author.name FROM author, author_book, book WHERE author.id_author=author_book.id_author AND author_book.id_book=book.id_book AND book.name = ?", input)

			if err != nil {
				panic(err.Error())
			}

			var authorName string
			var msg string

			//Формируем сообщение, если Книг/Авторов больше одного складываем их в одну строку
			for result.Next() {
				err = result.Scan(&authorName)
				if err != nil {
					panic(err.Error())
				}
				if msg == "" {
					msg += authorName
				} else {
					msg += ", " + authorName
				}
			}
			fmt.Println(msg)

			result.Close()
			//Если в БД нет подходящих строк - выводим "Не найдено"
			if msg == "" {
				msg = "Не найдено"
			}

			return msg
		}
	}
	return "Неверный формат"
}
