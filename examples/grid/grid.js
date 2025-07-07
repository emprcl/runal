function setup(c) {
  c.nLoop();
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

function onKey(c, key) {
  if (key == " ") {
    c.redraw();
  }
}
