function setup(c) {
  c.savedCanvasFontSize(24);
  c.cellPaddingDouble();
}

function draw(c) {
  for (let i = 0; i < c.width; i++) {
    for (let j = 0; j < c.height; j++) {
      let color = c.map(
        c.noise2D(
          i * 0.009 + c.framecount / 1000,
          j * 0.009 + c.framecount / 1000,
        ),
        0,
        1,
        150,
        231,
      );
      c.stroke("§", color, "#000000");
      c.point(i, j);
    }
  }
}

function onKey(c, e) {
  if (e.Key == "c") {
    c.saveCanvasToPNG(`canvas_${Date.now()}.png`);
  }
  if (e.key == "space") {
    c.noiseSeed(Date.now());
    c.redraw();
  }
}

function onMouse(c, e) {
  if (e.type == "click") {
    c.noiseSeed(Date.now());
    c.redraw();
  }
}
