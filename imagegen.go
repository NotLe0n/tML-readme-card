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

func generateImage(steamId string) ([]byte, error) {
	getJson("https://tmlapis.repl.co/author_api/"+steamId, &author)
	return run(steamId)
}

var imageWidth float64
var imageHeight float64

const margin float64 = 20.0
const padding float64 = 5.0

func run(steamId string) ([]byte, error) {
	if steamId == "" {
		return nil, errors.New("please enter a valid steamid64")
	}

	imageWidth = 878.0
	imageHeight = (35.0+padding)*float64(len(author.Mods)) + (35 * 2) + margin*2 + 10
	dc := gg.NewContext(int(imageWidth), int(imageHeight))

	// Draw light gray rounded rectangle
	dc.SetColor(color.RGBA{35, 39, 42, 255})
	dc.DrawRoundedRectangle(0, 0, imageWidth, imageHeight, 15)
	dc.Fill()

	// Draw dark gray rectangle and leave 20px border
	x := margin
	y := margin
	w := float64(imageWidth) - (2.0 * margin)
	h := float64(imageHeight) - (2.0 * margin)
	dc.SetColor(color.RGBA{25, 28, 30, 255})
	dc.DrawRectangle(x, y, w, h)
	dc.Fill()

	// Load font
	fontPath := filepath.Join("fonts", "Andy Bold.ttf")
	dc.LoadFontFace(fontPath, 35)

	// Draw Text
	userNameWidth, _ := dc.MeasureString(author.SteamName + "'s Stats")
	DrawText(dc, author.SteamName+"'s Stats", (imageWidth-userNameWidth)/2, margin*2+10, 35, color.White)

	if len(author.Mods) == 0 {
		DrawText(dc, "No mods found", 30, margin*4+10, 35, color.White)
	} else {
		DrawText(dc, "Rank", 30, margin*4+10, 35, color.White)
		DrawText(dc, "Display Name", 120, margin*4+10, 35, color.White)
		DrawText(dc, "Downloads", imageWidth-190, margin*4+10, 35, color.White)

		dc.SetLineWidth(2)
		dc.DrawLine(30, margin*4+15, imageWidth-30, margin*4+15)
		dc.Stroke()

		for i := 0; i < len(author.Mods); i++ {
			_, nameTextHeight := dc.MeasureString(author.Mods[i].DisplayName)
			dowloadsTextWidth, _ := dc.MeasureString(strconv.Itoa(author.Mods[i].DownloadsTotal))

			DrawText(dc, strconv.Itoa(author.Mods[i].RankTotal), 30, (nameTextHeight+padding)*float64(i)+(nameTextHeight*2)+margin*4+10, 35, color.White)

			// NEW: parsing chat tags using regexp
			displayNameColor, displayName := ParseChatTags(author.Mods[i].DisplayName)
			DrawText(dc, displayName, 120, (nameTextHeight+padding)*float64(i)+(nameTextHeight*2)+margin*4+10, 35, displayNameColor)

			DrawText(dc, strconv.Itoa(author.Mods[i].DownloadsTotal), imageWidth-dowloadsTextWidth-50, (nameTextHeight+padding)*float64(i)+(nameTextHeight*2)+margin*4+10, 35, color.White)
		}
	}
	DrawText(dc, time.Now().Format("2006-01-02 15:04:05"), imageWidth-160, imageHeight-20, 15, color.White)

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

func ParseChatTags(str string) (textColor color.Color, text string) {
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

	return color.White, str
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
