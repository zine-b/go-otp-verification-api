package provider

import (
	"context"
	"fmt"
	"time"
)

type SMSProvider struct{}

func NewSMSProvider()*SMSProvider{
	return &SMSProvider{}
}

func (p *SMSProvider) Send(ctx context.Context, phone string, code string)error{
	select {
	case <-time.After(300 * time.Millisecond):
		fmt.Printf("Sending SMS OTP %s to %s\n", code, phone)
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}