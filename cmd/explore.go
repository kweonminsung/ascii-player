package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/gocolly/colly/v2"
	"github.com/kweonminsung/console-cinema/pkg/tui/player"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var exploreCmd = &cobra.Command{
	Use:   "explore",
	Short: "Explore youtube videos",
	Long:  `Explore youtube videos with fuzzy finder`,
	Run:   explore,
}

type Video struct {
	Title       string
	URL         string
	Thumbnail   string
	Uploader    string
	PublishedAt string
}

func explore(cmd *cobra.Command, args []string) {
	app := tview.NewApplication()
	inputField := tview.NewInputField().
		SetLabel("Search YouTube: ")

	videoList := tview.NewList()
	suggestionsList := tview.NewList().
		ShowSecondaryText(false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(inputField, 0, 1, true).
			AddItem(suggestionsList, 0, 1, false), 0, 1, true).
		AddItem(videoList, 0, 2, false)

	inputField.SetChangedFunc(func(text string) {
		if text == "" {
			suggestionsList.Clear()
			return
		}
		suggestions := getSuggestions(text)
		suggestionsList.Clear()
		for _, s := range suggestions {
			suggestionsList.AddItem(s, "", 0, nil)
		}
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			var query string
			if inputField.HasFocus() {
				query = inputField.GetText()
			} else if suggestionsList.HasFocus() {
				query, _ = suggestionsList.GetItemText(suggestionsList.GetCurrentItem())
			}

			if query != "" {
				suggestionsList.Clear()
				videoList.Clear().AddItem("Searching...", "", 0, nil)
				app.SetFocus(videoList)

				go func() {
					videos := searchYoutube(query)
					app.QueueUpdateDraw(func() {
						updateVideoList(videoList, videos, app)
					})
				}()
			}
			return nil // Consume the event
		case tcell.KeyTab:
			if suggestionsList.HasFocus() && suggestionsList.GetItemCount() > 0 {
				query, _ := suggestionsList.GetItemText(suggestionsList.GetCurrentItem())
				inputField.SetText(query)
				app.SetFocus(inputField)
			}
			return nil // Consume the event
		case tcell.KeyDown:
			if inputField.HasFocus() && suggestionsList.GetItemCount() > 0 {
				app.SetFocus(suggestionsList)
				return nil // Consume the event
			}
		case tcell.KeyUp:
			if suggestionsList.HasFocus() && suggestionsList.GetCurrentItem() == 0 {
				app.SetFocus(inputField)
				return nil // Consume the event
			}
		}
		return event
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		log.Fatal(err)
	}
}

func getSuggestions(query string) []string {
	c := colly.NewCollector()
	var suggestions []string

	c.OnResponse(func(r *colly.Response) {
		var result []interface{}
		if err := json.Unmarshal(r.Body, &result); err != nil {
			return
		}
		if len(result) > 1 {
			if s, ok := result[1].([]interface{}); ok {
				for _, item := range s {
					if str, ok := item.(string); ok {
						suggestions = append(suggestions, str)
					}
				}
			}
		}
	})

	suggestURL := fmt.Sprintf("https://suggestqueries.google.com/complete/search?client=firefox&q=%s", url.QueryEscape(query))
	c.Visit(suggestURL)

	return suggestions
}

func searchYoutube(query string) []Video {
	c := colly.NewCollector()
	videos := []Video{}

	c.OnHTML("ytd-video-renderer", func(e *colly.HTMLElement) {
		videoURL := "https://www.youtube.com" + e.ChildAttr("a#video-title", "href")
		// Skip shorts and other non-standard videos
		if !strings.Contains(videoURL, "watch?v=") {
			return
		}

		video := Video{
			Title:       e.ChildAttr("a#video-title", "title"),
			URL:         videoURL,
			Thumbnail:   e.ChildAttr("img.yt-core-image", "src"),
			Uploader:    e.ChildText("a.yt-simple-endpoint.style-scope.yt-formatted-string"),
			PublishedAt: e.ChildText("span.style-scope.ytd-video-meta-block"),
		}
		videos = append(videos, video)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	searchURL := fmt.Sprintf("https://www.youtube.com/results?search_query=%s", url.QueryEscape(query))
	c.Visit(searchURL)

	return videos
}

func updateVideoList(list *tview.List, videos []Video, app *tview.Application) {
	list.Clear()
	for _, video := range videos {
		videoCopy := video // Create a copy to avoid closure issues
		// Displaying thumbnail URL as secondary text.
		secondaryText := fmt.Sprintf("By: %s | Published: %s | Thumbnail: %s", videoCopy.Uploader, videoCopy.PublishedAt, videoCopy.Thumbnail)
		list.AddItem(videoCopy.Title, secondaryText, 0, func() {
			app.Stop()
			playVideo(videoCopy.URL)
		})
	}
}

func playVideo(url string) {
	// This logic is borrowed from cmd/youtube.go
	fps, _ := youtubeCmd.Flags().GetInt("fps")
	loop, _ := youtubeCmd.Flags().GetBool("loop")
	color, _ := youtubeCmd.Flags().GetBool("color")
	mode, _ := youtubeCmd.Flags().GetString("mode")

	fmt.Printf("Starting %s player for YouTube video: %s\n", mode, url)
	fmt.Printf("Settings - FPS: %d, Loop: %t, Color: %t, Mode: %s\n", fps, loop, color, mode)

	// Create and start TUI player
	player := player.NewPlayer(url, fps, loop, color, mode)

	err := player.Play()
	if err != nil {
		fmt.Printf("Error during playback: %v\n", err)
		return
	}
}

func init() {
	youtubeCmd.AddCommand(exploreCmd)
}
