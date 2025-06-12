function setup() {
  c.cellPadding(" ");
}

function draw() {
  c.clear();
  c.stroke("TEST", "#ffffff", "#555555");
  c.fill("blop", "#ffffff", "#000000");
  c.translate(c.width / 2, c.height / 2);
  c.rotate(c.framecount);
  c.scale(1);
  c.ellipse(0, 0, 10, 20);
}
