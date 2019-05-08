package TaibaiJson

type Object interface{}
type Array []Object

type JsonObject map[string]Object
type JsonArray Array
