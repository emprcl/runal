function setup() {
  c.cellPadding(".");
}

function draw() {
  c.clear();
  c.stroke("A", "#ffffff", "#555555");
  c.fill("B", "#ffffff", "#eeeeee");
  c.translate(c.width / 2, c.height / 2);
  c.rotate(c.framecount);
  c.scale(1);
  c.square(0, 0, 10);
}
