# tML-readme-card
This program can generate widgets to display tModLoader mod or author information, for use in embeds, a github README or elsewhere.

A hosted version of this program can be found on:
Use `https://tml-card.le0n.dev/` 

Mod names with color chat tags will have their color changed and the tags removed.
tModLoader 1.3 and 1.4 have different mod browsers. Use the `v` parameter to specify the version (1.4 by default) 

## Parameters
| Parameter       | Value                         | Effect                                              |
|-----------------|-------------------------------|-----------------------------------------------------|
| `modname`       | string                        | Generates a mod widget using the mod's internal name|
| `steamid64`     | number                        | Generates a author widget                           |
| `text_color`    | Hex color value without the # | Changes the color of all text, except the mod names |
| `bg_color`      | Hex color value without the # | Changes the background color                        |
| `border_color`  | Hex color value without the # | Changes the border color                            |
| `border_width`  | number                        | Border width in pixels                              |
| `corner_radius` | number                        | corner radius, 0 for Rectangle                      |
| `v`             | "1.3" or "1.4"                | changes the tML version                             |
| `font`          | "Andy" or "Sans"              | changes the font                                    |

## Examples

### 1.4 Mod - default styling
* `https://tml-card.le0n.dev/?modname=NoFishingQuests`
![example-widget](https://tml-card.le0n.dev/?modname=NoFishingQuests)

### 1.3 Mod - default styling
* `https://tml-card.le0n.dev/?modname=CraftablePaint&v=1.3`
![example-widget](https://tml-card.le0n.dev/?modname=CraftablePaint&v=1.3)

### 1.3 Mod - custom styling
* `https://tml-card.le0n.dev/?modname=SummonHeart&border_width=1&corner_radius=60&border_color=FFFFFF&bg_color=0D1116&v=1.3`
![example-widget](https://tml-card.le0n.dev/?modname=SummonHeart&border_width=1&corner_radius=60&border_color=FFFFFF&bg_color=0D1116&v=1.3)

### 1.4 Mod - custom styling
* `https://tml-card.le0n.dev/?modname=BetterRussian&bg_color=F1F2F4&border_color=DCDFE4`
![example-widget](https://tml-card.le0n.dev/?modname=BetterRussian&bg_color=F1F2F4&border_color=DCDFE4)

### 1.4 Author - default styling
* `https://tml-card.le0n.dev/?steamid64=76561198278789341`
![example-widget](https://tml-card.le0n.dev/?steamid64=76561198278789341)

### 1.3 Author - custom styling
* `https://tml-card.le0n.dev/?steamid64=76561198278789341&border_width=1&corner_radius=60&border_color=FFFFFF&bg_color=0D1116&v=1.3`
![example-widget-parameters](https://tml-card.le0n.dev/?steamid64=76561198278789341&border_width=1&corner_radius=60&border_color=FFFFFF&bg_color=0D1116&v=1.3)

## Hosting Locally
1. Create config.json (for default values, write: `{}`)
2. Run the server using `go run ./src`

### Config
The default config is as follows:
```json
{
	"port": "8005",
	"useHTTPS": false,
	"certPath": "",
	"keyPath": ""
}
```

## License
All source code is licensed under GPL-2.0. 

The file `Andy-Bold-RU-CN.ttf` is a modification of the Andy typeface with cyrillic and chinese characters added.
* Andy Bold designed by Steve Matteson and owned by Monotype Imaging
* Cyrillic characters added by the [Translators Team](https://vk.com/translators_team). 
* Chinese Characters are from the font Hanyi Senty Rubber Candy by 王维 Wang Wei

**All rights belong to their respective owners.**