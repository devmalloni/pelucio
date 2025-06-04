package wallet

import (
	"errors"
	"math/big"
	"net/http"
	"pelucio/x/xerrors"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var (
	ErrNotFound            = xerrors.New("ErrNotFound")
	ErrWalletNotFound      = xerrors.New("ErrWalletNotFound")
	ErrBadRequest          = xerrors.New("ErrBadRequest")
	ErrExternalIDHasExists = xerrors.New("ErrExternalIDHasExists")
)

type (
	handlerDependencies interface {
		ManagerProvider
	}
	Handler struct {
		d handlerDependencies
	}
	TransferModel struct {
		FromWalletID uuid.UUID
		ToWalletID   uuid.UUID
		Amount       *string
		Currency     string
	}
	AmountModel struct {
		Amount   *string `json:"amount,omitempty"`
		Currency string  `json:"currency,omitempty"`
	}
	CreateWalletModel struct {
		ID         *uuid.UUID `json:"id,omitempty"`
		ExternalID string     `json:"externalID,omitempty"`
		Balances   *[]string  `json:"balances,omitempty"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		Version   uuid.UUID  `json:"version,omitempty"`
	}
)

func NewHandler(d handlerDependencies) *Handler {
	return &Handler{
		d,
	}
}

// Transfer godoc
// @Summary      		Transfer Transaction
// @Description  		transfer funds from one wallet to another
// @Tags         		wallet
// @Accept       		json
// @Produce      		json
// @Param      			model     		body   		TransferModel 		true	"Transfer data"
// @Success      		200
// @Failure      		400
// @Router       		/v1/open/wallet/transfer [post]
func (p *Handler) Transfer(c *gin.Context) error {
	var t TransferModel
	err := c.BindJSON(&t)
	if err != nil {
		return err
	}

	amount, ok := new(big.Int).SetString(*t.Amount, 10)
	if !ok {
		return ErrBadRequest.WithDescription("Cannot cast amount string to big.Int")
	}

	return p.d.WalletManager().Transfer(c, t.FromWalletID, t.ToWalletID, amount, WalletCurrency(t.Currency))
}

// GetWallet godoc
// @Summary      		Get Wallet
// @Description  		Get wallet infos by ID
// @Tags         		wallet
// @Accept       		json
// @Produce      		json
// @Param        		id  			path		string				true	"Wallet id"
// @Success      		200 			{object}   	WalletResponse				"wallet"
// @Failure      		400
// @Router       		/v1/open/wallet/{id} [get]
func (p *Handler) WalletByID(c *gin.Context) error {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, errors.New("WalletNotFound"))
		return errors.New("WalletNotFound")
	}
	w, err := p.d.WalletManager().WalletByID(c, uuid.FromStringOrNil(id))
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, w.ToWalletResponse())
	return nil
}

// GetWallets godoc
// @Summary      		Get Wallets
// @Description  		Get wallets
// @Tags         		wallet
// @Accept       		json
// @Produce      		json
// @Success      		200 			{array}   	WalletResponse				"wallet"
// @Failure      		400
// @Router       		/v1/admin/wallets [get]
func (p *Handler) GetWallets(c *gin.Context) error {
	w, err := p.d.WalletManager().GetWallets(c)
	if err != nil {
		return err
	}

	wr := []*WalletResponse{}

	for _, ww := range w {
		wr = append(wr, ww.ToWalletResponse())
	}

	c.JSON(http.StatusOK, wr)
	return nil
}

// GetWalletByExternalID godoc
// @Summary      		Get Wallet by externalID
// @Description  		Get wallet infos by ExternalID
// @Tags         		wallet
// @Accept       		json
// @Produce      		json
// @Param        		id  			path		string				true	"External id"
// @Success      		200 			{object}   	WalletResponse				"wallet"
// @Failure      		400
// @Router       		/v1/open/wallet/external/{id} [get]
func (p *Handler) WalletByExternalID(c *gin.Context) error {
	id := c.Param("id")
	if id == "" {
		return ErrWalletNotFound
	}
	w, err := p.d.WalletManager().WalletByExternalID(c, id)
	if err != nil {
		return err
	}

	if w == nil {
		return ErrWalletNotFound
	}

	c.JSON(http.StatusOK, w.ToWalletResponse())
	return nil
}

// GetWalletRecords godoc
// @Summary      		Get Wallet records
// @Description  		Get wallet records infos by ID
// @Tags         		wallet
// @Accept       		json
// @Produce      		json
// @Param        		id  			path		string				true	"Wallet id"
// @Success      		200 			{array}   	WalletRecord				"wallet records"
// @Failure      		400
// @Router       		/v1/open/wallet/{id}/records [get]
func (p *Handler) WalletRecords(c *gin.Context) error {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return ErrBadRequest
	}
	w, err := p.d.WalletManager().WalletRecordsByID(c, uuid.FromStringOrNil(id))
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, w)
	return nil
}

// Burn godoc
// @Summary      		Burn Transaction
// @Description  		burn funds from one wallet
// @Tags         		wallet
// @Accept       		json
// @Produce      		json
// @Param        		id  			path		string			true	"Wallet id"
// @Param      			model     		body   		AmountModel 		true	"Amount model"
// @Success      		200
// @Failure      		400
// @Router       		/v1/admin/wallet/{id}/burn [post]
func (p *Handler) Burn(c *gin.Context) error {
	var b AmountModel
	err := c.BindJSON(&b)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return err
	}

	wID := c.Param("id")
	if wID == "" {
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return ErrBadRequest
	}

	amount, ok := new(big.Int).SetString(*b.Amount, 10)
	if !ok {
		return ErrBadRequest.WithDescription("Cannot cast amount string to big.Int")
	}

	return p.d.WalletManager().Burn(c, uuid.FromStringOrNil(wID), amount, WalletCurrency(b.Currency))
}

// Mint godoc
// @Summary      		Mint Transaction
// @Description  		mint funds from one wallet
// @Tags         		wallet
// @Accept       		json
// @Produce      		json
// @Param        		id  			path		string			true	"Wallet id"
// @Param      			model     		body   		AmountModel 		true	"Amount model"
// @Success      		200
// @Failure      		400
// @Router       		/v1/admin/wallet/{id}/mint [post]
func (p *Handler) Mint(c *gin.Context) error {
	var m AmountModel
	err := c.BindJSON(&m)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return err
	}

	wID := c.Param("id")
	if wID == "" {
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return ErrBadRequest
	}

	amount, ok := new(big.Int).SetString(*m.Amount, 10)
	if !ok {
		return ErrBadRequest.WithDescription("Cannot cast amount string to big.Int")
	}

	return p.d.WalletManager().Mint(c, uuid.FromStringOrNil(wID), amount, WalletCurrency(m.Currency))
}

// Lock godoc
// @Summary      		Lock Transaction
// @Description  		lock funds from one wallet
// @Tags         		wallet
// @Accept       		json
// @Produce      		json
// @Param        		id  			path		string			true	"Wallet id"
// @Param      			model     		body   		AmountModel 		true	"Amount model"
// @Success      		200
// @Failure      		400
// @Router       		/v1/admin/wallet/{id}/lock [post]
func (p *Handler) Lock(c *gin.Context) error {
	var m AmountModel
	err := c.BindJSON(&m)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return err
	}

	wID := c.Param("id")
	if wID == "" {
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return ErrBadRequest
	}

	amount, ok := new(big.Int).SetString(*m.Amount, 10)
	if !ok {
		return ErrBadRequest.WithDescription("Cannot cast amount string to big.Int")
	}

	return p.d.WalletManager().Lock(c, uuid.FromStringOrNil(wID), amount, WalletCurrency(m.Currency))
}

// Unlock godoc
// @Summary      		Unlock Transaction
// @Description  		unlock funds from one wallet
// @Tags         		wallet
// @Accept       		json
// @Produce      		json
// @Param        		id  			path		string			true	"Wallet id"
// @Param      			model     		body   		AmountModel 		true	"Amount model"
// @Success      		200
// @Failure      		400
// @Router       		/v1/admin/wallet/{id}/unlock [post]
func (p *Handler) Unlock(c *gin.Context) error {
	var m AmountModel
	err := c.BindJSON(&m)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return err
	}

	wID := c.Param("id")
	if wID == "" {
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return ErrBadRequest
	}

	amount, ok := new(big.Int).SetString(*m.Amount, 10)
	if !ok {
		return ErrBadRequest.WithDescription("Cannot cast amount string to big.Int")
	}

	return p.d.WalletManager().Unlock(c, uuid.FromStringOrNil(wID), amount, WalletCurrency(m.Currency))
}

// MintAndLock godoc
// @Summary      		Mint and Lock Transaction
// @Description  		mint funds and lock that same funds from one wallet
// @Tags         		wallet
// @Accept       		json
// @Produce      		json
// @Param        		id  			path		string			true	"Wallet id"
// @Param      			model     		body   		AmountModel 		true	"Amount model"
// @Success      		200
// @Failure      		400
// @Router       		/v1/admin/wallet/{id}/mintandlock [post]
func (p *Handler) MintAndLock(c *gin.Context) error {
	var m AmountModel
	err := c.BindJSON(&m)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return err
	}

	wID := c.Param("id")
	if wID == "" {
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return ErrBadRequest
	}

	amount, ok := new(big.Int).SetString(*m.Amount, 10)
	if !ok {
		return ErrBadRequest.WithDescription("Cannot cast amount string to big.Int")
	}

	return p.d.WalletManager().MintAndLock(c, uuid.FromStringOrNil(wID), amount, WalletCurrency(m.Currency))
}

// UnlockAndBurn godoc
// @Summary      		Unlock and burn Transaction
// @Description  		unlock funds and burn that same funds from one wallet
// @Tags         		wallet
// @Accept       		json
// @Produce      		json
// @Param        		id  			path		string			true	"Wallet id"
// @Param      			model     		body   		AmountModel 		true	"Amount model"
// @Success      		200
// @Failure      		400
// @Router       		/v1/admin/wallet/{id}/mintandlock [post]
func (p *Handler) UnlockAndBurn(c *gin.Context) error {
	var m AmountModel
	err := c.BindJSON(&m)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return err
	}

	wID := c.Param("id")
	if wID == "" {
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return ErrBadRequest
	}

	amount, ok := new(big.Int).SetString(*m.Amount, 10)
	if !ok {
		return ErrBadRequest.WithDescription("Cannot cast amount string to big.Int")
	}

	return p.d.WalletManager().UnlockAndBurn(c, uuid.FromStringOrNil(wID), amount, WalletCurrency(m.Currency))
}

// CreateWallet godoc
// @Summary      		Create a wallet
// @Description  		create a wallet
// @Tags         		wallet
// @Accept       		json
// @Produce      		json
// @Param      			model     		body   		CreateWalletModel 		true	"Create wallet model"
// @Success      		200 			{object}   	WalletResponse					"wallet"
// @Failure      		400
// @Router       		/v1/admin/wallet [post]
func (p *Handler) CreateWallet(c *gin.Context) error {
	var w CreateWalletModel
	err := c.BindJSON(&w)
	if err != nil {
		return err
	}

	ww, err := p.d.WalletManager().d.WalletPersister().FindWalletByExternalID(c, w.ExternalID)
	if err != nil {
		return err
	}

	if ww != nil && ww.ID.String() != "" {
		return ErrExternalIDHasExists
	}

	if w.ID == nil {
		id := uuid.NewV4()
		w.ID = &id
	}

	balance := make(map[WalletCurrency]*big.Int)
	lockedBalance := make(map[WalletCurrency]*big.Int)

	if w.Balances != nil {
		for _, v := range *w.Balances {
			balance[WalletCurrency(v)] = big.NewInt(0)
		}
	}

	if w.Balances != nil {
		for _, v := range *w.Balances {
			lockedBalance[WalletCurrency(v)] = big.NewInt(0)
		}
	}

	wallet := &Wallet{
		ID:            *w.ID,
		Balance:       balance,
		LockedBalance: lockedBalance,
		Version:       w.Version,
		ExternalID:    w.ExternalID,
		CreatedAt:     time.Now(),
	}

	err = p.d.WalletManager().d.WalletPersister().SaveWallet(c, []*Wallet{wallet}, make([]*WalletRecord, 0), make([]*WalletTransaction, 0))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return err
	}

	c.JSON(http.StatusOK, wallet.ToWalletResponse())
	return nil
}

func (p *Handler) RegisterOpenRoutes(r *gin.RouterGroup) {
	r.POST("/wallet/transfer", xerrors.HandleWithError(p.Transfer))
	r.GET("/wallet/:id", xerrors.HandleWithError(p.WalletByID))
	r.GET("/wallet/external/:id", xerrors.HandleWithError(p.WalletByExternalID))
	r.GET("/wallet/:id/records", xerrors.HandleWithError(p.WalletRecords))
}

func (p *Handler) RegisterAdminRoutes(r *gin.RouterGroup) {
	r.POST("/wallet/:id/burn", xerrors.HandleWithError(p.Burn))
	r.POST("/wallet/:id/mint", xerrors.HandleWithError(p.Mint))
	r.POST("/wallet/:id/lock", xerrors.HandleWithError(p.Lock))
	r.POST("/wallet/:id/unlock", xerrors.HandleWithError(p.Unlock))
	r.POST("/wallet/:id/mintandlock", xerrors.HandleWithError(p.MintAndLock))
	r.POST("/wallet/:id/unlockandburn", xerrors.HandleWithError(p.UnlockAndBurn))
	r.POST("/wallet", xerrors.HandleWithError(p.CreateWallet))
	r.GET("/wallets", xerrors.HandleWithError(p.GetWallets))
}
