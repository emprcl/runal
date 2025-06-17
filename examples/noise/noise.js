function setup() {}

function draw() {
  for (let i = 0; i < c.width; i++) {
    for (let j = 0; j < c.height; j++) {
      let color = c.map(
        c.noise2D(
          i * 0.005 + c.framecount / 1000,
          j * 0.005 + c.framecount / 1000,
        ),
        0,
        1,
        0,
        255,
      );
      c.stroke("@", color, "#000000");
      c.point(i, j);
    }
  }
}

function onKey() {
  if (key == "c") {
    c.saveCanvas(`canvas_${Date.now()}.png`);
  }
  if (key == " ") {
    c.noiseSeed(Date.now());
    c.redraw();
  }
}
