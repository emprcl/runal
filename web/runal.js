// runal.js — browser runtime for the runal wasm build.
//
// The sketch runs in the browser's native JS engine. The `c` canvas object is a
// thin proxy: each primitive call goes straight into the Go engine through the
// //go:wasmexport functions (see internal/canvas/web_js.go). Numbers cross
// directly; strings go via the exported scratch buffer; colors are parsed by
// the Go engine (rColor*/rResolveColor) and cached here. Each frame the Go
// engine writes the cell blob into wasm linear memory and we paint it to the
// <canvas>. Image/video export is done here in the browser.

(function () {
  const BYTES_PER_CELL = 12;
  const enc = new TextEncoder();

  let X = null; // wasm exports
  let cvs = null, glCtx = null;
  let S = null; // { cellW, cellH, dpr, cssW, cssH }
  let sketch = {};
  let proxy = null;
  let rafId = 0, running = false, looping = true;
  let fontSize = 16, autoResize = true, cellModeOn = false;
  let fontFamily = "monospace", fontSeq = 0;
  let mouseX = 0, mouseY = 0;
  const colorCache = new Map();

  function fontSpec() { return fontSize + "px " + fontFamily; }

  // ---------- wasm boot ----------
  async function boot(wasmUrl) {
    const go = new Go();
    const buf = await (await fetch(wasmUrl)).arrayBuffer();
    const res = await WebAssembly.instantiate(buf, go.importObject);
    go.run(res.instance); // initializes package state; main blocks (runtime stays resident)
    X = res.instance.exports;
  }

  // ---------- linear-memory helpers ----------
  function writeScratch(str) {
    const bytes = enc.encode(str);
    const ptr = X.rScratchPtr() >>> 0;
    new Uint8Array(X.mem.buffer, ptr, bytes.length).set(bytes);
    return bytes.length;
  }

  function resolveColor(c) {
    if (c === undefined || c === null) return 0;
    if (typeof c === "number") {
      const k = "n" + c;
      let v = colorCache.get(k);
      if (v === undefined) { v = X.rResolveColor(writeScratch("" + c)) >>> 0; colorCache.set(k, v); }
      return v;
    }
    let v = colorCache.get(c);
    if (v === undefined) {
      v = X.rResolveColor(writeScratch(String(c))) >>> 0;
      if (colorCache.size > 16384) colorCache.clear();
      colorCache.set(c, v);
    }
    return v;
  }

  function toHex(packed) { return "#" + (packed >>> 0).toString(16).padStart(6, "0"); }
  function colorFn(fn) {
    return (a, b, c) => {
      const p = fn(a, b, c) >>> 0;
      const hexStr = toHex(p);
      colorCache.set(hexStr, p); // pre-seed so the resulting string resolves without a crossing
      return hexStr;
    };
  }

  // ---------- canvas metrics ----------
  function measure() {
    const dpr = window.devicePixelRatio || 1;
    glCtx = cvs.getContext("2d", { alpha: false });
    glCtx.font = fontSpec();
    const cellW = Math.max(1, Math.round(glCtx.measureText("M").width));
    const cellH = Math.max(1, Math.round(fontSize * 1.2));
    const cssW = cvs.clientWidth || cvs.width || 640;
    const cssH = cvs.clientHeight || cvs.height || 384;
    S = { cellW, cellH, dpr, cssW, cssH };
    return S;
  }

  function displayGrid() {
    return { cols: Math.max(1, Math.floor(S.cssW / S.cellW)), rows: Math.max(1, Math.floor(S.cssH / S.cellH)) };
  }

  function sizeBackingStore() {
    const pxW = Math.floor(S.cssW * S.dpr), pxH = Math.floor(S.cssH * S.dpr);
    if (cvs.width !== pxW) cvs.width = pxW;
    if (cvs.height !== pxH) cvs.height = pxH;
  }

  // Re-measure the cell box after a font change and, when auto-resizing, refit
  // the engine grid to the new cell size so it keeps filling the canvas.
  function applyFontMetrics() {
    measure();
    if (autoResize) {
      const g = displayGrid();
      X.rResize(g.cols, g.rows);
      refreshDims();
    }
    sizeBackingStore();
  }

  // Load a custom font file (URL) via the FontFace API and switch the canvas to
  // it once ready. Loading is async; the current frame keeps the previous font
  // until it resolves, then metrics are recomputed.
  function loadFont(src) {
    if (!src) return;
    if (!window.FontFace) { console.error("runal: FontFace API unavailable"); return; }
    const family = "runal-font-" + ++fontSeq;
    const face = new FontFace(family, "url(" + JSON.stringify(String(src)) + ")");
    face.load()
      .then((f) => { document.fonts.add(f); fontFamily = family; applyFontMetrics(); })
      .catch((e) => console.error("runal: font load failed:", e));
  }

  // ---------- render ----------
  const hexCache = new Map();
  function hex(v) {
    let s2 = hexCache.get(v);
    if (s2 === undefined) { s2 = "#" + v.toString(16).padStart(6, "0"); hexCache.set(v, s2); }
    return s2;
  }

  function paint() {
    const cols = X.rDispCols(), rows = X.rRows();
    const ptr = X.rRender() >>> 0; // builds blob in wasm memory, returns pointer
    if (!ptr) return;
    // Re-view after rRender: the ArrayBuffer identity changes if memory grew.
    const bytes = new Uint8Array(X.mem.buffer, ptr, cols * rows * BYTES_PER_CELL);
    const view = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength);
    const { cellW, cellH, dpr, cssW, cssH } = S;

    glCtx.setTransform(dpr, 0, 0, dpr, 0, 0);
    glCtx.textBaseline = "top";
    glCtx.font = fontSpec();
    glCtx.clearRect(0, 0, cssW, cssH);

    let last = null;
    const n = cols * rows;
    for (let i = 0; i < n; i++) {
      const off = i * BYTES_PER_CELL;
      const rune = view.getUint32(off, true);
      const fg = view.getUint32(off + 4, true);
      const bg = view.getUint32(off + 8, true);
      const x = (i % cols) * cellW;
      const y = ((i / cols) | 0) * cellH;
      const bgc = hex(bg);
      if (bgc !== last) { glCtx.fillStyle = bgc; last = bgc; }
      glCtx.fillRect(x, y, cellW, cellH);
      if (rune !== 0 && rune !== 32) {
        const fgc = hex(fg);
        if (fgc !== last) { glCtx.fillStyle = fgc; last = fgc; }
        glCtx.fillText(String.fromCodePoint(rune), x, y);
      }
    }
    captureGifFrame();
  }

  // ---------- the sketch-facing proxy ----------
  function s(str) { return writeScratch(str == null ? "" : String(str)); }

  function buildProxy() {
    const p = {
      width: 0, height: 0, mouseX: 0, mouseY: 0, framecount: 0,

      // config
      size(w, h) { X.rSize(w, h); autoResize = false; refreshDims(); },
      clear() { X.rClear(); },
      fps(f) { X.rFps(f); },
      background(t, fg, bg) { X.rBackground(s(t), resolveColor(fg), resolveColor(bg)); },
      backgroundText(t) { X.rBackgroundText(s(t)); },
      backgroundFg(fg) { X.rBackgroundFg(resolveColor(fg)); },
      backgroundBg(bg) { X.rBackgroundBg(resolveColor(bg)); },

      // stroke / fill
      stroke(t, fg, bg) { X.rStroke(s(t), resolveColor(fg), resolveColor(bg)); },
      strokeText(t) { X.rStrokeText(s(t)); },
      strokeFg(fg) { X.rStrokeFg(resolveColor(fg)); },
      strokeBg(bg) { X.rStrokeBg(resolveColor(bg)); },
      fill(t, fg, bg) { X.rFill(s(t), resolveColor(fg), resolveColor(bg)); },
      fillText(t) { X.rFillText(s(t)); },
      fillFg(fg) { X.rFillFg(resolveColor(fg)); },
      fillBg(bg) { X.rFillBg(resolveColor(bg)); },
      noStroke() { X.rNoStroke(); },
      noFill() { X.rNoFill(); },

      // transforms
      push() { X.rPush(); },
      pop() { X.rPop(); },
      translate(x, y) { X.rTranslate(x, y); },
      rotate(a) { X.rRotate(a); },
      scale(f) { X.rScale(f); },

      // shapes
      point(x, y) { X.rPoint(x, y); },
      line(x1, y1, x2, y2) { X.rLine(x1, y1, x2, y2); },
      rect(x, y, w, h) { X.rRect(x, y, w, h); },
      square(x, y, sz) { X.rSquare(x, y, sz); },
      ellipse(cx, cy, rx, ry) { X.rEllipse(cx, cy, rx, ry); },
      circle(cx, cy, r) { X.rCircle(cx, cy, r); },
      triangle(x1, y1, x2, y2, x3, y3) { X.rTriangle(x1, y1, x2, y2, x3, y3); },
      quad(x1, y1, x2, y2, x3, y3, x4, y4) { X.rQuad(x1, y1, x2, y2, x3, y3, x4, y4); },
      bezier(x1, y1, x2, y2, x3, y3, x4, y4) { X.rBezier(x1, y1, x2, y2, x3, y3, x4, y4); },
      text(t, x, y) { X.rText(s(t), x, y); },

      // colors (parsed in the Go engine)
      colorRGB: colorFn((r, g, b) => X.rColorRGB(r, g, b)),
      colorHSL: colorFn((h, sa, l) => X.rColorHSL(h, sa, l)),
      colorHSV: colorFn((h, sa, v) => X.rColorHSV(h, sa, v)),

      // math (pure — identical formulas to the Go engine)
      map(v, a, b, c2, d) { return c2 + ((d - c2) / (b - a)) * (v - a); },
      dist(x1, y1, x2, y2) { const dx = x2 - x1, dy = y2 - y1; return Math.sqrt(dx * dx + dy * dy); },

      // noise / random (stateful — in the Go engine)
      noise1D(x) { return X.rNoise1D(x); },
      noise2D(x, y) { return X.rNoise2D(x, y); },
      noiseSeed(seed) { X.rNoiseSeed(BigInt(Math.trunc(seed))); },
      noiseLoop(a, r) { return X.rNoiseLoop(a, r); },
      noiseLoop1D(a, r, x) { return X.rNoiseLoop1D(a, r, x); },
      noiseLoop2D(a, r, x, y) { return X.rNoiseLoop2D(a, r, x, y); },
      loopAngle(d) { return X.rLoopAngle(d); },
      random(min, max) { return X.rRandom(min, max); },
      randomSeed(seed) { X.rRandomSeed(BigInt(Math.trunc(seed))); },

      // cell mode
      cellModeDouble() { X.rCellModeDouble(); cellModeOn = true; refreshDims(); },
      cellModeCustom(ch) { X.rCellModeCustom((ch || " ").codePointAt(0)); cellModeOn = true; refreshDims(); },
      cellModeDefault() { X.rCellModeDefault(); cellModeOn = false; refreshDims(); },
      cellPadding(ch) { this.cellModeCustom(ch); },
      cellPaddingDouble() { this.cellModeDouble(); },
      noCellPadding() { this.cellModeDefault(); },

      // loop control (driven here, in JS)
      loop() { looping = true; },
      noLoop() { looping = false; },
      redraw() { drawOnce(); },
      get isLooping() { return looping; },
      exit() { stop(); },

      // debug + export (browser-side)
      debug(...m) { console.log(...m); },
      saveCanvasToPNG(name) { savePNG(name); },
      saveCanvasToGIF(name, dur) { recordGIF(name, dur, X.rGetFps()); },
      saveCanvasToMP4(name, dur) { recordVideo(name, dur); },
      savedCanvasFont(src) { loadFont(src); },
      savedCanvasFontSize(px) { fontSize = px; applyFontMetrics(); },
    };
    return p;
  }

  function refreshDims() {
    proxy.width = X.rWidth();
    proxy.height = X.rHeight();
  }

  // ---------- loop ----------
  function drawOnce() {
    proxy.framecount = X.rFramecount();
    proxy.mouseX = mouseX;
    proxy.mouseY = mouseY;
    try {
      if (sketch.draw) sketch.draw(proxy);
    } catch (e) { reportError(e); return; }
    paint();
  }

  let lastFrame = -1e12;
  function frame(now) {
    if (!running) return;
    if (looping) {
      const interval = 1000 / Math.max(1, X.rGetFps());
      if (now - lastFrame >= interval) { drawOnce(); lastFrame = now; }
    }
    rafId = requestAnimationFrame(frame);
  }

  // ---------- input ----------
  const listeners = [];
  function on(el, type, fn) { el.addEventListener(type, fn); listeners.push([el, type, fn]); }
  function clearListeners() { for (const [el, t, fn] of listeners) if (fn) el.removeEventListener(t, fn); else el.removeEventListener(); listeners.length = 0; }

  function cellCoords(e) {
    let x = S.cellW > 0 ? Math.floor(e.offsetX / S.cellW) : 0;
    let y = S.cellH > 0 ? Math.floor(e.offsetY / S.cellH) : 0;
    if (cellModeOn) x = Math.floor(x / 2);
    return { x, y };
  }
  const BTN = { 0: "left", 1: "middle", 2: "right" };
  function mapKey(k) {
    const m = { " ": "space", ArrowUp: "up", ArrowDown: "down", ArrowLeft: "left", ArrowRight: "right", Enter: "enter", Escape: "esc", Backspace: "backspace", Tab: "tab", Delete: "delete" };
    if (m[k]) return m[k];
    return k.length === 1 ? k.toLowerCase() : k;
  }

  function attachInput() {
    if (cvs.tabIndex < 0) cvs.tabIndex = 0;
    on(cvs, "keydown", (e) => { if (sketch.onKey) safe(() => sketch.onKey(proxy, { key: mapKey(e.key), code: e.keyCode })); });
    on(cvs, "mousemove", (e) => { const p = cellCoords(e); mouseX = p.x; mouseY = p.y; if (sketch.onMouseMove) safe(() => sketch.onMouseMove(proxy, { x: p.x, y: p.y })); });
    on(cvs, "mousedown", (e) => { const p = cellCoords(e); mouseX = p.x; mouseY = p.y; if (sketch.onMouseClick) safe(() => sketch.onMouseClick(proxy, { x: p.x, y: p.y, button: BTN[e.button] || "left" })); });
    on(cvs, "mouseup", (e) => { const p = cellCoords(e); if (sketch.onMouseRelease) safe(() => sketch.onMouseRelease(proxy, { x: p.x, y: p.y, button: BTN[e.button] || "left" })); });
    on(cvs, "wheel", (e) => { const p = cellCoords(e); if (sketch.onMouseWheel) safe(() => sketch.onMouseWheel(proxy, { x: p.x, y: p.y, button: e.deltaY > 0 ? "down" : "up" })); });
    if (window.ResizeObserver) {
      const ro = new ResizeObserver(() => { if (autoResize) onResize(); });
      ro.observe(cvs);
      listeners.push([{ removeEventListener() { ro.disconnect(); } }, "", null]);
    }
  }

  function onResize() {
    measure();
    const g = displayGrid();
    X.rResize(g.cols, g.rows);
    refreshDims();
    sizeBackingStore();
  }

  function safe(fn) { try { fn(); } catch (e) { reportError(e); } }
  let onErrorCb = null;
  function reportError(e) { console.error(e); if (onErrorCb) onErrorCb(String(e && e.stack ? e.stack : e)); stop(); }

  // ---------- lifecycle ----------
  function start(source, canvasEl, opts) {
    stop();
    opts = opts || {};
    onErrorCb = opts.onError || null;
    fontSize = opts.fontSize || 16;
    cvs = canvasEl;
    autoResize = true; cellModeOn = false; looping = true; lastFrame = -1e12;
    fontFamily = "monospace"; // reset any custom font from a previous sketch

    // parse sketch source in the browser's JS engine
    try {
      const factory = new Function(
        source + "\n;return {" +
        ["setup", "draw", "onKey", "onMouseMove", "onMouseClick", "onMouseRelease", "onMouseWheel"]
          .map((n) => `${n}: typeof ${n}==='function'?${n}:null`).join(",") + "};"
      );
      sketch = factory();
    } catch (e) { reportError(e); return String(e); }
    if (!sketch.setup || !sketch.draw) { const m = "sketch must define setup and draw"; reportError(m); return m; }

    measure();
    const g = displayGrid();
    X.rInit(g.cols, g.rows);
    proxy = buildProxy();
    refreshDims();

    safe(() => sketch.setup(proxy)); // may call size / cellMode / fps
    sizeBackingStore();

    attachInput();
    running = true;
    drawOnce();
    rafId = requestAnimationFrame(frame);
    return "";
  }

  function stop() {
    running = false;
    if (rafId) cancelAnimationFrame(rafId);
    rafId = 0;
    clearListeners();
    if (cvs && cvs.__runalGif) cvs.__runalGif = null;
  }

  // ---------- export (browser-side) ----------
  function download(blob, filename) {
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url; a.download = filename;
    document.body.appendChild(a); a.click(); a.remove();
    setTimeout(() => URL.revokeObjectURL(url), 2000);
  }
  function fixExt(name, ext) {
    if (!name) return "runal." + ext;
    return /\.[a-z0-9]+$/i.test(name) ? name.replace(/\.[a-z0-9]+$/i, "." + ext) : name + "." + ext;
  }
  function savePNG(name) { cvs.toBlob((b) => download(b, fixExt(name, "png")), "image/png"); }

  function captureGifFrame() {
    const gs = cvs.__runalGif;
    if (!gs) return;
    gs.tmpCtx.drawImage(cvs, 0, 0, gs.capW, gs.capH);
    const img = gs.tmpCtx.getImageData(0, 0, gs.capW, gs.capH);
    gs.frames.push(RunalGIF.quantize(img.data, gs.capW * gs.capH));
    if (gs.frames.length >= gs.need) {
      cvs.__runalGif = null;
      download(new Blob([RunalGIF.encode(gs.frames, gs.capW, gs.capH, gs.delayCentis)], { type: "image/gif" }), gs.filename);
    }
  }
  function recordGIF(name, dur, fps) {
    if (cvs.__runalGif) return;
    const scale = Math.min(1, 400 / S.cssW);
    const capW = Math.max(1, Math.round(S.cssW * scale)), capH = Math.max(1, Math.round(S.cssH * scale));
    const tmp = document.createElement("canvas"); tmp.width = capW; tmp.height = capH;
    cvs.__runalGif = { frames: [], need: Math.max(1, Math.round(dur * fps)), delayCentis: Math.max(2, Math.round(100 / fps)), filename: fixExt(name, "gif"), capW, capH, tmpCtx: tmp.getContext("2d", { willReadFrequently: true }) };
  }
  function recordVideo(name, dur) {
    if (!cvs.captureStream || !window.MediaRecorder) { console.error("runal: video capture unsupported"); return; }
    const types = ["video/mp4;codecs=avc1.42E01E", "video/mp4", "video/webm;codecs=vp9", "video/webm"];
    const mime = types.find((t) => MediaRecorder.isTypeSupported(t)) || "";
    const rec = new MediaRecorder(cvs.captureStream(30), mime ? { mimeType: mime } : undefined);
    const chunks = [];
    rec.ondataavailable = (e) => e.data && e.data.size && chunks.push(e.data);
    rec.onstop = () => { const type = rec.mimeType || mime || "video/webm"; download(new Blob(chunks, { type }), fixExt(name, type.includes("mp4") ? "mp4" : "webm")); };
    rec.start();
    setTimeout(() => rec.state !== "inactive" && rec.stop(), Math.max(100, dur * 1000));
  }

  window.Runal = {
    boot, start, stop,
    // export helpers operating on the running sketch's canvas
    savePNG: (name) => savePNG(name),
    recordGIF: (name, dur, fps) => recordGIF(name, dur, fps || (X ? X.rGetFps() : 30)),
    recordVideo: (name, dur) => recordVideo(name, dur),
    framecount: () => (X ? X.rFramecount() : 0),
  };
})();
