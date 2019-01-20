# epaper
Driver of Waveshare Electronics e-paper display for Raspberry Pi - in Go lang

**Work in progress** - Only 2.9" Waveshare Electronics BW display supported for now.

### What it can do (so far)

package `epaper` (comunicates with display over SPI):

  - Initialize e-paper display to use either `full` or `partial` update
  - Swap frame buffer of e-paper display
  - Clear frame buffer using black / or white color
  - Display arbitraty monochromatic bitmap image
  - Put display to Sleep
  
package `epaper/image` (creates in-memmory monochromatic bitmap `image.Mono`):
  - Clear whole image to black or white color
  - Draw black or white horizontal / vertical **lines**
  - Draw black or white stroked / filled **rectangle**
  - Draw black or white stroked / filled **circle**
  - Write black or white **text** using Go font (chars from [WGL4](https://en.wikipedia.org/wiki/Windows_Glyph_List_4) charset)
  - **Paste another image** (while converting it to monochromatic color mode) using go's `image.Image` interface.
  - **Rotate** bitmap 90° in each direction
  - **Flip** (mirror) bitmap vertically or horizontally
  - **Invert** colors
  
<img src="/../images/photo.jpg" height="296"/><img src="/../images/image.png" height="296"/>

### Wiring 

| e-paper | Raspberry Pi |
|---------|--------------|
| 3.3V    | 3v3          |
| GND     | Ground       |
| DIN     | MOSI (BCM 10)|
| CLK     | SCLK (BCM 11)|
| CS      | CE0 (BCM 8)  |
| CD      | BCM 25       |
| RST     | BCM 22       |
| BUSY    | BCM 24       |             

Note that RST is on BCM 22 instead of BCM 17 as in https://pinout.xyz/pinout/213_inch_e_paper_phat - the rest is the same.
