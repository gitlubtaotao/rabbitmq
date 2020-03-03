package controllers

import (
	"github.com/kataras/iris/v12/mvc"
	"rabbitmq/iris/repositories"
	"rabbitmq/iris/service"
)

type MovieController struct {
}

func (c *MovieController) Get() mvc.View {
	res := repositories.NewMovieManager()
	movie := service.NewMovieServiceManager(res)
	MovieResult := movie.ShowMovieName()
	return mvc.View{
		Name: "movie/index.html",
		Data: MovieResult,
	}
}
