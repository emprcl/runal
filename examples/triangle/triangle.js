function setup() {
  c.cellPadding(" ");
}

function draw() {
  c.clear();
  c.stroke(".", "#ffffff", "#000000");
  c.fill(" ", "#ffffff", "#eeeeee");
  c.translate(c.width / 2, c.height / 2);
  c.rotate(c.framecount);
  c.scale(1);
  c.triangle(5, 5, 15, 15, 2, 15);
}
