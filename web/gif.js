// gif.js — a small self-contained animated GIF89a encoder.
//
// Used by runal.js to implement SaveCanvasToGIF on the web build without
// pulling any image/gif encoder into the wasm binary. Frames are quantized to
// a fixed 3-3-2 RGB palette (256 colors, no dithering) — good enough for the
// flat colors of textmode output — and compressed with standard GIF LZW.
//
// window.RunalGIF = { quantize(rgba, n), encode(frames, w, h, delayCentis) }

(function () {
  // Fixed 3-3-2 RGB global color table: 8 R levels, 8 G levels, 4 B levels.
  const PALETTE = (function () {
    const p = new Uint8Array(256 * 3);
    for (let i = 0; i < 256; i++) {
      p[i * 3] = (((i >> 5) & 7) * 255 / 7) | 0;
      p[i * 3 + 1] = (((i >> 2) & 7) * 255 / 7) | 0;
      p[i * 3 + 2] = ((i & 3) * 255 / 3) | 0;
    }
    return p;
  })();

  // Map an RGBA buffer (length n*4) to n palette indices.
  function quantize(rgba, n) {
    const out = new Uint8Array(n);
    for (let i = 0, j = 0; i < n; i++, j += 4) {
      out[i] = ((rgba[j] >> 5) << 5) | ((rgba[j + 1] >> 5) << 2) | (rgba[j + 2] >> 6);
    }
    return out;
  }

  // GIF LZW using the "early-clear" scheme: codes are a fixed
  // (minCodeSize+1) bits wide, and a clear code is emitted before the decoder's
  // table would reach 2^codeSize (which would force a width change). This keeps
  // the encoder and the (inherently one-entry-behind) decoder perfectly in sync
  // with no variable-width bookkeeping — simple and provably correct. It trades
  // compression ratio for robustness, which is the right call for an export.
  function lzw(minCodeSize, pixels) {
    const clear = 1 << minCodeSize;
    const eoi = clear + 1;
    const size = minCodeSize + 1;
    const out = [];
    let acc = 0,
      nbits = 0;

    function emit(code) {
      acc |= code << nbits;
      nbits += size;
      while (nbits >= 8) {
        out.push(acc & 0xff);
        acc >>= 8;
        nbits -= 8;
      }
    }

    emit(clear);
    let sinceClear = 0;
    for (let i = 0; i < pixels.length; i++) {
      emit(pixels[i]);
      if (++sinceClear === clear - 2) {
        emit(clear);
        sinceClear = 0;
      }
    }
    emit(eoi);
    if (nbits > 0) out.push(acc & 0xff);
    return out;
  }

  // Split a byte array into GIF sub-blocks (<=255 bytes, length-prefixed,
  // 0x00-terminated) and append them to `bytes`.
  function writeSubBlocks(bytes, data) {
    for (let i = 0; i < data.length; i += 255) {
      const chunk = data.slice(i, i + 255);
      bytes.push(chunk.length);
      for (const b of chunk) bytes.push(b);
    }
    bytes.push(0);
  }

  function u16(bytes, v) {
    bytes.push(v & 0xff, (v >> 8) & 0xff);
  }

  // frames: array of Uint8Array palette-index buffers (each length w*h).
  function encode(frames, w, h, delayCentis) {
    const b = [];
    // Header
    for (const c of "GIF89a") b.push(c.charCodeAt(0));
    // Logical Screen Descriptor: global color table, 256 entries (size code 7)
    u16(b, w);
    u16(b, h);
    b.push(0xf7, 0, 0); // packed (GCT flag=1, color res=7, size=7), bg, aspect
    for (let i = 0; i < 768; i++) b.push(PALETTE[i]);
    // NETSCAPE2.0 looping extension (loop forever)
    b.push(0x21, 0xff, 0x0b);
    for (const c of "NETSCAPE2.0") b.push(c.charCodeAt(0));
    b.push(0x03, 0x01, 0x00, 0x00, 0x00);

    for (const frame of frames) {
      // Graphic Control Extension (delay, no transparency)
      b.push(0x21, 0xf9, 0x04, 0x00);
      u16(b, delayCentis);
      b.push(0x00, 0x00);
      // Image Descriptor (full frame, no local color table)
      b.push(0x2c);
      u16(b, 0);
      u16(b, 0);
      u16(b, w);
      u16(b, h);
      b.push(0x00);
      // Image data
      b.push(8); // LZW minimum code size
      writeSubBlocks(b, lzw(8, frame));
    }

    b.push(0x3b); // trailer
    return new Uint8Array(b);
  }

  window.RunalGIF = { quantize, encode, PALETTE };
})();
