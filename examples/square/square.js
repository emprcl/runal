function setup(c) {
  c.cellModeCustom(" ");
}

function draw(c) {
  c.clear();
  c.stroke("BORDER", "#ffffff", "#555555");
  c.fill("square", "#ffffff", "#000000");
  c.translate(c.width / 2, c.height / 2);
  c.rotate(c.framecount * 0.08);
  c.scale(1);
  c.square(0, 0, 10);
}
