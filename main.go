package main

import (
	"database/sql"
	"fmt"
	"strings"
	"log"
	"flag"
	_ "github.com/alexbrainman/odbc"
)

type Mssql struct {
	*sql.DB
	server string
	database   string
	windows    bool
	sa         SA
	port string
}

type SA struct {
	user   string
	passwd string
}
//数据库配置
func (m *Mssql) Open() (err error) {
	var conf []string
	conf = append(conf,"driver={sql server}")
	conf = append(conf,"server="+m.server+","+m.port)
	conf = append(conf, "Provider=SQLOLEDB")
	if m.windows {
		conf = append(conf, "integrated security=SSPI")
	}else {
		conf = append(conf, "Initial Catalog="+m.database)
		conf = append(conf, "user id="+m.sa.user)
		conf = append(conf, "password="+m.sa.passwd)
	}

	m.DB, err = sql.Open("odbc", strings.Join(conf, ";"))
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
//获取所有数据库名
func get_db_name(db Mssql)(db_name_list []string){
	// 列数据库
	query := "SELECT Name FROM Master..SysDatabases ORDER BY Name"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("query: ", err.Error()+"!!!")
		return
	}
	for rows.Next() {
		var dbname string
		if err := rows.Scan(&dbname); err != nil {
			log.Fatal(err)
		}
		db_name_list = append(db_name_list,dbname)
	}
	rows.Close()
	return
}
//获取所有表名
func get_tables_name(db Mssql,dbname string)(tables_list []string){
	query := "SELECT Name FROM "+dbname+"..SysObjects Where XType='U' ORDER BY Name "
	tables,err := db.Query(query)
	if err != nil {
		fmt.Println("query: ", err)
		return
	}
	for tables.Next() {
		var tables_name string
		if err := tables.Scan(&tables_name); err != nil {
			log.Fatal(err)
		}
		tables_list = append(tables_list,tables_name)
	}
	tables.Close()
	return
}
//获取所有字段名
func get_columns_list(db Mssql,dbname string,tableName string)(columns_list []string){
	db.Exec("USE " + dbname + ";")
	query := "SELECT Name from SysColumns WHERE id=Object_Id('" + tableName + "')"
	columns_query,err := db.Query(query)
	if err != nil {
		fmt.Println("query: ", err)
		return
	}

	for columns_query.Next() {
		var columns_name string
		if err := columns_query.Scan(&columns_name); err != nil {
			log.Fatal(err)
		}
		columns_list = append(columns_list,columns_name)
	}
	defer columns_query.Close()
	return
}

//获取数据条数
func getDataCount(db Mssql,dbname string,tableName string)(columns_list int) {
	db.Exec("USE " + dbname+";")
	query := "SELECT COUNT(1) AS COUNT FROM " + "\""+ tableName + "\""
	count , err := db.Query(query)
	if err != nil {
		fmt.Println("query: ", err)
		return
	}

	for count.Next() {
		var count_name int
		if err := count.Scan(&count_name ); err != nil {
			log.Fatal(err)
		}
		columns_list = count_name
	}
	defer count.Close()
	return
}
//取样
func getDataSamp(db Mssql,dbname string,tableName string,columns_list []string)(DataSamp_list map[string][]string){
	DataSamp_list = make(map[string][]string)
	for _,columnsName := range columns_list {
		db.Exec("USE " + dbname +";")
		query := "SELECT TOP 20 " +columnsName+ " FROM " + "\"" + tableName + "\""
		var DataSamp string
		Samp, err := db.Query(query)
		if err != nil {
			fmt.Println("query: ", err)
			return
		}
		columns_list = nil
		for Samp.Next() {
			if err := Samp.Scan(&DataSamp);Samp == nil {
				log.Print(err)
			}
			columns_list = append(columns_list,DataSamp )
		}

		DataSamp_list[columnsName] = columns_list
		columns_list = nil
		defer Samp.Close()
	}
	return DataSamp_list
}


func main() {

	IP := flag.String("IP","127.0.0.1","databases IP address , default is '127.0.0.1'")
	windows := flag.Bool("Windows_verification",true,"use Windows verification(true or false), default is true")
	username :=flag.String("username","sa","databases username , default is 'sa'")
	password:=flag.String("password","password","databases password , default is 'password'")
	port := flag.String("port","1433","databases port , default is 1433")

	//fmt.Println(*IP,*windows,*username,*password)

	flag.Parse()

	db := Mssql{
		server: *IP ,
		windows: *windows,
		port:	*port,
		sa: SA{
			user:   *username,
			passwd: *password,
		},
	}
	// 连接数据库
	err := db.Open()
	if err != nil {
		fmt.Println("sql open:", err)
		return
	}
	defer db.Close()

	fmt.Println("<!DOCTYPE html>")
	fmt.Println("<body>")

	DatabaseList := get_db_name(db)
	var DataCount int
	TablesList := make(map[string][]string)

	for _,GetDatabasesName := range DatabaseList {
		if GetDatabasesName != "master" && GetDatabasesName != "model" && GetDatabasesName != "msdb" && GetDatabasesName != "ReportServer" && GetDatabasesName != "ReportServerTempDB" && GetDatabasesName != "tempdb" {

			Tables_list := get_tables_name(db,GetDatabasesName)
			for _,GetTablesName := range Tables_list{
				DataCount = getDataCount(db,GetDatabasesName,GetTablesName)
				if GetTablesName != "sysdiagrams" {
					TablesList[GetDatabasesName] = Tables_list
					ColumnsList := get_columns_list(db,GetDatabasesName,GetTablesName)
					GetColumnsList := ColumnsList
					//Columns_List[GetTablesName] = ColumnsList
					GetSamp := getDataSamp(db,GetDatabasesName,GetTablesName,ColumnsList)
					fmt.Println("<table border=\"1\" cellspacing=\"0\">")
					fmt.Println("<tr><td>"+"databases name"+"</td>"+"<td>"+GetDatabasesName+"</td></tr>")
					fmt.Println("<tr><td>"+"DataCount "+"</td><td>",DataCount,"</td></tr>")
					fmt.Println("<tr><td>"+"table name"+"</td><td>"+GetTablesName+"</td>")
					fmt.Println("<tr><td>"+"columns name </td>")
					for _,text := range GetColumnsList{
						fmt.Println("<tr>")
						ColumnsName := text
						fmt.Print("<td>"+ColumnsName,"</td>\n")
						Samp_list := GetSamp[ColumnsName]
						//fmt.Print("GETSamp is :",Samp_list,"\n")
						for _,Samp_value := range Samp_list{
							fmt.Println("<td>"+Samp_value+"</td>")
						}
						fmt.Println("</tr>")
					}
					fmt.Println("</table>")
					fmt.Println("<br>")
					fmt.Println()
				}
			}
		}
	}
	fmt.Println("</body>")
	fmt.Println("</html>")
}
