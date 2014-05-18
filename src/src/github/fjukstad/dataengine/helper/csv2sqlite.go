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
    log.Print("Creating database from file: ", csvFilename)

    file, err := os.Open(csvFilename)
    if err != nil {
        log.Panic(err)
    }

    defer file.Close() 

    reader := csv.NewReader(file) 
    firstRow := true

    tableName := strings.Split(csvFilename,".")[0]

    log.Print("Tablename: ", tableName) 
    
    transaction, err := db.Begin()
    if err != nil{

        log.Panic("Transaction begin:", err)
    }
    
    var columnNames [] string
    var stmt *sql.Stmt

    for {

        record, err := reader.Read() 
        
        if err == io.EOF{
            break
        } else if err != nil {
            log.Panic(err) 
        }

        columnNames = record
        // First 
        if(firstRow){ 
            
            firstRow = false
            
            record, err := reader.Read() 
            
            if err != nil {
                log.Print(columnNames)
                log.Panic("Reader.read(): ", err)
            }
            
            createTableWithColumns(db, tableName, columnNames, record)
            stmt = prepareInsert(transaction, tableName, columnNames)
            defer stmt.Close()
            
            insertIntoTable(stmt, record) 
            continue

        } else {
            insertIntoTable(stmt, record) 
        }
    }
    transaction.Commit()
    
} 

func prepareInsert(transaction *sql.Tx , tableName string, 
        columnNames [] string) *sql.Stmt  {

    // Generate prepare statement
    sql := "insert into "+tableName+" ("
    for i, columnName := range columnNames { 
        sql += columnName 
        if(i < len(columnNames) - 1) {
            sql += ","
        }
    }

    sql += ") values(" + strings.Repeat("?,", len(columnNames))
    sql = strings.Trim(sql, ",") + ")"
    

    stmt, err := transaction.Prepare(sql)
    if err != nil { 
        log.Panic("Transaction prepare: " , err)
    }


    return stmt

}

func insertIntoTable(stmt *sql.Stmt, record []string) {

    var err error
    // Generate arguments from record
    args := make([]interface{}, len(record))
    for i, v := range(record) {
        switch determineType(v) {
        case "integer":
            args[i], err = strconv.Atoi(v)
            if(err != nil){
                log.Panic(err)
            }
        case "text":
            args[i] = v
        case "real":
            args[i], err = strconv.ParseFloat(v, 32)
            if(err != nil){
                log.Panic(err)
            }
        }
    }

    log.Print("Arguments:", args)

    // Execute query
    _, err = stmt.Exec(args...)
    if err != nil {

        log.Fatal("stmt.Exec: ", err)
    }

//    log.Print("Should insert ", record, "into table '", tableName, "'")

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
