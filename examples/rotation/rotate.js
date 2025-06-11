function setup() {
  //c.noLoop();
  c.cellPadding(" ");
}

function draw() {
  c.clear();
  c.stroke(".", "#ffffff", "#000000");
  c.fill(" ", "#ffffff", "#ffffff");
  c.translate(c.width / 2, c.height / 2);
  c.rotate(c.framecount * 0.5);
  c.circle(0, 0, 6);
}
