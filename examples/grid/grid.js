function setup() {
  c.noLoop();
}

function draw() {
  c.flush();
  for (let i = 0; i < c.width; i++) {
    for (let j = 0; j < c.height; j++) {
      c.text(".", i, j);
    }
  }
}
