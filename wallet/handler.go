package wallet

import "github.com/gin-gonic/gin"

type (
	handlerDependencies interface {
		ManagerProvider
	}
	Handler struct {
		d handlerDependencies
	}
)

// /v1/open/wallet/:id/transfer
func (p *Handler) Transfer(c *gin.Context) error {
	return nil
}

// /v1/admin/wallet/:id/burn
func (p *Handler) Burn(c *gin.Context) error {
	return nil
}

// /v1/admin/wallet/:id/mint
func (p *Handler) Mint(c *gin.Context) error {
	return nil
}

// /v1/admin/wallet/:id/transfer
func (p *Handler) Lock(c *gin.Context) error {
	return nil
}

func (p *Handler) Unlock(c *gin.Context) error {
	return nil
}

func (p *Handler) MintAndLock(c *gin.Context) error {
	return nil
}

func (p *Handler) UnlockAndBurn(c *gin.Context) error {
	return nil
}
