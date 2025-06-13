function setup() {
  c.cellPadding(" ");
}

function draw() {
  c.clear();
  c.stroke("TEST", "#ffffff", "#555555");
  c.fill("blop", "#ffffff", "#000000");
  //c.rotate(c.framecount);
  //c.scale(2);
  c.quad(
    c.map(Math.sin(c.framecount * 0.1), -1, 1, 1, 35),
    1,
    c.map(Math.cos(c.framecount * 0.1), -1, 1, 1, 35),
    3,
    16,
    12,
    2,
    18,
  );
}
