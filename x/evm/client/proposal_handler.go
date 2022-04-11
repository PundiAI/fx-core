package client

//var (
//	InitEvmProposalHandler = govclient.NewProposalHandler(cli.InitEvmProposalCmd, InitEvmProposalRESTHandler)
//)
//
//// InitEvmProposalRequest defines a request for a new init evm proposal.
//type InitEvmProposalRequest struct {
//	BaseReq         rest.BaseReq           `json:"base_req" yaml:"base_req"`
//	Title           string                 `json:"title" yaml:"title"`
//	Description     string                 `json:"description" yaml:"description"`
//	Deposit         sdk.Coins              `json:"deposit" yaml:"deposit"`
//	EvmParams       *types.Params          `json:"evm_params" yaml:"evm_params"`
//	FeemarketParams *feemarkettypes.Params `json:"feemarket_params" yaml:"feemarket_params"`
//	Erc20Params     *types.Erc20Params     `json:"erc20_params" yaml:"erc20_params"`
//	Metadata        []banktypes.Metadata   `json:"metadata" yaml:"metadata"`
//}
//
//func InitEvmProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
//	return govrest.ProposalRESTHandler{
//		SubRoute: types.ModuleName,
//		Handler:  newInitEvmProposalHandler(clientCtx),
//	}
//}

// nolint: dupl
//func newInitEvmProposalHandler(clientCtx client.Context) http.HandlerFunc {
//return func(w http.ResponseWriter, r *http.Request) {
//	var req InitEvmProposalRequest
//
//	if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
//		return
//	}
//
//	req.BaseReq = req.BaseReq.Sanitize()
//	if !req.BaseReq.ValidateBasic(w) {
//		return
//	}
//
//	fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
//	if rest.CheckBadRequestError(w, err) {
//		return
//	}
//
//	content := types.NewInitEvmProposal(req.Title, req.Description, req.EvmParams, req.FeemarketParams, req.Erc20Params, req.Metadata)
//	msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, fromAddr)
//	if rest.CheckBadRequestError(w, err) {
//		return
//	}
//
//	if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
//		return
//	}
//
//	tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
//}
//}
