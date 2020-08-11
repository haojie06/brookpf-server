package main

type Server struct {
	ID       int    `id`
	Name     string `name`
	IP       string `ip`
	Port     string `port`
	UserName string `username`
	Password string `password`
	Desc     string `desc`
}
type Config struct {
	UserName string   `username`
	Password string   `password`
	Port     string   `port`
	Servers  []Server `servers`
}
type DataResponse struct {
	Code int
	Msg  string
	Data map[string]interface{}
}
type MessageResponse struct {
	Code int
	Msg  string
}
