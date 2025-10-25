function setup(c) {
  c.cellModeCustom(" ");
}

function draw(c) {
  c.clear();
  c.stroke(".", "#ffffff", "#000000");
  c.fill("triangle", "#ffffff", "#000000");
  c.translate(c.width / 2, c.height / 2);
  c.rotate(c.framecount * 0.008);
  c.scale(1);
  c.triangle(5, 5, 15, 15, 2, 15);
}
