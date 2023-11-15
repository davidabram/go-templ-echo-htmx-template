package main

import (
	"context"
	"log"
	"os"

	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"

	"davidabram/go-templ-echo-htmx-template/internals/handlers"

	"github.com/donseba/go-htmx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := &handlers.App{
		HTMX: htmx.New(),
	}

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(HtmxMiddleware)

	e.GET("/", app.Hello)
	e.GET("/about", app.About)
	e.GET("/books", app.BooksTable)
	e.GET("/charts", app.Charts)
	e.GET("/contact", app.Contact)

	e.Static("/", "dist")
	e.Static("/fonts", "static/fonts")

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))

	clerkClient := clerk.NewClient(os.Getenv("CLERK_SECRET_KEY"))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	user, err := clerk.GetUser(r)
	if err != nil {

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	templates.ExecuteTemplate(w, "layout", user != nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	clerk.Login(w, r)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clerk.Logout(w, r)
}

func HtmxMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		hxh := htmx.HxRequestHeader{
			HxBoosted:               htmx.HxStrToBool(c.Request().Header.Get("HX-Boosted")),
			HxCurrentURL:            c.Request().Header.Get("HX-Current-URL"),
			HxHistoryRestoreRequest: htmx.HxStrToBool(c.Request().Header.Get("HX-History-Restore-Request")),
			HxPrompt:                c.Request().Header.Get("HX-Prompt"),
			HxRequest:               htmx.HxStrToBool(c.Request().Header.Get("HX-Request")),
			HxTarget:                c.Request().Header.Get("HX-Target"),
			HxTriggerName:           c.Request().Header.Get("HX-Trigger-Name"),
			HxTrigger:               c.Request().Header.Get("HX-Trigger"),
		}

		ctx = context.WithValue(ctx, htmx.ContextRequestHeader, hxh)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
