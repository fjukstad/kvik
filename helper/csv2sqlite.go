package main

import(
    "fmt" 
    "log"
    "os"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
       
)

func main() {
    
    os.Remove("./foo.db")

    db, err := sql.Open("sqlite3", "./foo.db")
    if err != nil {
            log.Panic(err)
    }

    defer db.Close()

    sql := `
        create table foo (id integer not null primary key, name text, age integer);
        delete from foo;
        `
    _, err = db.Exec(sql)
    if err != nil {
            log.Printf("%q: %s\n", err, sql)
            return
    }

    tx, err := db.Begin()
    if err != nil{
        log.Panic(err)
    }

    stmt, err := tx.Prepare("insert into foo(id,name,age) values(?,?,?)")
    if err != nil { 
        log.Panic(err)
    }

    defer stmt.Close()

    for i := 0; i < 100; i ++ {
        _, err = stmt.Exec(i, fmt.Sprintf("item %d",i), i*2)
        if err != nil{
            log.Panic(err) 
        } 
    } 

    tx.Commit() 

    rows, err := db.Query("select id, name, age from foo") 
    if err != nil {
        log.Panic(err)
    } 

    fmt.Println("Extracting from db:")
    for rows.Next() {
        var id int
        var name string
        var age int
        rows.Scan(&id, &name, &age)
        fmt.Println("id =",id,"name =",name,"age =",age)
    }

    rows.Close()



}
