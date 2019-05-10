package TaibaiJson

type JsonValue interface{}
type Array []JsonValue

type JsonObject map[string]JsonValue
type JsonArray Array
