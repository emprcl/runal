function setup(c) {
  c.noLoop();
  c.cellModeCustom(" ");
}

function draw(c) {
  c.clear();
  c.circle(0, 0, 5);
  c.translate(c.width / 2, c.height / 2);
  c.circle(0, 0, 5);
}
