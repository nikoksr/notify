package whatsapp

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	service, err := New()
	if err != nil {
		t.Errorf("New() error = %v, want nil", err)
		return
	}
	if service == nil {
		t.Error("New() returned nil service")
	}
}

func TestService_LoginWithSessionCredentials(t *testing.T) {
	service := &Service{}
	err := service.LoginWithSessionCredentials("", "", "", "", nil, nil)
	if err != nil {
		t.Errorf("LoginWithSessionCredentials() error = %v, want nil", err)
	}
}

func TestService_LoginWithQRCode(t *testing.T) {
	service := &Service{}
	err := service.LoginWithQRCode()
	if err != nil {
		t.Errorf("LoginWithQRCode() error = %v, want nil", err)
	}
}

func TestService_AddReceivers(_ *testing.T) {
	service := &Service{}
	// Should not panic
	service.AddReceivers("test1", "test2")
}

func TestService_Send(t *testing.T) {
	service := &Service{}
	err := service.Send(context.Background(), "subject", "message")
	if err != nil {
		t.Errorf("Send() error = %v, want nil", err)
	}
}
