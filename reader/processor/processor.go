// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package processor

import (
	"fmt"
	"miniflux.app/logger"
	"miniflux.app/model"
	"miniflux.app/reader/mercury"
	"miniflux.app/reader/rewrite"
	"miniflux.app/reader/sanitizer"
	"miniflux.app/reader/scraper"
	"miniflux.app/storage"
)

// ProcessFeedEntries downloads original web page for entries and apply filters.
func ProcessFeedEntries(store *storage.Storage, feed *model.Feed) {
	for _, entry := range feed.Entries {
		if feed.Crawler {
			if !store.EntryURLExists(feed.ID, entry.URL) {
				content, err := scraper.Fetch(entry.URL, feed.ScraperRules, feed.UserAgent)
				if err != nil {
					logger.Error(`[Filter] Unable to crawl this entry: %q => %v`, entry.URL, err)
				} else if content != "" {
					// We replace the entry content only if the scraper doesn't return any error.
					entry.Content = content
				}
			}
		}

		if feed.UseMercury {
			if !store.EntryURLExists(feed.ID, entry.URL) {
				user, err := store.UserByID(feed.UserID)
				if err != nil && user.MercuryAPIURL != "" {
					content, err := mercury.Fetch(entry.URL, user.MercuryAPIURL)
					if err != nil {
						logger.Error(`[Merucry] Unable to crawl this entry: %q => %v`, entry.URL, err)
					} else if content != "" {
						// We replace the entry content only if the scraper doesn't return any error.
						entry.Content = content
					}
				}
			}
		}

		entry.Content = rewrite.Rewriter(entry.URL, entry.Content, feed.RewriteRules)

		// The sanitizer should always run at the end of the process to make sure unsafe HTML is filtered.
		entry.Content = sanitizer.Sanitize(entry.URL, entry.Content)
	}
}

// ProcessEntryWebPage downloads the entry web page and apply rewrite rules.
func ProcessEntryWebPage(entry *model.Entry, store *storage.Storage) error {
	if entry.Feed.UseMercury {
		user, err := store.UserByID(entry.UserID)
		if err != nil {
			return fmt.Errorf(`[Merucry] Unable to crawl this entry: %q => %v`, entry.URL, err)
		}
		if user.MercuryAPIURL == "" {
			return fmt.Errorf(`[Merucry] Unable to crawl this entry: %q => %v`, entry.URL, "API URL NOT SET")
		}
		content, err := mercury.Fetch(entry.URL, user.MercuryAPIURL)
		if err != nil {
			return err
		}
		if content != "" {
			entry.Content = content
		}
	} else {
		content, err := scraper.Fetch(entry.URL, entry.Feed.ScraperRules, entry.Feed.UserAgent)
		if err != nil {
			return err
		}

		content = rewrite.Rewriter(entry.URL, content, entry.Feed.RewriteRules)
		content = sanitizer.Sanitize(entry.URL, content)

		if content != "" {
			entry.Content = content
		}
	}

	return nil
}
