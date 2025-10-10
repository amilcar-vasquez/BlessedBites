# Manifest icons

This folder contains scalable SVG icons used by the PWA manifest. For maximum compatibility (especially with older Android/Chrome versions), you may want to generate PNG versions at 192x192 and 512x512.

Commands (ImageMagick) to generate PNGs from the SVGs:

```bash
# from project root
mkdir -p ui/static/img/manifest/png
convert ui/static/img/manifest/icon-192.svg -resize 192x192 ui/static/img/manifest/png/icon-192.png
convert ui/static/img/manifest/icon-512.svg -resize 512x512 ui/static/img/manifest/png/icon-512.png
```

After generating PNGs, update `ui/static/manifest.json` to reference the PNG files first for better compatibility.

Example snippet for manifest icons:

```json
"icons": [
  {
    "src": "/static/img/manifest/png/icon-192.png",
    "sizes": "192x192",
    "type": "image/png"
  },
  {
    "src": "/static/img/manifest/png/icon-512.png",
    "sizes": "512x512",
    "type": "image/png"
  }
]
```
