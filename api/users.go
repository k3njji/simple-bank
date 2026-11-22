package api

import (
	"net/http"
	db "simple_bank/db/sqlc"
	"simple_bank/util"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"time"
)

type createUserRequest struct {
	Username    string `json:"username" binding:"required,alphanum"`
	FullName    string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required"`
}

type createUserResponse struct {
    Username string `json:"username"`
    FullName string `json:"full_name"`
    Email    string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if(err != nil){
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Email: req.Email,
		HashedPassword:  hashedPassword,
		FullName: req.FullName,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok{
			switch pqErr.Code.Name(){
			case "unique_violation", "unique_validation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))

			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := createUserResponse{
		Username: user.Username,
		Email: user.Email,
		FullName: user.FullName,
		CreatedAt: user.CreatedAt,
	}


	ctx.JSON(http.StatusOK, res)
}