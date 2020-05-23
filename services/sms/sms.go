package sms

import (
	"context"
	"time"

	"github.com/ttlv/common_utils/services/config"
	"google.golang.org/grpc"
)

type SmsServerInterface interface {
	Send(brand string, country string, phone string, content string) (*SendResp, error)
}

type Service struct {
	Client     SmsClient
	Connection *grpc.ClientConn
}

func NewSmsService() (*Service, *grpc.ClientConn) {
	conn, err := grpc.Dial(config.MustGetSmsServiceConfig().ServerAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return &Service{Client: NewSmsClient(conn), Connection: conn}, conn
}

func (ser *Service) Send(brand string, country string, phone string, content string) (*SendResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return ser.Client.Send(ctx, &SendParams{Brand: brand, Country: country, Phone: phone, Content: content})
}

func (ser *Service) Close() {
	ser.Connection.Close()
}
