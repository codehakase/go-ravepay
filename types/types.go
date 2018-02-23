package types

var (
	Envrionment string = "staging"
	Version     string = "1"
)

// RaveErr describes errors which are usually associated with the 400 http status code
type RaveErr struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}

// ValidationErr a pointer to RaveErr, are returned when one or more validation rules fail
// Examples include not passing required parameters e.g.
// not passing the transaction / provider ref during a re-query call will result in the error below:
//  {
//		"status":"error",
//		"message":"Cardno is required",
//		"data": null
//	}
type ValidationErr *RaveErr
