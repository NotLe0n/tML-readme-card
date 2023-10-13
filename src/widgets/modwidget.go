package widgets

import (
	"bytes"
	"errors"
	"image/color"
	"strconv"
	"time"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

func GenerateModWidget(modname string, config ImgConfig) ([]byte, error) {
	if modname == "" {
		return nil, errors.New("please enter a valid modname")
	}

	var modStruct mod
	if err := getJson("https://tmlapis.le0n.dev/"+config.Version+"/mod/"+modname, &modStruct); err != nil {
		return nil, err
	}

	return drawModWidget(modStruct, config)
}

func drawModWidget(mod mod, config ImgConfig) ([]byte, error) {
	imageWidth := 878.0
	imageHeight := 240.0

	dc := gg.NewContext(int(imageWidth), int(imageHeight)) // draw context

	// Draw card background
	dc.SetColor(config.BorderColor)
	dc.DrawRoundedRectangle(0, 0, imageWidth, imageHeight, float64(config.CornerRadius))
	dc.Fill()

	// draw mod icon
	iconDim, err := drawIcon(dc, config, mod)
	if err != nil {
		return nil, err
	}

	// draw info text
	if err := drawModInfoText(iconDim, imageWidth, dc, config, mod); err != nil {
		return nil, err
	}

	// return generated image
	var b bytes.Buffer
	err = dc.EncodePNG(&b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// draws the mod icon and returns the dimensions of the drawn area, or an error if it fails to load the image
func drawIcon(dc *gg.Context, config ImgConfig, mod mod) (float64, error) {
	cardPadding := 20.0
	iconSize := 80.0 * 2
	iconPadding := 40.0

	// draw icon background
	dc.SetColor(config.BgColor)
	dc.DrawRoundedRectangle(cardPadding, cardPadding, iconSize+iconPadding, iconSize+iconPadding, float64(config.CornerRadius)-5)
	dc.Fill()

	// draw mod icon
	var iconURL string
	if config.Version == "1.3" {
		iconURL = mod.Icon
	} else {
		iconURL = mod.Workshop_icon_url
	}

	img, imgErr := getImage(iconURL)
	if imgErr != nil {
		return 0, imgErr
	}

	img = resize.Resize(uint(iconSize), uint(iconSize), img, resize.NearestNeighbor)
	dc.DrawImage(img, int(iconPadding), int(iconPadding))

	return cardPadding + iconPadding + iconSize, nil
}

func drawModInfoText(iconDim, imageWidth float64, dc *gg.Context, config ImgConfig, mod mod) error {
	// load header font
	fontHeight, fontErr := loadFontSized(dc, config, 40)
	if fontErr != nil {
		return fontErr
	}

	// draw text
	yPos := 60.0
	displayNameSnippets := parseChatTags(mod.Display_name, color.White)

	// get the combined string
	fullStr := ""
	for _, snippet := range displayNameSnippets {
		fullStr += snippet.text
	}

	// draw all displayNameSnippets centered
	drawSnippets(dc, displayNameSnippets, func(snippet textSnippet, prevTextWidth float64) {
		// calculate the centered position
		textStart := calculateCenteredInfoTextStart(fullStr, iconDim, imageWidth, dc)
		dc.SetColor(snippet.color) // set color to snippet color
		dc.DrawString(snippet.text, textStart+prevTextWidth, yPos)
	})

	// load dataFont
	fontHeight, fontErr = loadFontSized(dc, config, 32)
	if fontErr != nil {
		return fontErr
	}

	yPos += fontHeight + 15
	drawTextCentered(dc, strconv.Itoa(mod.Downloads_total)+" Downloads Total", 0, yPos, iconDim, imageWidth, color.White)
	yPos += fontHeight + 15

	if config.Version == "1.3" {
		drawTextCentered(dc, strconv.Itoa(mod.Downloads_yesterday)+" Downloads Yesterday", 0, yPos, iconDim, imageWidth, color.White)
	} else {
		drawTextCentered(dc, strconv.Itoa(int(mod.Favorited))+" Favorites", 0, yPos, iconDim, imageWidth, color.White)
	}

	yPos += fontHeight + 15

	var lastUpdateTime string
	var v string
	if config.Version == "1.3" {
		lastUpdateTime = mod.Last_updated
		v = mod.Version
	} else {
		lastUpdateTime = time.Unix(int64(mod.Time_updated), 0).Format(time.RFC822)
		v = "v" + mod.Versions[0].Mod_version
	}
	drawTextCentered(dc, "Last updated: "+lastUpdateTime+" ("+v+")", 0, yPos, iconDim, imageWidth, color.White)

	return nil
}

func calculateCenteredInfoTextStart(text string, iconDim, imageWidth float64, dc *gg.Context) float64 {
	textWidth, _ := dc.MeasureString(text)
	return (imageWidth - textWidth + iconDim) / 2
}
