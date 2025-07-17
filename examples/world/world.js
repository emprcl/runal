function setup(c) {
  c.size(82, 41);
  c.background("/", "237", "#000000");
  c.stroke("0", "255", "#000000");
  c.cellPaddingDouble();
}

function draw(c) {
  c.clear();
  if (c.framecount % 30 < 15) {
    c.backgroundText("\\");
  } else {
    c.backgroundText("/");
  }
  let theta = c.loopAngle(10);
  for (let i = 0; i < c.width; i++) {
    for (let j = 0; j < c.height; j++) {
      let color = c.map(
        c.noise2D(
          c.framecount * 0.03 + i * 0.05,
          c.framecount * 0.03 + j * 0.05,
        ),
        0,
        1,
        231,
        255,
      );
      if (c.dist(i, j, c.width / 2, c.height / 2) <= c.width / 2 - 4) {
        c.strokeFg(color);
        c.point(i, j);
      }
    }
  }
}

function onKey(c, e) {
  if (e.key == "c") {
    c.saveCanvasToPNG(`canvas_${Date.now()}.png`);
  }
  if (e.key == "space") {
    c.noiseSeed(Date.now());
    c.redraw();
  }
}
