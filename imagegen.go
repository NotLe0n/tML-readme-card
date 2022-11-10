package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"image/color"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/fogleman/gg"
)

type Mod struct {
	Rank                int
	Display_name        string
	Downloads_total     int
	Downloads_yesterday int
}

type Author struct {
	Total      uint32
	Mods       []Mod
	Steam_name string
}

type TextSnippet struct {
	text  string
	color color.Color
}

var myClient = &http.Client{Timeout: 20 * time.Second}

var author Author

func generateImage(steamId string, config ImgConfig) ([]byte, error) {
	if steamId == "" {
		return nil, errors.New("please enter a valid steamid64")
	}

	if err := getJson("https://tmlapis.tomat.dev/"+config.version+"/author/"+steamId, &author); err != nil {
		return nil, err
	}

	return run(config)
}

var imageWidth float64
var imageHeight float64

const padding float64 = 15.0

func run(config ImgConfig) ([]byte, error) {

	bw := float64(config.borderWidth) // stands for border width
	imageWidth = 878.0
	modListHeight := (32 + padding) * math.Max(float64(len(author.Mods)), 2)
	imageHeight = 26.25 + 40 + bw*2 + modListHeight + 10
	dc := gg.NewContext(int(imageWidth), int(imageHeight)) // draw context

	// Draw border
	dc.SetColor(config.borderColor)
	dc.DrawRoundedRectangle(0, 0, imageWidth, imageHeight, float64(config.cornerRadius))
	dc.Fill()

	// Draw background
	w := imageWidth - (2.0 * bw)
	h := imageHeight - (2.0 * bw)
	dc.SetColor(config.bgColor)
	dc.DrawRoundedRectangle(bw, bw, w, h, float64(config.cornerRadius))
	dc.Fill()

	fontPath := ""
	// Load font
	switch config.font {
	case "Andy":
		fontPath = filepath.Join("fonts", "Andy Bold.ttf")
	case "Sans":
		fontPath = filepath.Join("fonts", "FreeSans.ttf")
	}
	fontSize := 35.0
	fontErr := dc.LoadFontFace(fontPath, fontSize)
	if fontErr != nil {
		return nil, fontErr
	}

	// Draw Text
	userNameWidth, userNameHeight := dc.MeasureString(author.Steam_name + "'s Stats")
	DrawText(dc, author.Steam_name+"'s Stats", (imageWidth-userNameWidth)/2, bw+35, fontSize, config.textColor)

	headerY := userNameHeight + 40 + bw*2
	if len(author.Mods) == 0 {
		DrawText(dc, "No mods found", 30, headerY+fontSize/2, fontSize, config.textColor)
	} else {
		// Draw header
		DrawText(dc, "Rank", 30, headerY, fontSize, config.textColor)
		DrawText(dc, "Display Name", 120, headerY, fontSize, config.textColor)
		DrawText(dc, "Downloads", imageWidth-190, headerY, fontSize, config.textColor)

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
			DrawText(dc, strconv.Itoa(author.Mods[i].Rank), 30, modY+headerY, fontSize, config.textColor)

			// Draw Display Name
			displayNameSnippets := ParseChatTags(html.UnescapeString(author.Mods[i].Display_name), config.textColor)
			for i, snippet := range displayNameSnippets {
				lastTextWidth := 0.0
				for _, prevSnippet := range displayNameSnippets[:i] {
					measuredWidth, _ := dc.MeasureString(prevSnippet.text)
					lastTextWidth += measuredWidth
				}

				DrawText(dc, snippet.text, 120+lastTextWidth, modY+headerY, fontSize, snippet.color)
			}

			// Draw downloads
			DrawText(dc, strconv.Itoa(author.Mods[i].Downloads_total), imageWidth-downloadsTextWidth-50, modY+headerY, fontSize, config.textColor)
		}
	}

	var b bytes.Buffer
	err := dc.EncodePNG(&b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func DrawText(dc *gg.Context, s string, x float64, y float64, pnt float64, col color.Color) {
	dc.SetColor(col)
	textWidth, textHeight := dc.MeasureString(s)
	x = ClampFloat(x, 0, imageWidth-textWidth)
	y = ClampFloat(y, textHeight, imageHeight-textHeight)
	dc.DrawString(s, x, y)
}

func ParseChatTags(str string, defaultColor color.Color) []TextSnippet {
	snippets := make([]TextSnippet, 0)

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

			snippets = append(snippets, TextSnippet{
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
			snippets = append(snippets, TextSnippet{
				text:  text,
				color: defaultColor,
			})
		}
	}

	return snippets
}

func ClampFloat(v float64, min float64, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("request returned with status code: %d", r.StatusCode)
	}

	return json.NewDecoder(r.Body).Decode(&target)
}
