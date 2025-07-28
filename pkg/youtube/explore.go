package youtube

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Video struct {
	Title       string
	URL         string
	Uploader    string
	PublishedAt string
	Views       string
	Length      string
}

func GetSuggestions(query string) []string {
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

func SearchYoutube(query string) []Video {
	c := colly.NewCollector()
	videos := []Video{}

	c.OnHTML("script", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "var ytInitialData") {
			// Extract the JSON part
			jsonString := strings.TrimPrefix(e.Text, "var ytInitialData = ")
			jsonString = strings.TrimSuffix(jsonString, ";")

			var data map[string]interface{}
			if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
				log.Println("Failed to parse ytInitialData JSON:", err)
				return
			}

			contents, ok := data["contents"].(map[string]interface{})
			if !ok || contents == nil {
				log.Println("contents not found in ytInitialData")
				return
			}
			twoCol, ok := contents["twoColumnSearchResultsRenderer"].(map[string]interface{})
			if !ok || twoCol == nil {
				log.Println("twoColumnSearchResultsRenderer not found")
				return
			}
			primary, ok := twoCol["primaryContents"].(map[string]interface{})
			if !ok || primary == nil {
				log.Println("primaryContents not found")
				return
			}
			sectionList, ok := primary["sectionListRenderer"].(map[string]interface{})
			if !ok || sectionList == nil {
				log.Println("sectionListRenderer not found")
				return
			}
			sectionContents, ok := sectionList["contents"].([]interface{})
			if !ok || len(sectionContents) == 0 {
				log.Println("contents array not found or empty in sectionListRenderer")
				return
			}
			itemSectionInterface := sectionContents[0]
			if itemSectionInterface == nil {
				return
			}
			itemSection, ok := itemSectionInterface.(map[string]interface{})
			if !ok || itemSection == nil {
				log.Println("itemSection not found in contents")
				return
			}
			itemSectionRenderer, ok := itemSection["itemSectionRenderer"].(map[string]interface{})
			if !ok || itemSectionRenderer == nil {
				log.Println("itemSectionRenderer not found")
				return
			}
			videoItems, ok := itemSectionRenderer["contents"].([]interface{})
			if !ok {
				log.Println("contents not found in itemSectionRenderer")
				return
			}

			for _, item := range videoItems {
				if item == nil {
					continue
				}
				itemMap, ok := item.(map[string]interface{})
				if !ok || itemMap == nil {
					continue
				}
				videoRenderer, ok := itemMap["videoRenderer"].(map[string]interface{})
				if !ok || videoRenderer == nil {
					continue // Not a video item, could be a playlist or ad
				}

				videoId, ok := videoRenderer["videoId"].(string)
				if !ok {
					continue
				}
				titleMap, ok := videoRenderer["title"].(map[string]interface{})
				if !ok || titleMap == nil {
					continue
				}
				titleRuns, ok := titleMap["runs"].([]interface{})
				if !ok || len(titleRuns) == 0 {
					continue
				}
				titleRun, ok := titleRuns[0].(map[string]interface{})
				if !ok || titleRun == nil {
					continue
				}
				title, ok := titleRun["text"].(string)
				if !ok {
					continue
				}

				ownerTextMap, ok := videoRenderer["ownerText"].(map[string]interface{})
				if !ok || ownerTextMap == nil {
					continue
				}
				ownerTextRuns, ok := ownerTextMap["runs"].([]interface{})
				if !ok || len(ownerTextRuns) == 0 {
					continue
				}
				ownerTextRun, ok := ownerTextRuns[0].(map[string]interface{})
				if !ok || ownerTextRun == nil {
					continue
				}
				uploader, ok := ownerTextRun["text"].(string)
				if !ok {
					continue
				}

				publishedTime := "N/A"
				if publishedTimeText, ok := videoRenderer["publishedTimeText"].(map[string]interface{}); ok && publishedTimeText != nil {
					if simpleText, ok := publishedTimeText["simpleText"].(string); ok {
						publishedTime = simpleText
					}
				}

				viewCount := "N/A"
				if viewCountText, ok := videoRenderer["viewCountText"].(map[string]interface{}); ok && viewCountText != nil {
					if simpleText, ok := viewCountText["simpleText"].(string); ok {
						viewCount = simpleText
					}
				}

				length := "N/A"
				if lengthText, ok := videoRenderer["lengthText"].(map[string]interface{}); ok && lengthText != nil {
					if simpleText, ok := lengthText["simpleText"].(string); ok {
						length = simpleText
					}
				}

				video := Video{
					Title:       title,
					URL:         "https://www.youtube.com/watch?v=" + videoId,
					Uploader:    uploader,
					PublishedAt: publishedTime,
					Views:       viewCount,
					Length:      length,
				}
				videos = append(videos, video)
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	searchURL := fmt.Sprintf("https://www.youtube.com/results?search_query=%s", url.QueryEscape(query))
	c.Visit(searchURL)

	return videos
}
