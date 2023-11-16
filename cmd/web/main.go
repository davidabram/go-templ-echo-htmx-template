package main

import (
	"context"
	"davidabram/go-templ-echo-htmx-template/internals/handlers"
	"fmt"
	"log"
	"os"

	"github.com/clerkinc/clerk-sdk-go/clerk"
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
	e.Use(clerkMiddleware)

	e.GET("/", app.Hello)
	e.GET("/about", app.About)
	e.GET("/books", app.BooksTable)
	e.GET("/charts", app.Charts)
	e.GET("/contact", app.Contact)

	e.Static("/", "dist")
	e.Static("/fonts", "static/fonts")

	client, err := clerk.NewClient(os.Getenv("CLERK_SECRET_KEY"))
	if err != nil {
		log.Fatalf("Error creating Clerk client: %v", err)
	}

	retrieveUsers(client)
	retrieveSessions(client)
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func retrieveUsers(client clerk.Client) {
	users, err := client.Users().ListAll(clerk.ListAllUsersParams{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Users:")
	for _, user := range users {
		fmt.Printf("%v. %v\n", user.FirstName, user.LastName)

	}
}

func retrieveSessions(client clerk.Client) {
	sessions, err := client.Sessions().ListAll()
	if err != nil {
		panic(err)
	}

	fmt.Println("\nSessions:")
	for i, session := range sessions {
		fmt.Printf("%v. %v (%v)\n", i+1, session.ID, session.Status)
	}
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

func clerkMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		isSignedIn := true

		c.Set("isSignedIn", isSignedIn)

		return next(c)
	}
}
