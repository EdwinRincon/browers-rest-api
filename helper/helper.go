package helper

import (
	"log"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

type ResponseJSON struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseJSONSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(200, ResponseJSON{message, data})
}

func HandleError(c *gin.Context, code int, message string, err error) {
	gormMessage := GormErrorFriendlyHandling(err)
	log.Println("Error: ", err)

	if gormMessage != "" {
		c.JSON(code, gin.H{
			"error":   gormMessage,
			"message": message,
		})
	}
	c.JSON(code, gin.H{
		"message": message,
	})
}

func IsStrongPassword(password string) bool {
	var hasUpperCase = false
	var hasLowerCase = false
	var hasDigit = false
	var hasSpecialChar = false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpperCase = true
		}
		if unicode.IsLower(char) {
			hasLowerCase = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecialChar = true
		}
		if len(password) < 8 {
			return false
		}
	}
	if !hasUpperCase || !hasLowerCase || !hasDigit || !hasSpecialChar {
		return false
	}
	return true
}

/*
// ErrRecordNotFound record not found error
ErrRecordNotFound = logger.ErrRecordNotFound
// ErrInvalidTransaction invalid transaction when you are trying to `Commit` or `Rollback`
ErrInvalidTransaction = errors.New("invalid transaction")
// ErrNotImplemented not implemented
ErrNotImplemented = errors.New("not implemented")
// ErrMissingWhereClause missing where clause
ErrMissingWhereClause = errors.New("WHERE conditions required")
// ErrUnsupportedRelation unsupported relations
ErrUnsupportedRelation = errors.New("unsupported relations")
// ErrPrimaryKeyRequired primary keys required
ErrPrimaryKeyRequired = errors.New("primary key required")
// ErrModelValueRequired model value required
ErrModelValueRequired = errors.New("model value required")
// ErrModelAccessibleFieldsRequired model accessible fields required
ErrModelAccessibleFieldsRequired = errors.New("model accessible fields required")
// ErrSubQueryRequired sub query required
ErrSubQueryRequired = errors.New("sub query required")
// ErrInvalidData unsupported data
ErrInvalidData = errors.New("unsupported data")
// ErrUnsupportedDriver unsupported driver
ErrUnsupportedDriver = errors.New("unsupported driver")
// ErrRegistered registered
ErrRegistered = errors.New("registered")
// ErrInvalidField invalid field
ErrInvalidField = errors.New("invalid field")
// ErrEmptySlice empty slice found
ErrEmptySlice = errors.New("empty slice found")
// ErrDryRunModeUnsupported dry run mode unsupported
ErrDryRunModeUnsupported = errors.New("dry run mode unsupported")
// ErrInvalidDB invalid db
ErrInvalidDB = errors.New("invalid db")
// ErrInvalidValue invalid value
ErrInvalidValue = errors.New("invalid value, should be pointer to struct or slice")
// ErrInvalidValueOfLength invalid values do not match length
ErrInvalidValueOfLength = errors.New("invalid association values, length doesn't match")
// ErrPreloadNotAllowed preload is not allowed when count is used
ErrPreloadNotAllowed = errors.New("preload is not allowed when count is used")
// ErrDuplicatedKey occurs when there is a unique key constraint violation
ErrDuplicatedKey = errors.New("duplicated key not allowed")
// ErrForeignKeyViolated occurs when there is a foreign key constraint violation
ErrForeignKeyViolated = errors.New("violates foreign key constraint")
*/
var customGormErrorMessages = map[string]string{
	"record not found":                 "registro no encontrado",
	"invalid transaction":              "transacción no válida al intentar 'Commit' o 'Rollback'",
	"not implemented":                  "operación no implementada",
	"WHERE conditions required":        "se requieren condiciones 'WHERE'",
	"unsupported relations":            "relaciones no admitidas",
	"primary key required":             "clave primaria requerida",
	"model value required":             "valor del modelo requerido",
	"model accessible fields required": "campos accesibles del modelo requeridos",
	"sub query required":               "se requiere una subconsulta",
	"unsupported data":                 "datos no admitidos",
	"unsupported driver":               "controlador no admitido",
	"registered":                       "ya registrado",
	"invalid field":                    "campo no válido",
	"empty slice found":                "se encontró una lista vacía",
	"dry run mode unsupported":         "modo de prueba en seco no admitido",
	"invalid db":                       "base de datos no válida",
	"invalid value, should be pointer to struct or slice": "valor no válido, debe ser un puntero a una estructura o una lista",
	"invalid association values, length doesn't match":    "valores no válidos, la longitud no coincide",
	"preload is not allowed when count is used":           "la precarga no está permitida cuando se usa la cuenta",
	"duplicated key not allowed":                          "violación de restricción de clave única",
	"violates foreign key constraint":                     "viola la restricción de clave foránea",
}

func GormErrorFriendlyHandling(err error) string {
	if err != nil {
		errorGorm := strings.ToLower(err.Error())
		for substr, customMessage := range customGormErrorMessages {
			if strings.Contains(errorGorm, substr) {
				return customMessage
			}
		}
	}
	return ""
}
