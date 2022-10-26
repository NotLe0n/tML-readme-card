# tML-readme-card
This is a little widget to display the download count of your all of your tmodloader mods.

Use `https://tml-readme-card.repl.co/?steamid64=<your steam64id>` to get the widget

Mod names with color chat tags will have their color changed.

## Parameters
| Parameter       | Value                         | Effect                                              |
|-----------------|-------------------------------|-----------------------------------------------------|
| `steamid64`     | number                        | The user of which to display the mods from          |
| `text_color`    | Hex color value without the # | Changes the color of all text, except the mod names |
| `bg_color`      | Hex color value without the # | Changes the background color                        |
| `border_color`  | Hex color value without the # | Changes the border color                            |
| `border_width`  | number                        | Border width in pixels                              |
| `corner_radius` | number                        | corner radius, 0 for Rectangle                      |
| `v`             | "1.3" or "1.4"                | changes the tML version                             |
| `font`          | "Andy" or "Serif"             | changes the font                                    |

## Examples
* `https://tml-readme-card.repl.co/?steamid64=76561198278789341`
![example-widget](https://tml-readme-card.repl.co/?steamid64=76561198278789341&)

* `https://tml-readme-card.repl.co/?steamid64=76561198278789341&border_width=1&corner_radius=60&border_color=FFFFFF&bg_color=0D1116`
![example-widget-parameters](https://tml-readme-card.repl.co/?steamid64=76561198278789341&border_width=1&corner_radius=60&border_color=FFFFFF&bg_color=0D1116&)
