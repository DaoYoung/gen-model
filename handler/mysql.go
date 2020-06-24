package handler

// Constants for return types of golang
const (
    golangByteArray  = "[]byte"
    gureguNullInt    = "null.Int"
    sqlNullInt       = "sql.NullInt64"
    golangInt        = "int"
    golangInt64      = "int64"
    gureguNullFloat  = "null.Float"
    sqlNullFloat     = "sql.NullFloat64"
    golangFloat32    = "float32"
    golangFloat64    = "float64"
    gureguNullString = "null.String"
    sqlNullString    = "sql.NullString"
    gureguNullTime   = "null.Time"
    golangTime       = "time.Time"
)
const (
    importSql     = "database/sql"
    importNull    = "gopkg.in/guregu/null.v3"
    importTime    = "time"
    importNothing = ""
)

func mysqlTypeToGoType(mysqlType string, nullable bool, gureguTypes bool) (goType, importPackage string) {
    switch mysqlType {
    case "tinyint", "int", "smallint", "mediumint":
        if nullable {
            if gureguTypes {
                return gureguNullInt, importNull
            }
            return sqlNullInt, importSql
        }
        return golangInt, importNothing
    case "bigint":
        if nullable {
            if gureguTypes {
                return gureguNullInt, importNull
            }
            return sqlNullInt, importSql
        }
        return golangInt64, importNothing
    case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "json":
        if nullable {
            if gureguTypes {
                return gureguNullString, importNull
            }
            return sqlNullString, importSql
        }
        return "string", importNothing
    case "date", "datetime", "time", "timestamp":
        if nullable && gureguTypes {
            return gureguNullTime, importNull
        }
        return golangTime, importTime
    case "decimal", "double":
        if nullable {
            if gureguTypes {
                return gureguNullFloat, importNull
            }
            return sqlNullFloat, importSql
        }
        return golangFloat64, importNothing
    case "float":
        if nullable {
            if gureguTypes {
                return gureguNullFloat, importNull
            }
            return sqlNullFloat, importSql
        }
        return golangFloat32, importNothing
    case "binary", "blob", "longblob", "mediumblob", "varbinary":
        return golangByteArray, importNothing
    }
    return "", importNothing
}
