package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"image/color"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/fogleman/gg"
)

type Mod struct {
	DisplayName        string
	RankTotal          int
	DownloadsTotal     int
	DownloadsYesterday int
}

type Author struct {
	SteamName string
	Mods      []Mod
}

var myClient = &http.Client{Timeout: 10 * time.Second}

var author Author

func generateImage(steamId string, config ImgConfig) ([]byte, error) {
	getJson("https://tmlapis.repl.co/author_api/"+steamId, &author)
	return run(steamId, config)
}

var imageWidth float64
var imageHeight float64

const padding float64 = 5.0

func run(steamId string, config ImgConfig) ([]byte, error) {
	if steamId == "" {
		return nil, errors.New("please enter a valid steamid64")
	}

	bw := float64(config.borderWidth) // stands for border width
	imageWidth = 878.0
	imageHeight = (35.0+padding)*float64(len(author.Mods)) + (35 * 2) + bw*2 + 10
	dc := gg.NewContext(int(imageWidth), int(imageHeight))

	// Draw light gray rounded rectangle
	dc.SetColor(config.borderColor)
	dc.DrawRoundedRectangle(0, 0, imageWidth, imageHeight, float64(config.cornerRadius))
	dc.Fill()

	// Draw dark gray rectangle and leave 20px border
	w := float64(imageWidth) - (2.0 * bw)
	h := float64(imageHeight) - (2.0 * bw)
	dc.SetColor(config.bgColor)
	dc.DrawRoundedRectangle(bw, bw, w, h, float64(config.cornerRadius))
	dc.Fill()

	// Load font
	fontPath := filepath.Join("fonts", "Andy Bold.ttf")
	fontSize := 35.0
	dc.LoadFontFace(fontPath, fontSize)

	// Draw Text
	userNameWidth, userNameHeight := dc.MeasureString(author.SteamName + "'s Stats")
	DrawText(dc, author.SteamName+"'s Stats", (imageWidth-userNameWidth)/2, bw+35, fontSize, config.textColor)

	headerY := userNameHeight + 20 + bw*2
	if len(author.Mods) == 0 {
		DrawText(dc, "No mods found", 30, headerY, fontSize, config.textColor)
	} else {
		DrawText(dc, "Rank", 30, headerY, fontSize, config.textColor)
		DrawText(dc, "Display Name", 120, headerY, fontSize, config.textColor)
		DrawText(dc, "Downloads", imageWidth-190, headerY, fontSize, config.textColor)

		dc.SetLineWidth(2)
		dc.DrawLine(30, headerY, imageWidth-30, headerY)
		dc.Stroke()

		for i := 0; i < len(author.Mods); i++ {
			_, nameTextHeight := dc.MeasureString(author.Mods[i].DisplayName)
			dowloadsTextWidth, _ := dc.MeasureString(strconv.Itoa(author.Mods[i].DownloadsTotal))

			DrawText(dc, strconv.Itoa(author.Mods[i].RankTotal), 30, (nameTextHeight+padding)*float64(i)+(nameTextHeight*2)+headerY, fontSize, config.textColor)

			// NEW: parsing chat tags using regexp
			displayNameColor, displayName := ParseChatTags(author.Mods[i].DisplayName, config.textColor)
			DrawText(dc, displayName, 120, (nameTextHeight+padding)*float64(i)+(nameTextHeight*2)+headerY, fontSize, displayNameColor)

			DrawText(dc, strconv.Itoa(author.Mods[i].DownloadsTotal), imageWidth-dowloadsTextWidth-50, (nameTextHeight+padding)*float64(i)+(nameTextHeight*2)+headerY, fontSize, config.textColor)
		}
	}
	DrawText(dc, time.Now().Format("2006-01-02 15:04:05"), imageWidth-160, imageHeight-20, 15, config.textColor)

	var b bytes.Buffer
	err := dc.EncodePNG(&b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func DrawText(dc *gg.Context, s string, x float64, y float64, pnt float64, col color.Color) {
	// Load font
	fontPath := filepath.Join("fonts", "Andy Bold.ttf")
	dc.LoadFontFace(fontPath, pnt)

	dc.SetColor(col)
	textWidth, textHeight := dc.MeasureString(s)
	x = ClampFloat(x, 0, imageWidth-textWidth)
	y = ClampFloat(y, textHeight, imageHeight-textHeight)
	dc.DrawString(s, x, y)
}

func ParseChatTags(str string, defaultColor color.Color) (textColor color.Color, text string) {
	var compRegEx = regexp.MustCompile(`\[c\/(?P<col>\w+):(?P<text>[\s\S]+?)\]`)

	if compRegEx.MatchString(str) {
		match := compRegEx.FindStringSubmatch(str)

		paramsMap := make(map[string]string)
		for i, name := range compRegEx.SubexpNames() {
			if i > 0 && i <= len(match) {
				paramsMap[name] = match[i]
			}
		}

		b, err := hex.DecodeString(paramsMap["col"])
		if err != nil {
			log.Println(err) //this should never happen so we don't need to 'throw' the error, but if it happens we know where
		}
		col := color.RGBA{b[0], b[1], b[2], 255}

		return col, paramsMap["text"]
	}

	return defaultColor, str
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

	return json.NewDecoder(r.Body).Decode(&target)
}
