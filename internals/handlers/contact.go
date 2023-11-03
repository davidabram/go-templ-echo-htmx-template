package handlers

import (
	"context"
	"davidabram/go-templ-echo-htmx-template/internals/templates"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

type JokeResponse struct {
	Failed    bool   `json:"error"`
	Setup     string `json:"setup"`
	Punchline string `json:"delivery"`
}

func (a *App) Contact(c echo.Context) error {
	url := "https://v2.jokeapi.dev/joke/Any?type=twopart"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var jokeResponse JokeResponse
	if err := json.Unmarshal(body, &jokeResponse); err != nil {
		return err
	}

	page := &templates.Page{
		Title:   "Contact",
		Boosted: true,
	}

	joke := &templates.Joke{
		Setup:     jokeResponse.Setup,
		Punchline: jokeResponse.Punchline,
	}

	components := templates.Contact(page, joke)
	return components.Render(context.Background(), c.Response().Writer)
}
