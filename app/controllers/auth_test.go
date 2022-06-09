package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/memnix/memnixrest/app/models"
	"reflect"
	"testing"
)

func TestAuthDebugMode(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name string
		args args
		want models.ResponseAuth
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AuthDebugMode(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthDebugMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckAuth(t *testing.T) {
	type args struct {
		c *fiber.Ctx
		p models.Permission
	}
	tests := []struct {
		name string
		args args
		want models.ResponseAuth
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckAuth(tt.args.c, tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsConnected(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 models.ResponseAuth
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := IsConnected(tt.args.c)
			if got != tt.want {
				t.Errorf("IsConnected() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("IsConnected() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Login(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Logout(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Logout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Register(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := User(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("User() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_extractToken(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractToken(tt.args.c); got != tt.want {
				t.Errorf("extractToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jwtKeyFunc(t *testing.T) {
	type args struct {
		token *jwt.Token
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jwtKeyFunc(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("jwtKeyFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("jwtKeyFunc() got = %v, want %v", got, tt.want)
			}
		})
	}
}
