package widgets

import (
	"bytes"
	"errors"
	"html"
	"sort"

	"github.com/spf13/viper"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type author struct {
	Total      uint32
	Mods       []mod
	Steam_name string
}

func GenerateAuthorWidget(steamId string, config ImgConfig) ([]byte, error) {
	if steamId == "" {
		return nil, errors.New("please enter a valid steamid64")
	}

	var author author
	if err := getJson(viper.GetString("api")+"/"+config.Version+"/author/"+steamId, &author); err != nil {
		return nil, err
	}

	return drawAuthorWidget(author, config)
}

func drawAuthorWidget(author author, config ImgConfig) ([]byte, error) {
	const padding float64 = 15.0
	prt := message.NewPrinter(language.AmericanEnglish)

	bw := float64(config.BorderWidth) // stands for border width
	imageWidth := 878.0
	modListHeight := (30 + padding) * max(float64(len(author.Mods)), 2)
	imageHeight := 26 + 40 + bw*2 + modListHeight + 10
	dc := drawCard(config, imageWidth, imageHeight)

	fontSize, fontErr := loadFont(dc, config)
	if fontErr != nil {
		return nil, fontErr
	}

	// Draw text
	userNameWidth, userNameHeight := dc.MeasureString(author.Steam_name + "'s Stats")
	drawBorderText(dc, author.Steam_name+"'s Stats", (imageWidth-userNameWidth)/2, bw+35, config.TextColor)

	headerY := userNameHeight + 40 + bw*2
	if len(author.Mods) == 0 {
		drawBorderText(dc, "No mods found", 30, headerY+fontSize/2, config.TextColor)
	} else {
		// Draw header
		startX := 30.0
		if config.Version == "1.3" {
			drawBorderText(dc, "Rank", startX, headerY, config.TextColor)
		} else {
			startX = -30.0
		}

		drawBorderText(dc, "Display Name", startX+90, headerY, config.TextColor)
		drawBorderText(dc, "Downloads", imageWidth-190, headerY, config.TextColor)

		// Draw line
		dc.SetLineWidth(2)
		dc.DrawLine(30, headerY+5, imageWidth-30, headerY+5)
		dc.Stroke()

		sort.Slice(author.Mods, func(i, j int) bool {
			return author.Mods[i].Downloads_total > author.Mods[j].Downloads_total
		})

		for i := 0; i < len(author.Mods); i++ {
			downloadsStr := prt.Sprintf("%d", author.Mods[i].Downloads_total)
			nameTextWidth, nameTextHeight := dc.MeasureString(author.Mods[i].Display_name)
			downloadsTextWidth, _ := dc.MeasureString(downloadsStr)

			modY := (nameTextHeight+padding)*float64(i) + (nameTextHeight * 2)

			if config.Version == "1.3" {
				// Draw Rank
				drawBorderText(dc, prt.Sprint(author.Mods[i].Rank), startX, modY+headerY, config.TextColor)
			}

			// Draw Display Name
			scale := 610.0 / nameTextWidth
			if scale < 1 {
				_, _ = loadFontSized(dc, config, 35*scale) // resize font
			}

			displayNameSnippets := parseChatTags(html.UnescapeString(author.Mods[i].Display_name), config.TextColor)
			drawSnippets(dc, displayNameSnippets, startX+90, headerY+modY)

			if scale < 1 {
				_, _ = loadFont(dc, config) // reset font size
			}

			// Draw downloads
			drawBorderText(dc, downloadsStr, imageWidth-downloadsTextWidth-50, modY+headerY, config.TextColor)
		}
	}

	var b bytes.Buffer
	err := dc.EncodePNG(&b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
