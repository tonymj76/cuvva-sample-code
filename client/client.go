package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/crpc"
	"github.com/cuvva/cuvva-public-go/lib/jsonclient"
	"github.com/sirupsen/logrus"
	Handler "github.com/tonymj76/cuvva-sample-code/handlers"
	Models "github.com/tonymj76/cuvva-sample-code/models"
)

type MerchantClient struct {
	*crpc.Client
}

func (mc *MerchantClient) CreateMerchant(ctx context.Context, req *Models.CreateRequest) (res *Models.CreateResponse, err error) {
	return res, mc.Do(ctx, "create_merchant", "2021-06-29", req, &res)
}

func main() {
	client := &http.Client{
		Transport: jsonclient.NewAuthenticatedRoundTripper(nil, "Bearer", "...someSecret"),
		Timeout:   5 * time.Second,
	}
	var hs Handler.Metcher = &MerchantClient{
		Client: crpc.NewClient("http://127.0.0.1:3003/v1", client),
	}

	ctx := context.Background()
	res, err := hs.CreateMerchant(ctx, &Models.CreateRequest{
		NumberOfProduct: 233,
		Email:           "cuvves@gmail.com",
		BusinessName:    "Cuvva",
	})
	if err != nil {
		logrus.WithError(err).Fatalln("failed to send request")
	}
	fmt.Println("created merchant: ", res)
}
