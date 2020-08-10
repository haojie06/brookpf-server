package main

type StatusResponse struct {
	Code      int
	Installed bool
	Enable    bool
	Records   []string
}
