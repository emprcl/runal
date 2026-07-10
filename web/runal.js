// runal.js — browser rendering backend for the runal wasm build.
//
// The Go side (internal/canvas/render_js.go) calls into the two functions
// exposed here, once per frame, through the global `__runal` object:
//
//   __runal.metrics(canvasEl, fontSize) -> [cols, rows, cellW, cellH]
//   __runal.draw(canvasEl, cols, rows, bytes, fontSize)
//
// `bytes` is a Uint8Array of cols*rows*12 bytes: for each cell, three
// little-endian uint32s — rune code point, packed 0xRRGGBB foreground,
// packed 0xRRGGBB background.

(function () {
  const BYTES_PER_CELL = 12;
  const colorCache = new Map();

  function hex(v) {
    let s = colorCache.get(v);
    if (s === undefined) {
      s = "#" + v.toString(16).padStart(6, "0");
      colorCache.set(v, s);
    }
    return s;
  }

  // Measure the monospace cell box for a given font size and lay the canvas
  // backing store out at device-pixel resolution. State is cached on the
  // element so draw() doesn't re-measure every frame.
  function setup(el, fontSize) {
    const dpr = window.devicePixelRatio || 1;
    const ctx = el.getContext("2d", { alpha: false });
    ctx.font = fontSize + "px monospace";

    const cellW = Math.max(1, Math.round(ctx.measureText("M").width));
    const cellH = Math.max(1, Math.round(fontSize * 1.2));

    const cssW = el.clientWidth || el.width || 640;
    const cssH = el.clientHeight || el.height || 384;

    const cols = Math.max(1, Math.floor(cssW / cellW));
    const rows = Math.max(1, Math.floor(cssH / cellH));

    const pxW = Math.floor(cssW * dpr);
    const pxH = Math.floor(cssH * dpr);
    if (el.width !== pxW) el.width = pxW;
    if (el.height !== pxH) el.height = pxH;

    const state = { ctx, cellW, cellH, dpr, fontSize, cssW, cssH, cols, rows };
    el.__runalRenderer = state;
    return state;
  }

  function download(blob, filename) {
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    a.remove();
    setTimeout(() => URL.revokeObjectURL(url), 2000);
  }

  // Force `name` to end in `.ext` (the actual produced format).
  function fixExt(name, ext) {
    if (!name) return "runal." + ext;
    return /\.[a-z0-9]+$/i.test(name)
      ? name.replace(/\.[a-z0-9]+$/i, "." + ext)
      : name + "." + ext;
  }

  // When a GIF recording is active on `el`, snapshot the just-rendered frame
  // (downscaled, quantized) and, once enough frames are collected, encode and
  // download the animated GIF. Called at the end of every draw().
  function captureGifFrame(el) {
    const g = el.__runalGif;
    if (!g) return;
    g.tmpCtx.drawImage(el, 0, 0, g.capW, g.capH);
    const img = g.tmpCtx.getImageData(0, 0, g.capW, g.capH);
    g.frames.push(RunalGIF.quantize(img.data, g.capW * g.capH));
    if (g.frames.length >= g.need) {
      el.__runalGif = null;
      const bytes = RunalGIF.encode(g.frames, g.capW, g.capH, g.delayCentis);
      download(new Blob([bytes], { type: "image/gif" }), g.filename);
    }
  }

  window.__runal = {
    metrics(el, fontSize) {
      const s = setup(el, fontSize);
      return [s.cols, s.rows, s.cellW, s.cellH];
    },

    draw(el, cols, rows, bytes, fontSize) {
      let s = el.__runalRenderer;
      if (!s || s.fontSize !== fontSize) s = setup(el, fontSize);

      const { ctx, cellW, cellH, dpr } = s;
      ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
      ctx.textBaseline = "top";
      ctx.font = fontSize + "px monospace";
      ctx.clearRect(0, 0, s.cssW, s.cssH);

      const view = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength);
      const n = cols * rows;
      let lastFill = null;

      for (let i = 0; i < n; i++) {
        const off = i * BYTES_PER_CELL;
        const rune = view.getUint32(off, true);
        const fg = view.getUint32(off + 4, true);
        const bg = view.getUint32(off + 8, true);

        const col = i % cols;
        const row = (i / cols) | 0;
        const x = col * cellW;
        const y = row * cellH;

        const bgc = hex(bg);
        if (bgc !== lastFill) {
          ctx.fillStyle = bgc;
          lastFill = bgc;
        }
        ctx.fillRect(x, y, cellW, cellH);

        if (rune !== 0 && rune !== 32) {
          const fgc = hex(fg);
          if (fgc !== lastFill) {
            ctx.fillStyle = fgc;
            lastFill = fgc;
          }
          ctx.fillText(String.fromCodePoint(rune), x, y);
        }
      }

      captureGifFrame(el);
    },

    // SaveCanvasToPNG: download the current frame as a PNG.
    savePNG(el, filename) {
      el.toBlob((blob) => download(blob, fixExt(filename, "png")), "image/png");
    },

    // SaveCanvasToGIF: record `duration` seconds of rendered frames at `fps`
    // and download an animated GIF. Frames are captured in draw().
    recordGIF(el, filename, duration, fps) {
      if (el.__runalGif) return; // already recording
      const s = el.__runalRenderer || setup(el, 16);
      const scale = Math.min(1, 400 / s.cssW); // cap width for size/memory
      const capW = Math.max(1, Math.round(s.cssW * scale));
      const capH = Math.max(1, Math.round(s.cssH * scale));
      const tmp = document.createElement("canvas");
      tmp.width = capW;
      tmp.height = capH;
      el.__runalGif = {
        frames: [],
        need: Math.max(1, Math.round(duration * fps)),
        delayCentis: Math.max(2, Math.round(100 / fps)),
        filename: fixExt(filename, "gif"),
        capW,
        capH,
        tmpCtx: tmp.getContext("2d", { willReadFrequently: true }),
      };
    },

    // SaveCanvasToMP4: record `duration` seconds of the canvas via
    // MediaRecorder and download. Emits mp4 if the browser supports it,
    // otherwise webm (and the filename extension is adjusted to match).
    recordVideo(el, filename, duration) {
      if (!el.captureStream || !window.MediaRecorder) {
        console.error("runal: video capture unsupported in this browser");
        return;
      }
      const types = [
        "video/mp4;codecs=avc1.42E01E",
        "video/mp4",
        "video/webm;codecs=vp9",
        "video/webm",
      ];
      const mime = types.find((t) => MediaRecorder.isTypeSupported(t)) || "";
      const rec = new MediaRecorder(el.captureStream(30), mime ? { mimeType: mime } : undefined);
      const chunks = [];
      rec.ondataavailable = (e) => e.data && e.data.size && chunks.push(e.data);
      rec.onstop = () => {
        const type = rec.mimeType || mime || "video/webm";
        const ext = type.includes("mp4") ? "mp4" : "webm";
        download(new Blob(chunks, { type }), fixExt(filename, ext));
      };
      rec.start();
      setTimeout(() => rec.state !== "inactive" && rec.stop(), Math.max(100, duration * 1000));
    },
  };
})();
