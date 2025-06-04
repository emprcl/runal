function setup() {
  //c.noLoop();
  c.cellPadding(" ");
}

function draw() {
  c.clear();
  c.stroke(" ", "255", "0");
  c.circle(0, 0, 5);
  c.circle(10, 10, 5);
  c.scale(c.map(Math.sin(c.framecount * 0.1), -1, 1, 0, 10));
  c.rotate(c.framecount * 0.2);
}
