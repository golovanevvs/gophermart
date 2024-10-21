package service

import "context"

func (os *orderServiceStr) UploadOrder(ctx context.Context, userID int, orderNumber int) (int, error) {
	orderID, err := os.st.SaveOrderNumberByUserID(ctx, userID, orderNumber)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}
