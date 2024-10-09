package service

import "fmt"

func (sv serviceStr) RegisterUser() string {
	return fmt.Sprintf("registerUser")
}
