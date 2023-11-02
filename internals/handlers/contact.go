package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

type JokeResponse struct {
	Success bool `json:"success"`
	Body    []struct {
		ID        string `json:"_id"`
		Type      string `json:"type"`
		Setup     string `json:"setup"`
		Punchline string `json:"punchline"`
	} `json:"body"`
}

func (a *App) Contact(c echo.Context) error {
	url := "https://dad-jokes.p.rapidapi.com/random/joke"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("x-rapidapi-host", "dad-jokes.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "612e44964cmshe2958fd8c359484p10457ejsn7bf2d4dea352")

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

	return c.JSON(http.StatusOK, jokeResponse)
}
