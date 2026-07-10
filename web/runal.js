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
    },
  };
})();
