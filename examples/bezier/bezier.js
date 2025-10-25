function setup(c) {
  c.cellModeCustom(" ");
}

function draw(c) {
  c.clear();
  c.stroke("0", "#ffffff", "#000000");
  c.bezier(10, 10, 20, 0, 30, 20, 40, 10);
}
