package router

import (
	"github.com/labstack/echo/v4"
	"github.com/matthewTechCom/progate_hackathon/controller"
)

func SetupRoutes(e *echo.Echo, boardController controller.WidgetControllerInterface) {
	e.POST("/process-board", boardController.ProcessBoard)
}
