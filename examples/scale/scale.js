function setup() {
  //c.noLoop();
  c.cellPadding(" ");
}

function draw() {
  c.clear();
  c.stroke(" ", "255", "0");
  c.circle(0, 0, 5);
  c.circle(10, 10, 5);
  c.scale((Math.sin(c.framecount * 0.1) / 2 + 0.5) * 10);
  c.rotate(c.framecount * 0.2);
}
