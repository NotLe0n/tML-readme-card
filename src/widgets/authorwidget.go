package widgets

import (
	"bytes"
	"errors"
	"github.com/fogleman/gg"
	"html"
	"math"
	"sort"
	"strconv"
)

func GenerateAuthorWidget(steamId string, config ImgConfig) ([]byte, error) {
	if steamId == "" {
		return nil, errors.New("please enter a valid steamid64")
	}

	var author author
	if err := getJson("https://tmlapis.tomat.dev/"+config.Version+"/author/"+steamId, &author); err != nil {
		return nil, err
	}

	return drawAuthorWidget(author, config)
}

func drawAuthorWidget(author author, config ImgConfig) ([]byte, error) {
	const padding float64 = 15.0

	bw := float64(config.BorderWidth) // stands for border width
	imageWidth := 878.0
	modListHeight := (32 + padding) * math.Max(float64(len(author.Mods)), 2)
	imageHeight := 26.25 + 40 + bw*2 + modListHeight + 10
	dc := gg.NewContext(int(imageWidth), int(imageHeight)) // draw context

	// Draw border
	dc.SetColor(config.BorderColor)
	dc.DrawRoundedRectangle(0, 0, imageWidth, imageHeight, float64(config.CornerRadius))
	dc.Fill()

	// Draw background
	w := imageWidth - (2.0 * bw)
	h := imageHeight - (2.0 * bw)
	dc.SetColor(config.BgColor)
	dc.DrawRoundedRectangle(bw, bw, w, h, float64(config.CornerRadius))
	dc.Fill()

	fontSize, fontErr := loadFont(dc, config)
	if fontErr != nil {
		return nil, fontErr
	}

	// Draw text
	userNameWidth, userNameHeight := dc.MeasureString(author.Steam_name + "'s Stats")
	drawText(dc, author.Steam_name+"'s Stats", (imageWidth-userNameWidth)/2, bw+35, imageWidth, imageHeight, config.TextColor)

	headerY := userNameHeight + 40 + bw*2
	if len(author.Mods) == 0 {
		drawText(dc, "No mods found", 30, headerY+fontSize/2, imageWidth, imageHeight, config.TextColor)
	} else {
		// Draw header
		drawText(dc, "Rank", 30, headerY, imageWidth, imageHeight, config.TextColor)
		drawText(dc, "Display Name", 120, headerY, imageWidth, imageHeight, config.TextColor)
		drawText(dc, "Downloads", imageWidth-190, headerY, imageWidth, imageHeight, config.TextColor)

		// Draw line
		dc.SetLineWidth(2)
		dc.DrawLine(30, headerY+5, imageWidth-30, headerY+5)
		dc.Stroke()

		sort.Slice(author.Mods, func(i, j int) bool {
			return author.Mods[i].Downloads_total > author.Mods[j].Downloads_total
		})

		for i := 0; i < len(author.Mods); i++ {
			_, nameTextHeight := dc.MeasureString(author.Mods[i].Display_name)
			downloadsTextWidth, _ := dc.MeasureString(strconv.Itoa(author.Mods[i].Downloads_total))

			modY := (nameTextHeight+padding)*float64(i) + (nameTextHeight * 2)
			// Draw Rank
			drawText(dc, strconv.Itoa(author.Mods[i].Rank), 30, modY+headerY, imageWidth, imageHeight, config.TextColor)

			// Draw Display Name
			displayNameSnippets := parseChatTags(html.UnescapeString(author.Mods[i].Display_name), config.TextColor)
			drawSnippets(dc, displayNameSnippets, func(snippet textSnippet, prevTextWidth float64) {
				drawText(dc, snippet.text, 120+prevTextWidth, modY+headerY, imageWidth, imageHeight, snippet.color)
			})

			// Draw downloads
			drawText(dc, strconv.Itoa(author.Mods[i].Downloads_total), imageWidth-downloadsTextWidth-50, modY+headerY, imageWidth, imageHeight, config.TextColor)
		}
	}

	var b bytes.Buffer
	err := dc.EncodePNG(&b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
