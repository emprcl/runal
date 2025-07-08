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

function onKey(c, key) {
  console.log(key);
  if (key == " ") {
    c.redraw();
  }
}
