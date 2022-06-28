package http

type MoveForm struct {
	FromUUID string `json:"from_uuid"`
	ToUUID   string `json:"to_uuid"`
}
