package out

import "context"

type OTPProvider interface{
	Send(ctx context.Context, phone string, code string) error
}