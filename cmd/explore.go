package cmd

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/kweonminsung/console-cinema/pkg/tui/player"
	"github.com/kweonminsung/console-cinema/pkg/youtube"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var exploreCmd = &cobra.Command{
	Use:   "explore",
	Short: "Explore youtube videos",
	Long:  `Explore youtube videos with fuzzy finder`,
	Run:   explore,
}

func explore(cmd *cobra.Command, args []string) {
	app := tview.NewApplication()
	inputField := tview.NewInputField().
		SetLabel("Search YouTube: ").
		SetLabelColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorWhite).
		SetFieldTextColor(tcell.ColorBlack)

	videoList := tview.NewList()
	suggestionsList := tview.NewList().
		ShowSecondaryText(false)
	suggestionsList.SetBorderPadding(0, 1, 16, 0)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(inputField, 1, 0, true).
			AddItem(suggestionsList, 0, 1, false), 0, 1, true).
		AddItem(videoList, 0, 2, false)

	inputField.SetChangedFunc(func(text string) {
		if text == "" {
			suggestionsList.Clear()
			return
		}
		suggestions := youtube.GetSuggestions(text)
		suggestionsList.Clear()
		for _, s := range suggestions {
			suggestionsList.AddItem(s, "", 0, nil)
		}
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if videoList.HasFocus() {
			if event.Key() == tcell.KeyUp && videoList.GetCurrentItem() == 0 {
				app.SetFocus(inputField)
				return nil
			}
			return event
		}
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
					videos := youtube.SearchYoutube(query)
					app.QueueUpdateDraw(func() {
						updateVideoList(videoList, videos, app)
					})
				}()
				return nil // Consume the event
			}
			return event
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

func updateVideoList(list *tview.List, videos []youtube.Video, app *tview.Application) {
	list.Clear()
	if len(videos) == 0 {
		list.AddItem("No videos found", "", 0, nil)
		return
	}
	for _, video := range videos {
		videoCopy := video // Create a copy to avoid closure issues
		secondaryText := fmt.Sprintf("By: %s | Views: %s | Length: %s | Published: %s", videoCopy.Uploader, videoCopy.Views, videoCopy.Length, videoCopy.PublishedAt)
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
