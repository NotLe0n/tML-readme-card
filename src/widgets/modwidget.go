package widgets

import (
	"bytes"
	"errors"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image/color"
	"strconv"
)

func GenerateModWidget(modname string, config ImgConfig) ([]byte, error) {
	if modname == "" {
		return nil, errors.New("please enter a valid modname")
	}

	var modStruct mod
	if err := getJson("https://tmlapis.tomat.dev/"+config.Version+"/mod/"+modname, &modStruct); err != nil {
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
	img, imgErr := getImage(mod.Icon)
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
	drawTextCentered(dc, strconv.Itoa(mod.Downloads_yesterday)+" Downloads Yesterday", 0, yPos, iconDim, imageWidth, color.White)
	yPos += fontHeight + 15
	drawTextCentered(dc, "Last updated: "+mod.Last_updated+" ("+mod.Version+")", 0, yPos, iconDim, imageWidth, color.White)

	return nil
}

func calculateCenteredInfoTextStart(text string, iconDim, imageWidth float64, dc *gg.Context) float64 {
	textWidth, _ := dc.MeasureString(text)
	return (imageWidth - textWidth + iconDim) / 2
}
