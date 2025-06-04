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
	Api          string
}

type textSnippet struct {
	text  string
	color color.Color
}

func parseChatTags(str string, defaultColor color.Color) []textSnippet {
	snippets := make([]textSnippet, 0)

	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
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
	return loadFontSized(dc, config, 35)
}

func loadFontSized(dc *gg.Context, config ImgConfig, fontSize float64) (float64, error) {
	fontPath := ""
	// Load font
	switch config.Font {
	case "Andy":
		fontPath = filepath.Join("fonts", "Andy-Bold-RU-CN.ttf")
	case "Sans":
		fontPath = filepath.Join("fonts", "FreeSans.ttf")
	}

	fontErr := dc.LoadFontFace(fontPath, fontSize)

	return fontSize, fontErr
}

func drawCard(config ImgConfig, imageWidth, imageHeight float64) *gg.Context {
	dc := gg.NewContext(int(imageWidth), int(imageHeight)) // draw context

	// Draw border
	dc.SetColor(config.BorderColor)
	dc.DrawRoundedRectangle(0, 0, imageWidth, imageHeight, float64(config.CornerRadius))
	dc.Fill()

	// Draw background
	bw := float64(config.BorderWidth) // stands for border width
	w := imageWidth - (2.0 * bw)
	h := imageHeight - (2.0 * bw)
	dc.SetColor(config.BgColor)
	dc.DrawRoundedRectangle(bw, bw, w, h, float64(config.CornerRadius-config.BorderWidth))
	dc.Fill()

	return dc
}

func drawBorderText(dc *gg.Context, str string, x, y float64, col color.Color) {
	dc.SetColor(color.Black)
	const n = 4 // "stroke" size
	for dy := -n; dy <= n; dy++ {
		for dx := -n; dx <= n; dx++ {
			if dx*dx+dy*dy >= n*n {
				// give it rounded corners
				continue
			}
			dc.DrawString(str, x+float64(dx), y+float64(dy))
		}
	}
	dc.SetColor(col)
	dc.DrawString(str, x, y)
}

func drawTextCentered(dc *gg.Context, str string, x, y, imageWidth float64, color color.Color) {
	textWidth, _ := dc.MeasureString(str)
	drawBorderText(dc, str, (imageWidth-textWidth+x)/2, y, color)
}

func drawSnippets(dc *gg.Context, snippets []textSnippet, x, y float64) {
	for i, snippet := range snippets {
		prevTextWidth := 0.0
		for _, prevSnippet := range snippets[:i] {
			measuredWidth, _ := dc.MeasureString(prevSnippet.text)
			prevTextWidth += measuredWidth
		}

		drawBorderText(dc, snippet.text, x+prevTextWidth, y, snippet.color)
	}
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
