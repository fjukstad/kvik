package main

import(
    "fmt" 
    "github.com/mattn/go-sqlite3"
)

func main() {
    fmt.Println("hello")

    db, err := sql.Open("sqlite3", "./foo.db")
    if err != nil {
            log.Fatal(err)
    }

    fmt.Println(db)
}
