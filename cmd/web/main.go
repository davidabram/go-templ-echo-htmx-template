package main

import (
	"context"
	"davidabram/go-templ-echo-htmx-template/internals/handlers"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/donseba/go-htmx"
	"github.com/go-jose/go-jose/v3/jwt"
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
	e.GET("/signin", app.SiginIn)

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
		if (user.FirstName == nil) || (user.LastName == nil) {
			continue
		}
		fmt.Printf("%v %v\n", *user.FirstName, *user.LastName)

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

func isDevelopmentOrStaging(c clerk.Client) bool {
	return strings.HasPrefix(c.APIKey(), "test_") || strings.HasPrefix(c.APIKey(), "sk_test_")
}

func isProduction(c clerk.Client) bool {
	return !isDevelopmentOrStaging(c)
}

var urlSchemeRe = regexp.MustCompile(`(^\w+:|^)\/\/`)

func isCrossOrigin(r *http.Request) bool {
	// origin contains scheme+host and optionally port (ommitted if 80 or 443)
	// ref. https://www.rfc-editor.org/rfc/rfc6454#section-6.1
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	origin = urlSchemeRe.ReplaceAllString(origin, "") // strip scheme
	if origin == "" {
		return false
	}

	// parse request's host and port, taking into account reverse proxies
	u := &url.URL{Host: r.Host}
	host := strings.TrimSpace(r.Header.Get("X-Forwarded-Host"))
	if host == "" {
		host = u.Hostname()
	}
	port := strings.TrimSpace(r.Header.Get("X-Forwarded-Port"))
	if port == "" {
		port = u.Port()
	}

	if port != "" && port != "80" && port != "443" {
		host = net.JoinHostPort(host, port)
	}

	return origin != host
}

func clerkMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		r := c.Request()

		client, err := clerk.NewClient(os.Getenv("CLERK_SECRET_KEY"))
		if err != nil {
			log.Fatalf("Error creating Clerk client: %v", err)
		}

		cookieToken, _ := r.Cookie("__session")
		clientUat, _ := r.Cookie("__client_uat")

		/* if isDevelopmentOrStaging(client) && (r.Referer() == "" || isCrossOrigin(r)) {
			c.Set("isSignedIn", false)
			// error! 401
			fmt.Println("isDevelopmentOrStaging")
			return next(c)
		} */

		if isProduction(client) && clientUat == nil {
			c.Set("isSignedIn", false)
			fmt.Println("isProduction")
			return next(c)
		}

		if clientUat != nil && clientUat.Value == "0" {
			c.Set("isSignedIn", false)
			fmt.Println("clientUat.Value == 0")
			return next(c)
		}

		if clientUat == nil {
			c.Set("isSignedIn", false)
			fmt.Println("clientUat == nil")
			// error! 401
			return next(c)
		}

		if cookieToken == nil {
			c.Set("isSignedIn", false)
			fmt.Println("cookieToken == nil")
			// error! 401
			return next(c)
		}

		var clientUatTs int64
		ts, err := strconv.ParseInt(clientUat.Value, 10, 64)
		if err == nil {
			clientUatTs = ts
		}

		claims, err := client.VerifyToken(cookieToken.Value)

		if err == nil {
			if claims.IssuedAt != nil && clientUatTs <= int64(*claims.IssuedAt) {
				fmt.Println("claims.IssuedAt != nil && clientUatTs <= int64(*claims.IssuedAt)")
				c.Set("clerk_user", &claims.Subject)
				c.Set("isSignedIn", true)
				return next(c)
			}

			c.Set("isSignedIn", false)
			// error! 401
			return next(c)
		}

		if errors.Is(err, jwt.ErrExpired) || errors.Is(err, jwt.ErrIssuedInTheFuture) {
			c.Set("isSignedIn", false)
			fmt.Println("errors.Is(err, jwt.ErrExpired) || errors.Is(err, jwt.ErrIssuedInTheFuture)")
			// error! 401
			return next(c)
		}

		c.Set("isSignedIn", false)
		fmt.Println("c.Set(\"isSignedIn\", false)")

		return next(c)
	}
}
