function setup(c) {
  c.cellModeCustom(" ");
}

function draw(c) {
  c.clear();
  c.fill("ellipse", "#ffffff", "#000000");
  c.translate(c.width / 2, c.height / 2);
  c.rotate(c.framecount * 0.008);
  c.ellipse(0, 0, 15, 5);
}
