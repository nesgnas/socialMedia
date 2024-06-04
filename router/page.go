package router

import (
	"github.com/gin-gonic/gin"
	"socialMedia/controler"
	"socialMedia/middleware"
)

func Page(incomingRouter *gin.Engine) {
	incomingRouter.Use(middleware.CORSMiddleware())
	incomingRouter.POST("/initiateUser", controler.NewUser())
	incomingRouter.POST("/createPost/:userid", controler.NewPost())
	incomingRouter.POST("/commentPost/:sender", controler.NewCommentPost())
	incomingRouter.POST("/likePost/:userid", controler.NewLikePost())
	incomingRouter.PUT("/deletePost/:userid", controler.DeletePost())
	incomingRouter.PUT("/deleteComment/:userid", controler.DeleteComment())
	incomingRouter.GET("/getPost/:userid", controler.GetPost())
}
