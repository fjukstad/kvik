package main

import(
    "log"
    "os"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "flag"
    "strings"
    "encoding/csv"
    "io"
    "strconv"
)

func main() {
    // cmd line flags
    var files = flag.String("files", "test.csv", "csv files to be imported")
    var dbPath = flag.String("db", "test.db" ,"path to db to import tables into")
    flag.Parse() 

    filelist := strings.Split(*files, ",")

    log.Print("Flags: files = ", filelist, ", db = ", *dbPath)

    db, err := sql.Open("sqlite3", *dbPath)

    if err != nil {
        log.Panic(err)
    }

    defer db.Close()

    for i := range(filelist){
        insertIntoDB(db, filelist[i])
    }

}

// Will create a table from csv file and insert into given db
func insertIntoDB(db *sql.DB, csvFilename string) {
    log.Print("Inserting ", csvFilename)

    file, err := os.Open(csvFilename)
    if err != nil {
        log.Panic(err)
    }

    defer file.Close() 

    reader := csv.NewReader(file) 
    firstRow := true

    tableName := strings.Split(csvFilename,".")[0]

    log.Print("TableName:", tableName) 
    
    for {

        record, err := reader.Read() 
        
        if err == io.EOF{
            break
        } else if err != nil {
            log.Panic(err) 
        }
        
        // First 
        if(firstRow){ 
            
            firstRow = false
            
            columnNames := record
            
            record, err := reader.Read() 
            
            if err != nil {
                log.Panic(err)
            }

            createTableWithColumns(db, tableName, columnNames, record)
        }

        insertIntoTable(db,tableName,record) 
    }
} 


func insertIntoTable(db *sql.DB, tableName string, record []string) {


    tx, err := db.Begin()
    if err != nil{
        log.Panic(err)
    }

    stmt, err := tx.Prepare("insert into "+tablename+"(id,name,age) values(?,?,?)")
    if err != nil { 
        log.Panic(err)
    }

    defer stmt.Close()
    log.Print("Should insert ", record, "into table '", tableName, "'")


}

func createTableWithColumns(db *sql.DB, tableName string, columnNames, firstRow [] string) {

    
    sql := `create table `+tableName+` (`
    

    types := make([]string, len(firstRow))
    for i, item := range firstRow { 
        item = strings.TrimSpace(item)
        t := determineType(item)
        types[i] = t
    }
    
    for i, columnName := range columnNames { 
        sql += columnName + " " + types[i] 
        if(i < len(columnNames) - 1) {
            sql += ","
        }
    }

    sql += ");"
    log.Print(sql)

    _, err := db.Exec(sql)
    if err != nil {
        log.Print(err,":", sql)
    }

}

// Determine type of element stored as string. Could be string/int/float
func determineType(item string) string {
    
    // Try integer
    _, err := strconv.Atoi(item)
    if err == nil {
        return "integer"
    }

    // Try float
    _, err = strconv.ParseFloat(item, 32)
    if err == nil {
        return "real"
    }

    // All else fails, its a string
    return "text" 

}
/*
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
*/
