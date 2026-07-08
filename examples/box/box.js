function setup(c) {
  c.cellModeDouble();
  c.light(-0.5, 0.8, 1);
}

function draw(c) {
  c.clear();
  c.fill("t", "#00ffcc", "#000000");
  c.translate(c.width / 2, c.height / 2);
  c.rotateX(c.framecount * 0.02);
  c.rotateY(c.framecount * 0.03);
  c.rotateZ(c.framecount * 0.01);
  c.box(20, 20, 20);
}
