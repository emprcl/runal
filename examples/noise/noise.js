function setup(c) {
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

function onKey(c, key) {
  if (key == "c") {
    c.saveCanvasToPNG(`canvas_${Date.now()}.png`);
  }
  if (key == " ") {
    c.noiseSeed(Date.now());
    c.redraw();
  }
}
