function setup() {
  c.cellPadding(" ");
  c.noLoop();
}

function draw() {
  for (let i = 0; i < c.width; i++) {
    for (let j = 0; j < c.height; j++) {
      let color = c.map(c.noise2D(i * 0.005, j * 0.005), 0, 1, 0, 255);
      c.stroke("0", color, "#000000");
      c.point(i, j);
    }
  }
}
