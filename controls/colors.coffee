colors = require '../data/colors.json'

osum = (a,b)->
  r = Object.keys(a)
  .map (k)-> {[k]: (a[k] ? 0) + (b[k] ? 0) }
  Object.assign {}, r...

objectsSum = (objects)->
  r = objects.reduce osum, black
  for k of r
    r[k] = r[k] / objects.length
  r

sum = ->
  objectsSum Array.from arguments

for k, v of colors
  colors[k] =
    r: v.rgb[0]
    g: v.rgb[1]
    b: v.rgb[2]

  black = {r: 0, g:0, b:0, a:0, w:0, uv:0}
  red = {r: 255, g:0, b:0, a:0, w:0, uv:0}
  green = {r: 0, g:255, b:0, a:0, w:0, uv:0}
  blue = {r: 0, g:0, b:255, a:0, w:0, uv:0}
  uv = {r: 0, g:0, b:0, a:0, w:0, uv:255}
  amber = {r: 0, g:0, b:0, a:255, w: 0, uv: 0}
  cyan = {r: 0, g:255, b:255, a:0, w:0, uv:0}
  purple = {r: 255, g:0, b:255, a:0, w:0, uv:0}
  yellow_green = {r: 255, g:255, b:0, a:0, w:0, uv:0}
  yellow = {r: 255, g:127, b:0, a:0, w:0, uv:0}
  amber2 = {r: 255, g:191, b:0 }
  gold = {r: 255, g:109, b:0, a:64, w:0, uv:0}
  violet = {r:20, b: 10, g: 0, w: 0, a:0, uv: 255}
  pink = {r: 255, g:0, b:0, a:0, w:64, uv:0}
  orange = {r: 255, g:64, b:0, a:64, w:0, uv:64}
  dark_purple = {r: 50, g:0, b:50, a:0, w:0, uv:255}
  white_purple = {r:255, g:192, b:203}
  dark_pink = {r: 255, g:0, b:0, a:0, w:0, uv:255}
  dark_green = {r:0, g:128, b:0, a:20, uv: 0}
  pink_white = {r: 255, g:127, b:64, a:30, w:0, uv:0}
  soft_white = {r: 255, g:120, b:0, a:37, w:72, uv:0}
  gold_white = {r: 255, g:127, b:0, a:36, w:72, uv:0}
  white = {r: 128, g:128, b:128, a:0, w:255, uv:0}
  red_uv = {r:255, g:0, b:0, w: 0, uv:255}
  candle = {r: 100, g: 56, b: 6, a: 200, w: 10, uv:8 }
  green_brown = {r: 20, g: 255, b: 6, a: 100, w: 0, uv:8 }

  dark_orange = sum orange, black
  deep_orange = sum yellow, red, red
  dark_cyan = sum cyan, black
  light_cyan = sum cyan, yellow
  dark_purple2 = sum purple, black
  dark_green2 = sum green, black
  dark_green3 = sum green, green, red, black
