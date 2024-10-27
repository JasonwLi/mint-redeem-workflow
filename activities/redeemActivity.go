package activities

import (
	"context"
	"fmt"
	"mint-redeem-workflow/deps"
)

type RedeemActivityResponse struct {
	RequestId string
}

func RedeemActivity(ctx context.Context, amount float64, recipient string, requestId string) (RedeemActivityResponse, error) {
	deps, err := deps.NewDependencies()
	if err != nil {
		return RedeemActivityResponse{
			RequestId: requestId,
		}, err
	}

	resp, err := deps.BraleClient.Redeem(amount, recipient, requestId)
	if err != nil {
		return RedeemActivityResponse{
			RequestId: requestId,
		}, err
	}

	if len(resp.Errors) > 0 {
		return RedeemActivityResponse{
			RequestId: requestId,
		}, fmt.Errorf(resp.Errors[0].Detail)
	}

	return RedeemActivityResponse{
		RequestId: requestId,
	}, nil
}
