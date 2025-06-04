function setup() {
  c.noLoop();
  c.cellPadding(" ");
}

function draw() {
  c.clear();
  c.circle(0, 0, 5);
  c.translate(c.width / 2, c.height / 2);
  c.circle(0, 0, 5);
}
