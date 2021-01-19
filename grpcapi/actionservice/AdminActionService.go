package actionservice

import (
	"fmt"

	"github.com/delta/dalal-street-server/models"
	actions_pb "github.com/delta/dalal-street-server/proto_build/actions"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func (d *dalalActionService) SendNews(ctx context.Context, req *actions_pb.SendNewsRequest) (*actions_pb.SendNewsResponse, error) {
	resp := &actions_pb.SendNewsResponse{}
	// now call functions from models
	resp.StatusMessage = "Done"
	resp.StatusCode = actions_pb.SendNewsResponse_OK
	return resp, nil
}

func (d *dalalActionService) OpenMarket(ctx context.Context, req *actions_pb.OpenMarketRequest) (*actions_pb.OpenMarketResponse, error) {

	var l = logger.WithFields(logrus.Fields{
		"method":        "OpenMarket",
		"param_session": fmt.Sprintf("%+v", ctx.Value("session")),
		"param_req":     fmt.Sprintf("%+v", req),
	})

	resp := &actions_pb.OpenMarketResponse{}
	makeError := func(st actions_pb.OpenMarketResponse_StatusCode, msg string) (*actions_pb.OpenMarketResponse, error) {
		resp.StatusCode = st
		resp.StatusMessage = msg
		return resp, nil
	}
	userId := getUserId(ctx)
	if !models.IsAdminAuth(userId) {
		return makeError(actions_pb.OpenMarketResponse_NotAdminUserError, "User is not admin")
	}

	err := models.OpenMarket(req.UpdateDayHighAndLow)

	if err != nil {
		l.Errorf("Request failed due to %+v: ", err)
		return makeError(actions_pb.OpenMarketResponse_InternalServerError, getInternalErrorMessage(err))
	}

	resp.StatusMessage = "Done"
	resp.StatusCode = actions_pb.OpenMarketResponse_OK
	return resp, nil
}

func (d *dalalActionService) CloseMarket(ctx context.Context, req *actions_pb.CloseMarketRequest) (*actions_pb.CloseMarketResponse, error) {

	var l = logger.WithFields(logrus.Fields{
		"method":        "CloseMarket",
		"param_session": fmt.Sprintf("%+v", ctx.Value("session")),
		"param_req":     fmt.Sprintf("%+v", req),
	})

	resp := &actions_pb.CloseMarketResponse{}
	makeError := func(st actions_pb.CloseMarketResponse_StatusCode, msg string) (*actions_pb.CloseMarketResponse, error) {
		resp.StatusCode = st
		resp.StatusMessage = msg
		return resp, nil
	}
	userId := getUserId(ctx)
	if !models.IsAdminAuth(userId) {
		return makeError(actions_pb.CloseMarketResponse_NotAdminUserError, "User is not admin")
	}

	err := models.CloseMarket(req.UpdatePrevDayClose)

	if err != nil {
		l.Errorf("Error closing the market due to %+v: ", err)
		return makeError(actions_pb.CloseMarketResponse_InternalServerError, getInternalErrorMessage(err))
	}

	resp.StatusCode = actions_pb.CloseMarketResponse_OK
	resp.StatusMessage = "OK"

	return resp, nil

}

func (d *dalalActionService) UpdateEndOfDayValues(ctx context.Context, req *actions_pb.UpdateEndOfDayValuesRequest) (*actions_pb.UpdateEndOfDayValuesResponse, error) {
	var l = logger.WithFields(logrus.Fields{
		"method":        "UpdateEndOfDayValues",
		"param_session": fmt.Sprintf("%+v", ctx.Value("session")),
		"param_req":     fmt.Sprintf("%+v", req),
	})

	resp := &actions_pb.UpdateEndOfDayValuesResponse{}
	makeError := func(st actions_pb.UpdateEndOfDayValuesResponse_StatusCode, msg string) (*actions_pb.UpdateEndOfDayValuesResponse, error) {
		resp.StatusCode = st
		resp.StatusMessage = msg
		return resp, nil
	}
	userId := getUserId(ctx)
	if !models.IsAdminAuth(userId) {
		return makeError(actions_pb.UpdateEndOfDayValuesResponse_NotAdminUser, "User is not admin")
	}

	err := models.UpdateEndOfDayValues()

	if err != nil {
		l.Errorf("Error updationg EndOfDayValues due to: %+v", err)
		return makeError(actions_pb.UpdateEndOfDayValuesResponse_InternalServerError, getInternalErrorMessage(err))
	}

	resp.StatusCode = actions_pb.UpdateEndOfDayValuesResponse_OK
	resp.StatusMessage = "OK"

	return resp, err
}

func (d *dalalActionService) SendNotifications(ctx context.Context, req *actions_pb.SendNotificationsRequest) (*actions_pb.SendNotificationsResponse, error) {

	var l = logger.WithFields(logrus.Fields{
		"method":        "SendNotifications",
		"param_session": fmt.Sprintf("%+v", ctx.Value("session")),
		"param_req":     fmt.Sprintf("%+v", req),
	})

	resp := &actions_pb.SendNotificationsResponse{}

	makeError := func(st actions_pb.SendNotificationsResponse_StatusCode, msg string) (*actions_pb.SendNotificationsResponse, error) {
		resp.StatusCode = st
		resp.StatusMessage = msg
		return resp, nil
	}
	userId := getUserId(ctx)
	if !models.IsAdminAuth(userId) {
		return makeError(actions_pb.SendNotificationsResponse_NotAdminUserError, "User is not admin")
	}

	if req.IsGlobal && req.UserId != 0 {
		l.Errorf("Cannot send Global Notification to Non Zero Id")
		return makeError(actions_pb.SendNotificationsResponse_InternalServerError, "Cannot send Global Notification to Non Zero Id")
	}

	err := models.SendNotification(req.UserId, req.Text, req.IsGlobal)

	if err != nil {
		l.Errorf("Request failed due to %+v: ", err)
		return makeError(actions_pb.SendNotificationsResponse_InternalServerError, getInternalErrorMessage(err))
	}

	resp.StatusMessage = "Done"
	resp.StatusCode = actions_pb.SendNotificationsResponse_OK
	return resp, nil
}

func (d *dalalActionService) LoadStocks(ctx context.Context, req *actions_pb.LoadStocksRequest) (*actions_pb.LoadStocksResponse, error) {

	var l = logger.WithFields(logrus.Fields{
		"method":        "LoadStocks",
		"param_session": fmt.Sprintf("%+v", ctx.Value("session")),
		"param_req":     fmt.Sprintf("%+v", req),
	})
	resp := &actions_pb.LoadStocksResponse{}

	userId := getUserId(ctx)
	makeError := func(st actions_pb.LoadStocksResponse_StatusCode, msg string) (*actions_pb.LoadStocksResponse, error) {
		resp.StatusCode = st
		resp.StatusMessage = msg
		return resp, nil
	}
	if !models.IsAdminAuth(userId) {
		return makeError(actions_pb.LoadStocksResponse_NotAdminUserError, "User is not admin")
	}
	err := models.LoadStocks()

	if err != nil {
		l.Errorf("Request failed due to %+v: ", err)
		return makeError(actions_pb.LoadStocksResponse_InternalServerError, getInternalErrorMessage(err))
	}

	resp.StatusMessage = "Done"
	resp.StatusCode = actions_pb.LoadStocksResponse_OK
	return resp, nil
}

func (d *dalalActionService) AddStocksToExchange(ctx context.Context, req *actions_pb.AddStocksToExchangeRequest) (*actions_pb.AddStocksToExchangeResponse, error) {

	var l = logger.WithFields(logrus.Fields{
		"method":        "AddStocksToExchange",
		"param_session": fmt.Sprintf("%+v", ctx.Value("session")),
		"param_req":     fmt.Sprintf("%+v", req),
	})
	resp := &actions_pb.AddStocksToExchangeResponse{}

	makeError := func(st actions_pb.AddStocksToExchangeResponse_StatusCode, msg string) (*actions_pb.AddStocksToExchangeResponse, error) {
		resp.StatusCode = st
		resp.StatusMessage = msg
		return resp, nil
	}
	userId := getUserId(ctx)
	if !models.IsAdminAuth(userId) {
		return makeError(actions_pb.AddStocksToExchangeResponse_NotAdminUserError, "User is not admin")
	}

	stock, err := models.GetStockCopy(req.StockId)

	l.Debugf("Adding new stocks to Exchange for %s", stock.FullName)

	if err != nil {
		l.Errorf("Request failed due to %+v: ", err)
		return makeError(actions_pb.AddStocksToExchangeResponse_InternalServerError, getInternalErrorMessage(err))

	}

	err = models.AddStocksToExchange(req.StockId, req.NewStocks)

	if err != nil {
		l.Errorf("Request failed due to %+v: ", err)
		return makeError(actions_pb.AddStocksToExchangeResponse_InternalServerError, getInternalErrorMessage(err))
	}

	resp.StatusMessage = "Done"
	resp.StatusCode = actions_pb.AddStocksToExchangeResponse_OK
	return resp, nil
}

func (d *dalalActionService) UpdateStockPrice(ctx context.Context, req *actions_pb.UpdateStockPriceRequest) (*actions_pb.UpdateStockPriceResponse, error) {

	var l = logger.WithFields(logrus.Fields{
		"method":        "UpdateStockPrice",
		"param_session": fmt.Sprintf("%+v", ctx.Value("session")),
		"param_req":     fmt.Sprintf("%+v", req),
	})

	userId := getUserId(ctx)
	resp := &actions_pb.UpdateStockPriceResponse{}

	makeError := func(st actions_pb.UpdateStockPriceResponse_StatusCode, msg string) (*actions_pb.UpdateStockPriceResponse, error) {
		resp.StatusCode = st
		resp.StatusMessage = msg
		return resp, nil
	}
	if !models.IsAdminAuth(userId) {
		return makeError(actions_pb.UpdateStockPriceResponse_NotAdminUserError, "User is not admin")
	}

	stock, err := models.GetStockCopy(req.StockId)

	l.Debugf("Adding new stocks to Exchange for %s", stock.FullName)

	if err != nil {
		l.Errorf("Request failed due to %+v: ", err)
		return makeError(actions_pb.UpdateStockPriceResponse_InternalServerError, getInternalErrorMessage(err))

	}

	err = models.UpdateStockPrice(req.StockId, req.NewPrice, 10000)

	if err != nil {
		l.Errorf("Request failed due to %+v: ", err)
		return makeError(actions_pb.UpdateStockPriceResponse_InternalServerError, getInternalErrorMessage(err))
	}

	resp.StatusMessage = "Done"
	resp.StatusCode = actions_pb.UpdateStockPriceResponse_OK
	return resp, nil
}

func (d *dalalActionService) AddMarketEvent(ctx context.Context, req *actions_pb.AddMarketEventRequest) (*actions_pb.AddMarketEventResponse, error) {

	var l = logger.WithFields(logrus.Fields{
		"method":        "AddMarketEvent",
		"param_session": fmt.Sprintf("%+v", ctx.Value("session")),
		"param_req":     fmt.Sprintf("%+v", req),
	})

	resp := &actions_pb.AddMarketEventResponse{}

	makeError := func(st actions_pb.AddMarketEventResponse_StatusCode, msg string) (*actions_pb.AddMarketEventResponse, error) {
		resp.StatusCode = st
		resp.StatusMessage = msg
		return resp, nil
	}
	userId := getUserId(ctx)
	if !models.IsAdminAuth(userId) {
		return makeError(actions_pb.AddMarketEventResponse_NotAdminUserError, "User is not admin")
	}

	if req.StockId != 0 {
		stock, err := models.GetStockCopy(req.StockId)
		l.Debugf("Adding Market Event for %s", stock.FullName)
		if err != nil {
			l.Errorf("Request failed due to %+v: ", err)
			return makeError(actions_pb.AddMarketEventResponse_InternalServerError, getInternalErrorMessage(err))
		}
	}

	if req.IsGlobal && req.StockId != 0 {
		l.Errorf("Cannot send Global Notification to Non Zero Stock Id")
		return makeError(actions_pb.AddMarketEventResponse_InternalServerError, "Cannot send Global Notification to Non Zero Stock Id")
	}

	err := models.AddMarketEvent(req.StockId, req.Headline, req.Text, req.IsGlobal, req.ImageUrl)

	if err != nil {
		l.Errorf("Request failed due to %+v: ", err)
		return makeError(actions_pb.AddMarketEventResponse_InternalServerError, getInternalErrorMessage(err))
	}

	resp.StatusMessage = "Done"
	resp.StatusCode = actions_pb.AddMarketEventResponse_OK
	return resp, nil
}

func (d *dalalActionService) SendDividends(ctx context.Context, req *actions_pb.SendDividendsRequest) (*actions_pb.SendDividendsResponse, error) {
	var l = logger.WithFields(logrus.Fields{
		"method":        "SendDividends",
		"param_session": fmt.Sprintf("%+v", ctx.Value("session")),
		"param_req":     fmt.Sprintf("%+v", req),
	})

	l.Infof("Request for dividends sent")
	resp := &actions_pb.SendDividendsResponse{}
	makeError := func(st actions_pb.SendDividendsResponse_StatusCode, msg string) (*actions_pb.SendDividendsResponse, error) {
		resp.StatusCode = st
		resp.StatusMessage = msg
		return resp, nil
	}
	userId := getUserId(ctx)
	if !models.IsAdminAuth(userId) {
		return makeError(actions_pb.SendDividendsResponse_NotAdminUserError, "User is not admin")
	}

	if !models.IsMarketOpen() {
		return makeError(actions_pb.SendDividendsResponse_MarketClosedError, "Market Is closed. You cannot send dividends right now.")
	}

	stockID := req.StockId
	dividendAmount := req.DividendAmount

	err := models.PerformDividendsTransaction(stockID, dividendAmount)

	if err == nil {
		resp.StatusCode = 0
		resp.StatusMessage = "OK"

	}

	switch e := err.(type) {
	case models.InvalidStockIdError:
		return makeError(actions_pb.SendDividendsResponse_InvalidStockIdError, e.Error())
	}
	if err != nil {
		l.Errorf("Request failed due to %+v: ", err)
		return makeError(actions_pb.SendDividendsResponse_InternalServerError, getInternalErrorMessage(err))
	}

	l.Infof("Request completed successfully")

	return resp, nil
}

func (d *dalalActionService) SetGivesDividends(ctx context.Context, req *actions_pb.SetGivesDividendsRequest) (*actions_pb.SetGivesDividendsResponse, error) {
	var l = logger.WithFields(logrus.Fields{
		"method":        "SetGivesDividends",
		"param_session": fmt.Sprintf("%+v", ctx.Value("session")),
		"param_req":     fmt.Sprintf("%+v", req),
	})

	l.Infof("Request for setting givesDividends")

	resp := &actions_pb.SetGivesDividendsResponse{}
	makeError := func(st actions_pb.SetGivesDividendsResponse_StatusCode, msg string) (*actions_pb.SetGivesDividendsResponse, error) {
		resp.StatusCode = st
		resp.StatusMessage = msg
		return resp, nil
	}

	userId := getUserId(ctx)
	if !models.IsAdminAuth(userId) {
		return makeError(actions_pb.SetGivesDividendsResponse_NotAdminUserError, "User is not admin")
	}

	if !models.IsMarketOpen() {
		return makeError(actions_pb.SetGivesDividendsResponse_MarketClosedError, "Market Is closed. You cannot set GivesDividends for stocks now.")
	}

	stockID := req.GetStockId()
	givesDividends := req.GetGivesDividends()

	err := models.SetGivesDividends(stockID, givesDividends)

	if err == models.InvalidStockError {
		return makeError(actions_pb.SetGivesDividendsResponse_InvalidStockIdError, "Invalid stock id provided.")
	}

	if err != nil {
		return makeError(actions_pb.SetGivesDividendsResponse_InternalServerError, getInternalErrorMessage(err))
	}

	resp.StatusCode = 0
	resp.StatusMessage = "GivesDividends set succesfully."

	return resp, nil
}

func (d *dalalActionService) SetBankruptcy(ctx context.Context, req *actions_pb.SetBankruptcyRequest) (*actions_pb.SetBankruptcyResponse, error) {
	var l = logger.WithFields(logrus.Fields{
		"method":        "SetBankruptcy",
		"param_session": fmt.Sprintf("%+v", ctx.Value("session")),
		"param_req":     fmt.Sprintf("%+v", req),
	})

	l.Infof("Request for setting bankruptcy")

	resp := &actions_pb.SetBankruptcyResponse{}
	makeError := func(st actions_pb.SetBankruptcyResponse_StatusCode, msg string) (*actions_pb.SetBankruptcyResponse, error) {
		resp.StatusCode = st
		resp.StatusMessage = msg
		return resp, nil
	}

	userId := getUserId(ctx)
	if !models.IsAdminAuth(userId) {
		return makeError(actions_pb.SetBankruptcyResponse_NotAdminUserError, "User is not admin")
	}

	if !models.IsMarketOpen() {
		return makeError(actions_pb.SetBankruptcyResponse_MarketClosedError, "Market Is closed. You cannot set bankruptcy for stocks now.")
	}

	stockID := req.GetStockId()
	isBankrupt := req.GetIsBankrupt()

	err := models.SetBankruptcy(stockID, isBankrupt)

	if err == models.InvalidStockError {
		return makeError(actions_pb.SetBankruptcyResponse_InvalidStockIdError, "Invalid stock id provided.")
	}

	if err != nil {
		return makeError(actions_pb.SetBankruptcyResponse_InternalServerError, getInternalErrorMessage(err))
	}

	resp.StatusCode = 0
	resp.StatusMessage = "Bankruptcy set succesfully."

	return resp, nil
}

func (d *dalalActionService) InspectUser(ctx context.Context, req *actions_pb.InspectUserRequest) (*actions_pb.InspectUserResponse, error) {

	var l = logger.WithFields(logrus.Fields{
		"method":        "InspectUser",
		"param_session": fmt.Sprintf("%+v", ctx.Value("session")),
		"param_req":     fmt.Sprintf("%+v", req),
	})

	resp := &actions_pb.InspectUserResponse{}

	makeError := func(st actions_pb.InspectUserResponse_StatusCode, msg string) (*actions_pb.InspectUserResponse, error) {
		resp.StatusCode = st
		resp.StatusMessage = msg
		return resp, nil
	}
	userId := getUserId(ctx)
	if !models.IsAdminAuth(userId) {
		return makeError(actions_pb.InspectUserResponse_NotAdminUserError, "User is not admin")
	}

	results, err := models.GetInspectUserDetails(req.UserId, req.TransactionType, req.Day)
	if err != nil {
		l.Errorf("Request failed due to %+v: ", err)
		return makeError(actions_pb.InspectUserResponse_InternalServerError, getInternalErrorMessage(err))
	}

	for _, result := range results {
		resp.List = append(resp.List, result.ToProto())
	}
	resp.StatusMessage = "Done"
	resp.StatusCode = actions_pb.InspectUserResponse_OK
	return resp, nil
}

func (d *dalalActionService) BlockUser(ctx context.Context, req *actions_pb.BlockUserRequest) (*actions_pb.BlockUserResponse, error) {
	var l = logger.WithFields(logrus.Fields{
		"method":        "BlockUser",
		"param_session": fmt.Sprintf("%+v", ctx.Value("session")),
		"param_req":     fmt.Sprintf("%+v", req),
	})

	l.Infof("Block User Requested")

	userId := req.GetUserId()

	resp := &actions_pb.BlockUserResponse{}
	makeError := func(st actions_pb.BlockUserResponse_StatusCode, msg string) (*actions_pb.BlockUserResponse, error) {
		resp.StatusCode = st
		resp.StatusMessage = msg
		return resp, nil
	}

	requesterId := getUserId(ctx)
	if !models.IsAdminAuth(requesterId) {
		return makeError(actions_pb.BlockUserResponse_NotAdminUserError, "User is not admin")
	}

	err := models.SetBlockUser(userId, true)

	if err == models.InternalServerError {
		return makeError(actions_pb.BlockUserResponse_InternalServerError, getInternalErrorMessage(err))
	} else if err == models.UserNotFoundError {
		return makeError(actions_pb.BlockUserResponse_InvalidUserIDError, "Invalid userId requested.")
	}

	return makeError(actions_pb.BlockUserResponse_OK, "User blocked successfully.")
}

func (d *dalalActionService) UnBlockUser(ctx context.Context, req *actions_pb.UnblockUserRequest) (*actions_pb.UnblockUserResponse, error) {
	var l = logger.WithFields(logrus.Fields{
		"method":        "UnblockUser",
		"param_session": fmt.Sprintf("%+v", ctx.Value("session")),
		"param_req":     fmt.Sprintf("%+v", req),
	})

	l.Infof("UnBlock User Requested")

	userId := req.GetUserId()

	resp := &actions_pb.UnblockUserResponse{}
	makeError := func(st actions_pb.UnblockUserResponse_StatusCode, msg string) (*actions_pb.UnblockUserResponse, error) {
		resp.StatusCode = st
		resp.StatusMessage = msg
		return resp, nil
	}

	requesterId := getUserId(ctx)
	if !models.IsAdminAuth(requesterId) {
		return makeError(actions_pb.UnblockUserResponse_NotAdminUserError, "User is not admin")
	}

	err := models.SetBlockUser(userId, false)

	if err == models.InternalServerError {
		return makeError(actions_pb.UnblockUserResponse_InternalServerError, getInternalErrorMessage(err))
	} else if err == models.UserNotFoundError {
		return makeError(actions_pb.UnblockUserResponse_InvalidUserIDError, "Invalid userId requested.")
	}

	return makeError(actions_pb.UnblockUserResponse_OK, "User unblocked successfully.")
}

func (d *dalalActionService) UnBlockAllUsers(ctx context.Context, req *actions_pb.UnblockAllUsersRequest) (*actions_pb.UnblockAllUsersResponse, error) {
	var l = logger.WithFields(logrus.Fields{
		"method":        "UnblockAllUsers",
		"param_session": fmt.Sprintf("%+v", ctx.Value("session")),
		"param_req":     fmt.Sprintf("%+v", req),
	})

	l.Infof("UnBlockAllUsers Requested")

	resp := &actions_pb.UnblockAllUsersResponse{}
	makeError := func(st actions_pb.UnblockAllUsersResponse_StatusCode, msg string) (*actions_pb.UnblockAllUsersResponse, error) {
		resp.StatusCode = st
		resp.StatusMessage = msg
		return resp, nil
	}

	requesterId := getUserId(ctx)
	if !models.IsAdminAuth(requesterId) {
		return makeError(actions_pb.UnblockAllUsersResponse_NotAdminUserError, "User is not admin")
	}

	err := models.UnBlockAllUsers()

	if err == models.InternalServerError {
		return makeError(actions_pb.UnblockAllUsersResponse_InternalServerError, getInternalErrorMessage(err))
	}

	return makeError(actions_pb.UnblockAllUsersResponse_OK, "All users unblocked successfully.")
}
