# anyconvert

A self-hosted file converter (images, video, audio) built as a personal tool — and an exercise in Docker image optimization that got slightly out of hand.

## The optimization story

This project went through four major rewrites, driven by one moment of horror:

```
$ docker images
REPOSITORY    TAG       SIZE
anyconvert    v1        545MB
```

What followed was an obsessive journey to make that number smaller.

### Image sizes (API)

| Version | Image Size | Stack | Key Change |
|---------|-----------|-------|------------|
| **v1** | **545 MB** | Python 3.10, ImageMagick (compiled from source) | "It works" |
| **v2** | **122 MB** | Python 3.10 (compiled from source with LTO), FFmpeg | Custom minimal CPython runtime |
| **v3** | **22 MB** | Go, FFmpeg (compiled from source) | Go rewrite + surgical FFmpeg build |
| **v4** | **26 MB** | Go, FFmpeg (compiled from source) + MP3 support | Added lame codec, Kubernetes/Helm |

**545 MB → 22 MB (96% reduction).** The v3 API image ended up smaller than its own React frontend (24.5 MB).

### What happened at each version

**v1 — Nov 2023** — Python + ImageMagick
- Multi-arch build (ARM64/AMD64) with separate base images
- ImageMagick compiled from source via [imei.sh](https://github.com/SoftCreatR/imei)
- Used Gradio for the UI
- Produced a 545 MB image and a "what the actual fuck" moment

**v2 — Dec 2023** — Python, but minimal
- Compiled CPython 3.10 from source with aggressive flags:
  ```
  --with-ensurepip=no --without-doc-strings --without-mimalloc
  --without-pymalloc --without-readline --disable-test-modules
  --with-lto
  ```
- Linked with `LDFLAGS="-Wl,--strip-all"` to strip all symbols
- Stripped all `.pyc`, `.pyo`, `.dist-info`, `__pycache__` from site-packages
- Switched from ImageMagick to FFmpeg (from Alpine packages)
- Split into API + frontend (Vue/Vite)
- **Result:** 122 MB (78% smaller than v1)

**v3 — Dec 2023** — Go rewrite + custom FFmpeg
- Rewrote the entire API in Go (Gin framework)
- Compiled FFmpeg from source with `--disable-everything`, then surgically re-enabled only the exact codecs needed:
  ```
  --enable-protocol=file
  --enable-encoder=libx264,aac,gif
  --enable-decoder=libx264,webp,vp9,vorbis,gif,png,apng,zlib
  --enable-filter=pad,scale
  --enable-muxer=mp4,gif,mov
  --enable-demuxer=gif,mov,image2,matroska
  --enable-parser=png,webp,aac,gif,h264,vorbis
  ```
- Go binary built with `-ldflags="-s -w"` (strip debug info + symbols)
- Swagger docs auto-generated in a separate build stage
- **Result:** 22 MB — smaller than the React frontend

**v4 — Jan 2024** — Kubernetes deployment
- Added MP3 support (lame encoder/decoder)
- Added Helm chart for Kubernetes deployment
- Redis for job queue with configurable TTLs
- Docker Compose with network isolation and health checks
- **Result:** 26 MB (4 MB increase for MP3 codec support)

## Architecture (v4)

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Frontend   │────▶│   API (Go)   │────▶│  Redis       │
│  Vue + Nginx │     │  Gin + FFmpeg│     │  Job Queue   │
│   24.5 MB    │     │    26 MB     │     │   Alpine     │
└─────────────┘     └─────────────┘     └─────────────┘
```

- **API:** Go (Gin), custom-compiled FFmpeg, Swagger docs
- **Frontend:** Vue 3 + Vite, served via Nginx with reverse proxy to API
- **Queue:** Redis 7 (Alpine) for async job processing
- **Deployment:** Docker Compose (local) or Helm chart (Kubernetes)

## Running

```bash
docker compose up --build
```

The frontend is available at `http://localhost` (or configure a Traefik host rule), the API at `:8000`.

## Version history

Each major version is preserved in its own branch:

- [`v1`](../../tree/v1) — Python + ImageMagick monolith (545 MB)
- [`v2`](../../tree/v2) — Minimal Python runtime + FFmpeg (122 MB)
- [`v3`](../../tree/v3) — Go rewrite + surgical FFmpeg (22 MB)
- [`v4`](../../tree/v4) — v3 + MP3 support + Kubernetes (26 MB) ← **you are here (main)**
