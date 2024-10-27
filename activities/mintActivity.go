package activities

import (
	"context"
	"fmt"
	"mint-redeem-workflow/deps"
)

type MintActivityResponse struct {
	RequestId string
}

func MintActivity(ctx context.Context, amount float64, recipient string, requestId string) (MintActivityResponse, error) {
	deps, err := deps.NewDependencies()
	if err != nil {
		return MintActivityResponse{
			RequestId: requestId,
		}, err
	}

	resp, err := deps.BraleClient.Mint(amount, recipient, requestId)
	if err != nil {
		return MintActivityResponse{
			RequestId: requestId,
		}, err
	}

	if len(resp.Errors) > 0 {
		return MintActivityResponse{
			RequestId: requestId,
		}, fmt.Errorf(resp.Errors[0].Detail)
	}

	return MintActivityResponse{
		RequestId: requestId,
	}, nil
}
