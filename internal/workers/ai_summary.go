package workers

import (
	"fmt"
	"log"
	"time"

	"github.com/ankitsingh/urlshortener/internal/models"
	"github.com/ankitsingh/urlshortener/internal/repository"
	"github.com/ankitsingh/urlshortener/internal/services"
)

// StartAISummaryWorker starts a goroutine that runs every 24 hours
// For testing purposes, we can trigger it immediately or run it more frequently.
func StartAISummaryWorker() {
	go func() {
		// Run once on startup after a small delay (useful for testing)
		time.Sleep(5 * time.Second)
		log.Println("Running initial AI Summary job...")
		runAISummaryJob()

		// Then run daily
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("Running scheduled AI Summary job...")
			runAISummaryJob()
		}
	}()
}

func runAISummaryJob() {
	// 1. Get the start of the current week (e.g., Sunday 00:00)
	now := time.Now()
	daysSinceSunday := int(now.Weekday())
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-daysSinceSunday, 0, 0, 0, 0, now.Location())
	sevenDaysAgo := now.Add(-7 * 24 * time.Hour)

	// 2. Find links that had clicks in the last 7 days
	var activeLinkIDs []uint
	err := repository.DB.Model(&models.ClickEvent{}).
		Where("clicked_at >= ?", sevenDaysAgo).
		Distinct("link_id").
		Pluck("link_id", &activeLinkIDs).Error

	if err != nil {
		log.Printf("Worker Error: Failed to fetch active links: %v", err)
		return
	}

	for _, linkID := range activeLinkIDs {
		// 3. Check if a summary already exists for this link and this weekStart
		var count int64
		repository.DB.Model(&models.AISummary{}).
			Where("link_id = ? AND week_start = ?", linkID, weekStart).
			Count(&count)

		if count > 0 {
			continue // Already processed this week
		}

		// 4. Fetch aggregated stats for this link
		var totalClicks int64
		repository.DB.Model(&models.ClickEvent{}).Where("link_id = ? AND clicked_at >= ?", linkID, sevenDaysAgo).Count(&totalClicks)

		type agg struct {
			Name  string
			Count int64
		}
		var topCountries []agg
		repository.DB.Raw("SELECT country as name, COUNT(*) as count FROM click_events WHERE link_id = ? AND clicked_at >= ? GROUP BY country ORDER BY count DESC LIMIT 3", linkID, sevenDaysAgo).Scan(&topCountries)

		var topDevices []agg
		repository.DB.Raw("SELECT device_type as name, COUNT(*) as count FROM click_events WHERE link_id = ? AND clicked_at >= ? GROUP BY device_type ORDER BY count DESC LIMIT 3", linkID, sevenDaysAgo).Scan(&topDevices)

		// 5. Construct the Prompt
		prompt := fmt.Sprintf("Analyze the traffic for this URL over the past 7 days. Total Clicks: %d. ", totalClicks)
		prompt += "Top Countries: "
		for _, c := range topCountries {
			prompt += fmt.Sprintf("%s (%d), ", c.Name, c.Count)
		}
		prompt += "Top Devices: "
		for _, d := range topDevices {
			prompt += fmt.Sprintf("%s (%d), ", d.Name, d.Count)
		}
		prompt += "Please summarize these trends in a short, engaging paragraph."

		// 6. Call Groq
		summaryText, err := services.GenerateSummary(prompt)
		if err != nil {
			log.Printf("Worker Error: Failed to generate summary for link %d: %v", linkID, err)
			continue
		}

		// 7. Save to Database
		summary := models.AISummary{
			LinkID:      linkID,
			WeekStart:   weekStart,
			SummaryText: summaryText,
		}
		if err := repository.DB.Create(&summary).Error; err != nil {
			log.Printf("Worker Error: Failed to save summary for link %d: %v", linkID, err)
		} else {
			log.Printf("Successfully generated and cached AI summary for link %d", linkID)
		}
	}
}
