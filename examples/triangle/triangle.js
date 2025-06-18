function setup(c) {
  c.cellPadding(" ");
}

function draw(c) {
  c.clear();
  c.stroke("test", "#ffffff", "#000000");
  c.fill("blop", "#ffffff", "#000000");
  c.translate(c.width / 2, c.height / 2);
  c.rotate(c.framecount);
  c.scale(1);
  c.triangle(5, 5, 15, 15, 2, 15);
}
