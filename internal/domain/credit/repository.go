package credit

import "context"

type AccountRepository interface {
	QueryUserCreditAccount(ctx context.Context, userID string) (AccountEntity, bool, error)
}
