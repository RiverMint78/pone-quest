# PoneQuest

MLP meme search, inspired by MemeMeow.

## Frontend build (Bun)

前端采用现代化本地构建链路：Bun + Tailwind CSS v4 + TypeScript，
不依赖 CDN、默认生产最小化、尽量零配置。

### Install

```bash
bun install
```

### Production build (minify)

```bash
bun run build
```

### Development watch

```bash
bun run dev
```

### Output files

- CSS source: `ui/src/style.css`
- CSS output: `ui/static/css/style.css`
- JS source entry: `ui/src/app.ts`
- JS output: `ui/static/js/ui.js`
