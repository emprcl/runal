function setup(c) {
  c.noLoop();
}

function draw(c) {
  c.clear();
  for (let i = 0; i < c.width; i++) {
    for (let j = 0; j < c.height; j++) {
      if (Math.random() < 0.8) {
        c.text(".", i, j);
      }
    }
  }
}

function onKey(c, e) {
  if (e.Key == " ") {
    c.redraw();
  }
}

function onMouse(c, e) {
  c.redraw();
}
