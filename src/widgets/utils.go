package widgets

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/fogleman/gg"
)

type ImgConfig struct {
	TextColor    color.Color
	BgColor      color.Color
	BorderColor  color.Color
	BorderWidth  uint64
	CornerRadius uint64
	Version      string
	Font         string
}

type mod struct {
	Rank                int
	Display_name        string
	Downloads_total     int
	Downloads_yesterday int    // 1.3
	Favorited           uint   // 1.4
	Icon                string // 1.3
	Workshop_icon_url   string // 1.4
	Last_updated        string // 1.3
	Time_updated        uint   // 1.4
	Version             string
	Versions            []version // 1.4
}

type author struct {
	Total      uint32
	Mods       []mod
	Steam_name string
}

type textSnippet struct {
	text  string
	color color.Color
}

type version struct {
	Mod_version        string
	Tmodloader_version string
}

func parseChatTags(str string, defaultColor color.Color) []textSnippet {
	snippets := make([]textSnippet, 0)

	runes := []rune(str)
	for i := 0; i < len(runes); {
		if runes[i] == '[' && runes[i+1] == 'c' && runes[i+2] == '/' {
			i += 3 // skip '[c/'

			// parse color code
			colorCode := ""
			// check for HEX number
			for (runes[i] >= 'A' && runes[i] <= 'F') || (runes[i] >= '0' && runes[i] <= '9') {
				colorCode += string(runes[i])
				i++
			}
			i++ // skip ':'

			// read text until ']'
			text := ""
			for runes[i] != ']' {
				text += string(runes[i])
				i++
			}

			b, err := hex.DecodeString(colorCode)
			if err != nil {
				log.Println("Error while decoding hex string: " + err.Error())
			}
			col := color.RGBA{R: b[0], G: b[1], B: b[2], A: 255}

			snippets = append(snippets, textSnippet{
				text:  text,
				color: col,
			})
		} else {
			text := ""
			for i < len(runes) {
				text += string(runes[i])
				if i+1 < len(runes) && runes[i+1] == '[' {
					break
				}
				i++
			}
			snippets = append(snippets, textSnippet{
				text:  text,
				color: defaultColor,
			})
		}
	}

	return snippets
}

func loadFont(dc *gg.Context, config ImgConfig) (float64, error) {
	fontPath := ""
	// Load font
	switch config.Font {
	case "Andy":
		fontPath = filepath.Join("fonts", "Andy Bold.ttf")
	case "Sans":
		fontPath = filepath.Join("fonts", "FreeSans.ttf")
	}
	fontSize := 35.0
	fontErr := dc.LoadFontFace(fontPath, fontSize)

	return fontSize, fontErr
}

func loadFontSized(dc *gg.Context, config ImgConfig, fontSize float64) (float64, error) {
	fontPath := ""
	// Load font
	switch config.Font {
	case "Andy":
		fontPath = filepath.Join("fonts", "Andy Bold.ttf")
	case "Sans":
		fontPath = filepath.Join("fonts", "FreeSans.ttf")
	}

	fontErr := dc.LoadFontFace(fontPath, fontSize)

	return fontSize, fontErr
}

func drawText(dc *gg.Context, s string, x, y, imagewidth, imageheight float64, col color.Color) {
	dc.SetColor(col)
	textWidth, textHeight := dc.MeasureString(s)
	x = clampFloat(x, 0, imagewidth-textWidth)
	y = clampFloat(y, textHeight, imageheight-textHeight)
	dc.DrawString(s, x, y)
}

func drawTextCentered(dc *gg.Context, str string, xOffset, y, iconDim, imageWidth float64, color color.Color) {
	textStart := calculateCenteredInfoTextStart(str, iconDim, imageWidth, dc)
	dc.SetColor(color)
	dc.DrawString(str, textStart+xOffset, y)
}

func drawSnippets(dc *gg.Context, snippets []textSnippet, drawFunc func(snippet textSnippet, prevTextWidth float64)) {
	for i, snippet := range snippets {
		prevTextWidth := 0.0
		for _, prevSnippet := range snippets[:i] {
			measuredWidth, _ := dc.MeasureString(prevSnippet.text)
			prevTextWidth += measuredWidth
		}

		drawFunc(snippet, prevTextWidth)
	}
}

func clampFloat(v float64, min float64, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

var myClient = &http.Client{Timeout: 20 * time.Second}

func request(url string) (*http.Response, error) {
	r, err := myClient.Get(url)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("JSON request returned with status code: %d", r.StatusCode)
	}

	return r, nil
}

func getJson(url string, target interface{}) error {
	r, err := request(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(&target)
}

func getImage(url string) (image.Image, error) {
	r, err := request(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	img, err := png.Decode(r.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}
