package service

type PutRequest struct {
	Key   string `form:"key"`
	Value string `form:"value"`
}

type CreateRequest struct {
	DBName string `form:"dbname"`
}

type DelRequest struct {
	Key string `form:"key"`
}

type GetRequest struct {
	Key string `form:"key"`
}
