module gomysql

go 1.18

replace origin => ./origin

replace model => ./model

replace dbsql => ./dbsql

require (
	github.com/go-sql-driver/mysql v1.7.0
	model v0.0.0-00010101000000-000000000000
)
