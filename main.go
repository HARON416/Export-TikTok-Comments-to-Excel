package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// Removed comment limit - now extracts unlimited comments

type CommentResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	Comments   []struct {
		AuthorPin             bool   `json:"author_pin"`
		AwemeID               string `json:"aweme_id"`
		ID                    string `json:"cid"`
		CollectStat           int    `json:"collect_stat"`
		CommentLanguage       string `json:"comment_language"`
		CommentPostItemIDs    any    `json:"comment_post_item_ids"`
		CreateTime            int    `json:"create_time"`
		DiggCount             int    `json:"digg_count"`
		ImageList             any    `json:"image_list"`
		IsAuthorDigged        bool   `json:"is_author_digged"`
		IsCommentTranslatable bool   `json:"is_comment_translatable"`
		IsHighPurchaseIntent  bool   `json:"is_high_purchase_intent"`
		LabelList             any    `json:"label_list"`
		NoShow                bool   `json:"no_show"`
		ReplyComment          any    `json:"reply_comment"`
		ReplyCommentTotal     int    `json:"reply_comment_total"`
		ReplyID               string `json:"reply_id"`
		ReplyToReplyID        string `json:"reply_to_reply_id"`
		ShareInfo             struct {
			ACL struct {
				Code  int    `json:"code"`
				Extra string `json:"extra"`
			} `json:"acl"`
			Desc  string `json:"desc"`
			Title string `json:"title"`
			URL   string `json:"url"`
		} `json:"share_info"`
		SortExtraScore struct {
			ReplyScore    float64 `json:"reply_score"`
			ShowMoreScore float64 `json:"show_more_score"`
		} `json:"sort_extra_score"`
		SortTags      string `json:"sort_tags"`
		Status        int    `json:"status"`
		StickPosition int    `json:"stick_position"`
		Text          string `json:"text"`
		TextExtra     []any  `json:"text_extra"`
		TransBtnStyle int    `json:"trans_btn_style"`
		User          struct {
			AccountLabels           any `json:"account_labels"`
			AdCoverURL              any `json:"ad_cover_url"`
			AdvanceFeatureItemOrder any `json:"advance_feature_item_order"`
			AdvancedFeatureInfo     any `json:"advanced_feature_info"`
			AvatarThumb             struct {
				URI       string   `json:"uri"`
				URLList   []string `json:"url_list"`
				URLPrefix any      `json:"url_prefix"`
			} `json:"avatar_thumb"`
			BoldFields                 any    `json:"bold_fields"`
			CanMessageFollowStatusList any    `json:"can_message_follow_status_list"`
			CanSetGeofencing           any    `json:"can_set_geofencing"`
			ChaList                    any    `json:"cha_list"`
			CoverURL                   any    `json:"cover_url"`
			CustomVerify               string `json:"custom_verify"`
			EnterpriseVerifyReason     string `json:"enterprise_verify_reason"`
			Events                     any    `json:"events"`
			FollowersDetail            any    `json:"followers_detail"`
			Geofencing                 any    `json:"geofencing"`
			HomepageBottomToast        any    `json:"homepage_bottom_toast"`
			ItemList                   any    `json:"item_list"`
			MutualRelationAvatars      any    `json:"mutual_relation_avatars"`
			NeedPoints                 any    `json:"need_points"`
			Nickname                   string `json:"nickname"`
			PlatformSyncInfo           any    `json:"platform_sync_info"`
			PredictedAgeGroup          string `json:"predicted_age_group"`
			RelativeUsers              any    `json:"relative_users"`
			SearchHighlight            any    `json:"search_highlight"`
			SecUid                     string `json:"sec_uid"`
			ShieldEditFieldInfo        any    `json:"shield_edit_field_info"`
			TypeLabel                  any    `json:"type_label"`
			UID                        string `json:"uid"`
			UniqueID                   string `json:"unique_id"`
			UserProfileGuide           any    `json:"user_profile_guide"`
			UserTags                   any    `json:"user_tags"`
			WhiteCoverURL              any    `json:"white_cover_url"`
		} `json:"user"`
		UserBuried bool `json:"user_buried"`
		UserDigged int  `json:"user_digged"`
	} `json:"comments"`
	Cursor  int `json:"cursor"`
	HasMore int `json:"has_more"`
	Total   int `json:"total"`
}

type TikTokComment struct {
	ID            string `json:"id"`
	Text          string `json:"text"`
	CreateTime    int    `json:"create_time"`
	CreateTimeStr string `json:"create_time_str"`
	DiggCount     int    `json:"likes"`
	ReplyCount    int    `json:"replies"`
	Author        struct {
		ID       string `json:"id"`
		Nickname string `json:"name"`
		Avatar   string `json:"avatar"`
	} `json:"author"`
}

type TikTokConfig struct {
	CommentsMsToken string
	CommentsXBogus  string
	CommentsXGnarly string
	CommentsCookies string
	Username        string
}

func getRandomDelay() time.Duration {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomMs := r.Intn(501) + 1000
	return time.Duration(randomMs) * time.Millisecond
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

func main() {
	fmt.Println("📱 TikTok Comment Scraper (Free Limited Version)")
	fmt.Println("=================================================")
	fmt.Println("This tool extracts main comments from TikTok videos.")
	fmt.Println("✅ Unlimited main comments extraction - no limits!")
	fmt.Println()
	fmt.Println("💎 Need unlimited comments with replies?")
	fmt.Println("📧 Email: haronkibetrutoh@gmail.com")
	fmt.Println("📱 WhatsApp: +254718448461")
	fmt.Println()

	var url string
	fmt.Print("🔗 Enter TikTok URL: ")
	fmt.Scanln(&url)

	if url == "" {
		fmt.Printf("❌ Error: No URL provided\n")
		fmt.Println("💡 Example URLs:")
		fmt.Println("   - https://www.tiktok.com/@username/video/1234567890123456789")
		fmt.Println("   - https://vm.tiktok.com/XXXXXX/")
		os.Exit(1)
	}

	startTime := time.Now()
	fmt.Printf("🔍 Analyzing TikTok URL: %s\n", url)
	fmt.Printf("⏱️ Extraction started at: %s\n", startTime.Format("15:04:05"))

	finalVideoID, err := extractVideoID(url)
	if err != nil {
		fmt.Printf("❌ Error extracting video ID from URL: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully extracted video ID: %s\n", finalVideoID)

	config := TikTokConfig{
		CommentsMsToken: "Q2_5J14X-7tmkMN0K14oh1vzv_tABpJjL_HaJId04kRpfmrQMqo5bREb2RnD6XRTitJlgw-O9y8KD6OrLjVLkhlBXz8XxnfWRMDCGYhxzplwpv6cmGm6TuU1lczYHCJOnsv-MZvi_AAEpv4=",
		CommentsXBogus:  "DFSzswVYrbTANJH/CrLF/gRDAYuO",
		CommentsXGnarly: "M8IQ0jUNugJRfqLpi4u5jcNCN/5iEkiH3RtHs0Z6fHiMQo2Wra3J7RbYA0lfGZup98cJCgmuiTCSucPrPUkLTfMjtJXNe3NAY3u9Wy9S-khy3288b0m64IiNf/k4-MqSSAoshfYqFWLcsTvDFyrbQ/xuP9Ve651TbOD-D38iBQsqxckkOL9NklnDf/D/yGhRqB1r69iQ-3xfRdoxvnlhX7bzFCfXcg7awvYLRA5D8sJoKk-cM57zppufJiPsTY6ADYJPSVk9eViVU519kUGa6AiYZK-BxLr4qmQrXlFXuWpo",
		CommentsCookies: "delay_guest_mode_vid=5; tiktok_webapp_theme=light; passport_csrf_token=a4d4a1fdfc11e9161e4bf7ff312c1cb3; passport_csrf_token_default=a4d4a1fdfc11e9161e4bf7ff312c1cb3; living_user_id=383622962067; passport_auth_status=833fc28b9c072bced652b09badcaf621%2C; passport_auth_status_ss=833fc28b9c072bced652b09badcaf621%2C; last_login_method=google; store-country-sign=MEIEDB9MjCbM2jpZsWk1vQQgfY8YF6nkkuwPQwdR8rERB-3s2JL6t5s_RDxwjFoBecwEEEJv50bhS8Q3E7u5Sf5RGYQ; tiktok_webapp_theme_source=system; odin_tt=c87d5f0f43a404da5485bd59496710a53ba522110e7e4f710b307862c4f872cec2a8aec9d5ca4a63dffbc0c5b9b3bb4062d2474d13bc6af08a8aff7d40db37935fd6e2eaec9db3b38f57fd574f867e00; tt_csrf_token=hpJ1LswI-X0yy3MAsSJwVJJC1aTiAhiVTn98; ttwid=1%7CAo18CWxPgX0G9zqyYvlpBXr_zxFs7kLC-DgZ272CiXg%7C1753969719%7Cfdbee65e607d9bdc48c74e59d1357116a47d0696cd9c3b755dc731107004195d; tt_chain_token=MB1Ywgn8JFwJse0Q+VSZIQ==; s_v_web_id=verify_mdu49mbb_O9NppxWT_J8Ny_4Tk3_8Zxc_52uNHVK770Dp; msToken=uBMZpyz9HYlSiflk8mnYb_JZOyxAajndJ-HA0C3GzyDKzBjjYvToReuDTeIxHA_AoomeNSQ4nH_swgJ1L-0kdBTLLHKVCqK9sq5dXAQ_mOPHgGtQrt4vBwaAnSnog0VD_3kjYxZS7b-VCqI=; perf_feed_cache={%22expireTimestamp%22:1754301600000%2C%22itemIds%22:[%227526066216314752262%22%2C%227515114470423743750%22%2C%227526535589798186246%22]}; msToken=YLGRF5rM8DhPePlWeWfcqGuvEGMrz869raHX9CxuClyB6PlPbZlXiDUTZZnS6bACxyJ6ceesV7kgQ3GNAdwnQf0ZxkUwBZuunthXWNHQK-YT0Sv-SL_sdMjMsMyAR3cJ5w6raTZGCbyniEY=",
		Username:        "mcexodus",
	}

	fmt.Printf("📥 Extracting main comments for video %s...\n", finalVideoID)
	fmt.Println("✅ Unlimited comments extraction - no limits!")

	allComments := fetchAllComments(finalVideoID, config)

	if len(allComments) == 0 {
		fmt.Printf("❌ No comments were retrieved. Check debug files for more information.\n")
		os.Exit(1)
	}

	fmt.Printf("📊 Total main comments retrieved: %d\n", len(allComments))

	fmt.Printf("📊 Exporting comments to Excel...\n")
	excelPath, err := exportTikTokCommentsToExcel(allComments, url)
	if err != nil {
		fmt.Printf("❌ Error exporting to Excel: %v\n", err)
		os.Exit(1)
	}

	endTime := time.Now()
	actualDuration := endTime.Sub(startTime)
	minutes := int(actualDuration.Minutes())
	seconds := int(actualDuration.Seconds()) % 60

	fmt.Printf("⏱️ Extraction completed at: %s\n", endTime.Format("15:04:05"))
	fmt.Printf("🎯 Actual processing time: %d minutes %d seconds (%.1f seconds total)\n",
		minutes, seconds, actualDuration.Seconds())

	fmt.Printf("✅ Exported %d comments to Excel file: %s\n", len(allComments), excelPath)

	fmt.Println()
	fmt.Println("🎉 Extraction completed successfully!")
	fmt.Printf("📊 Summary: %d comments extracted and exported\n", len(allComments))
	fmt.Printf("📁 Excel file saved to: %s\n", excelPath)
	fmt.Println("💡 You can now open the Excel file to view all comments")
	fmt.Println("💎 Need unlimited comments with replies?")
	fmt.Println("📧 Email: haronkibetrutoh@gmail.com")
	fmt.Println("📱 WhatsApp: +254718448461")
	fmt.Println()
}

func fetchAllComments(videoID string, config TikTokConfig) []TikTokComment {
	var allComments []TikTokComment
	cursor := 0
	retryCount := 0
	maxRetries := 10

	fmt.Println("📥 Fetching comments...")

	for {
		// Removed comment limit check - now extracts unlimited comments

		fmt.Printf("📄 Fetching comments page with cursor %d...\n", cursor)
		response := fetchComments(videoID, cursor, &config)

		if response.StatusCode == 0 && len(response.Comments) == 0 {
			retryCount++
			if retryCount <= maxRetries {
				fmt.Printf("⚠️ Received empty response from API. Retry %d/%d after longer delay...\n", retryCount, maxRetries)

				time.Sleep(getRandomDelay())
				continue
			} else {
				fmt.Printf("🛑 Max retries reached. Stopping pagination.\n")
				break
			}
		} else {
			retryCount = 0
		}

		for _, c := range response.Comments {
			comment := TikTokComment{
				ID:            c.ID,
				Text:          c.Text,
				CreateTime:    c.CreateTime,
				CreateTimeStr: formatUnixTimestamp(c.CreateTime),
				DiggCount:     c.DiggCount,
				ReplyCount:    c.ReplyCommentTotal,
			}
			comment.Author.ID = c.User.UID
			comment.Author.Nickname = c.User.Nickname
			if len(c.User.AvatarThumb.URLList) > 0 {
				comment.Author.Avatar = c.User.AvatarThumb.URLList[0]
			}
			allComments = append(allComments, comment)

			// Removed comment limit check - now extracts unlimited comments
		}

		if response.HasMore != 1 || len(response.Comments) == 0 {
			fmt.Printf("🏁 No more comments available (HasMore=%d)\n", response.HasMore)
			break
		}

		if response.Cursor > 0 {
			cursor = response.Cursor
		} else {
			cursor += 20
		}

		fmt.Printf("✅ Retrieved %d comments so far...\n", len(allComments))

		time.Sleep(getRandomDelay())
	}

	fmt.Printf("📊 Total main comments retrieved: %d\n", len(allComments))
	return allComments
}

func formatUnixTimestamp(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}

func fetchComments(videoID string, cursor int, config *TikTokConfig) CommentResponse {
	url := fmt.Sprintf("https://www.tiktok.com/api/comment/list/?WebIdLastTime=1753969719&aid=1988&app_language=en&app_name=tiktok_web&aweme_id=%s&browser_language=en-US&browser_name=Mozilla&browser_online=true&browser_platform=Linux%%20x86_64&browser_version=5.0%%20%%28X11%%3B%%20Linux%%20x86_64%%29%%20AppleWebKit%%2F537.36%%20%%28KHTML%%2C%%20like%%20Gecko%%29%%20Chrome%%2F138.0.0.0%%20Safari%%2F537.36&channel=tiktok_web&cookie_enabled=true&count=20&cursor=%d&data_collection_enabled=true&device_id=7533242546057070085&device_platform=web_pc&focus_state=true&from_page=video&history_len=13&is_fullscreen=false&is_page_visible=true&odinId=7533581167397438520&os=linux&priority_region=&referer=&region=KE&screen_height=768&screen_width=1366&tz_name=Africa%%2FNairobi&user_is_login=false&verifyFp=verify_mdu49mbb_O9NppxWT_J8Ny_4Tk3_8Zxc_52uNHVK770Dp&webcast_language=en&msToken=%s&X-Bogus=%s&X-Gnarly=%s",
		videoID,
		cursor,
		config.CommentsMsToken,
		config.CommentsXBogus,
		config.CommentsXGnarly)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("❌ Error creating request: %v\n", err)
		return CommentResponse{}
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", fmt.Sprintf("https://www.tiktok.com/@%s/video/%s?lang=en", config.Username, videoID))
	req.Header.Set("sec-ch-ua", "\"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"138\", \"Google Chrome\";v=\"138\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"Linux\"")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", config.CommentsCookies)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Error making request: %v\n", err)
		return CommentResponse{}
	}
	defer resp.Body.Close()

	var body []byte

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Printf("❌ Error creating gzip reader: %v\n", err)
			return CommentResponse{}
		}
		defer gzReader.Close()
		body, err = io.ReadAll(gzReader)
	} else {
		body, err = io.ReadAll(resp.Body)
	}

	if err != nil {
		fmt.Printf("❌ Error reading response: %v\n", err)
		return CommentResponse{}
	}

	fmt.Printf("📡 Response Status: %s\n", resp.Status)
	fmt.Printf("📄 Response Body Length: %d bytes\n", len(body))
	if len(body) < 200 {
		fmt.Printf("📝 Response Body: %s\n", string(body))
	} else {
		fmt.Printf("📝 Response Body Sample: %s...\n", string(body[:200]))
	}

	if newMsToken := resp.Header.Get("X-Ms-Token"); newMsToken != "" {
		fmt.Printf("🔄 Updating msToken from response: %s...\n", newMsToken[:20])
		config.CommentsMsToken = newMsToken
	}

	if bdturingHeader := resp.Header.Get("Bdturing-Verify"); bdturingHeader != "" {
		fmt.Printf("⚠️ Warning: TikTok is requesting verification. This may affect further requests.\n")
		fmt.Printf("🔐 Verification details: %s\n", bdturingHeader)
	}

	var response CommentResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("❌ Error parsing JSON: %v\n", err)
		return CommentResponse{}
	}

	if response.StatusCode != 0 {
		fmt.Printf("⚠️ API returned status code: %d - %s\n", response.StatusCode, response.StatusMsg)

		if response.StatusCode == 401 || response.StatusCode == 403 {
			fmt.Printf("🔐 Token may have expired. Consider refreshing tokens.\n")
		}
	}

	return response
}

func extractVideoID(url string) (string, error) {
	if strings.Contains(url, "tiktok.com") {
		if strings.Contains(url, "/video/") {
			parts := strings.Split(url, "/video/")
			if len(parts) > 1 {
				videoID := strings.Split(parts[1], "?")[0]
				videoID = strings.Split(videoID, "#")[0]
				return videoID, nil
			}
		}
	}

	if isNumeric(url) {
		return url, nil
	}

	return "", fmt.Errorf("could not extract video ID from URL: %s", url)
}

func exportTikTokCommentsToExcel(comments []TikTokComment, sourceURL string) (string, error) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing Excel file:", err)
		}
	}()

	sheetName := "Comments"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", fmt.Errorf("error creating sheet: %w", err)
	}
	f.SetActiveSheet(index)

	f.DeleteSheet("Sheet1")

	headers := []string{
		"Comment ID", "Author Name", "Author ID", "Comment Text", "Created Time",
		"Likes Count", "Reply Count", "URL",
	}

	columnWidths := map[string]float64{
		"A": 20,
		"B": 25,
		"C": 20,
		"D": 60,
		"E": 20,
		"F": 12,
		"G": 12,
		"H": 40,
	}

	for col, width := range columnWidths {
		f.SetColWidth(sheetName, col, col, width)
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E0E0E0"},
			Pattern: 1,
		},
		Border: []excelize.Border{
			{Type: "bottom", Color: "#000000", Style: 1},
		},
	})
	if err != nil {
		fmt.Println("Warning: Error creating header style:", err)
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
		if headerStyle != 0 {
			f.SetCellStyle(sheetName, cell, cell, headerStyle)
		}
	}

	row := 2

	for _, comment := range comments {
		var commentURL string
		if sourceURL != "" {
			baseURL := strings.Split(sourceURL, "?")[0]
			commentURL = fmt.Sprintf("%s?commentId=%s", baseURL, comment.ID)
		}

		rowData := []any{
			comment.ID,
			comment.Author.Nickname,
			comment.Author.ID,
			comment.Text,
			comment.CreateTimeStr,
			comment.DiggCount,
			comment.ReplyCount,
			commentURL,
		}

		for i, value := range rowData {
			cell := fmt.Sprintf("%s%d", string(rune('A'+i)), row)
			f.SetCellValue(sheetName, cell, value)
		}

		row++
	}

	metaSheetName := "Metadata"
	_, err = f.NewSheet(metaSheetName)
	if err != nil {
		fmt.Println("Warning: Error creating metadata sheet:", err)
	} else {
		f.SetCellValue(metaSheetName, "A1", "Source URL")
		f.SetCellValue(metaSheetName, "B1", sourceURL)

		f.SetCellValue(metaSheetName, "A2", "Extraction Date")
		f.SetCellValue(metaSheetName, "B2", time.Now().Format("2006-01-02 15:04:05"))

		f.SetCellValue(metaSheetName, "A3", "Total Comments")
		f.SetCellValue(metaSheetName, "B3", len(comments))
	}

	exportDir := "exports"
	if _, err := os.Stat(exportDir); os.IsNotExist(err) {
		if err := os.Mkdir(exportDir, 0755); err != nil {
			return "", fmt.Errorf("error creating exports directory: %w", err)
		}
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(exportDir, fmt.Sprintf("tiktok_comments_%s.xlsx", timestamp))

	if err := f.SaveAs(filename); err != nil {
		return "", fmt.Errorf("error saving Excel file: %w", err)
	}

	fmt.Printf("✅ Exported %d comments to Excel file: %s\n", len(comments), filename)
	return filename, nil
}
