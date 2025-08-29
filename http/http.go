package https

import (
	"auth-system/interfaces"
	request "auth-system/models/requests"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type Http struct {
	uc interfaces.UseCase
}

func NewHttp(usecase interfaces.UseCase) (Http) {
	return Http{uc: usecase}
}



func (h Http) Register(c *gin.Context) {
	requestBody := &request.RegisterRequest{}
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(400, gin.H {
			"error": err.Error(),
		})

		return
	}

	validate := validator.New()
	err = validate.Struct(requestBody)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid Request Body",
		})

		return
	}

	err = h.uc.Register(requestBody)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H {
		"message": "Successfully registered user",
	})
}

func (h Http) Login(c *gin.Context) {
	requestBody := &request.LoginRequest{}
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(400, gin.H {
			"error": err.Error(),
		})

		return
	}

	validate := validator.New()
	err = validate.Struct(requestBody)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid Request Body",
		})

		return
	}

	token, err := h.uc.Login(requestBody)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"token": token,
		"message": "User successfully logged in",
	})
}

func (h Http) Logout(c *gin.Context) {
	tokenID, _ := c.Get("tokenID")
	ID := tokenID.(float64)
	err := h.uc.Logout(int(ID))

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "User logged out successfully",
	})
}

func (h Http) Forgot(c *gin.Context) {
	requestBody := &request.ForgotRequest{}
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(400, gin.H {
			"error": err.Error(),
		})

		return
	}

	validate := validator.New()
	err = validate.Struct(requestBody)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid Request Body",
		})

		return
	}

	err = h.uc.Forgot(requestBody)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "recovery email sent", 
	})
}

func (h Http) Reset(c *gin.Context) {
		requestBody := &request.ResetRequest{}
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(400, gin.H {
			"error": err.Error(),
		})

		return
	}

	validate := validator.New()
	err = validate.Struct(requestBody)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid Request Body",
		})

		return
	}

	token := c.Query("token")
	err = h.uc.Reset(requestBody, token)

		if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "Password successfully changed",
	})
}
