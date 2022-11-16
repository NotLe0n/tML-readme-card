package widgets

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"path/filepath"
	"time"
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
	Downloads_yesterday int
	Icon                string
	Last_updated        string
	Version             string
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

func parseChatTags(str string, defaultColor color.Color) []textSnippet {
	snippets := make([]textSnippet, 0)

	for index := 0; index < len(str); index++ {
		if str[index] == '[' && str[index+1] == 'c' && str[index+2] == '/' {
			index += 3 // skip '[c/'

			// parse color code
			colorCode := ""
			// check for HEX number
			for (str[index] >= 'A' && str[index] <= 'F') || (str[index] >= '0' && str[index] <= '9') {
				colorCode += string(str[index])
				index++
			}
			index++ // skip ':'

			// read text until ']'
			text := ""
			for str[index] != ']' {
				text += string(str[index])
				index++
			}

			b, err := hex.DecodeString(colorCode)
			if err != nil {
				log.Println(err) //this should never happen, so we don't need to 'throw' the error, but if it happens we know where
			}
			col := color.RGBA{R: b[0], G: b[1], B: b[2], A: 255}

			snippets = append(snippets, textSnippet{
				text:  text,
				color: col,
			})
		} else {
			text := ""
			for index < len(str) {
				text += string(str[index])
				if index+1 < len(str) && str[index+1] == '[' {
					break
				}

				index++
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
