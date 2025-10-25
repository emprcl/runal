function setup(c) {
  c.cellModeCustom(" ");
}

function draw(c) {
  c.clear();
  c.stroke(".", "#ffffff", "#555555");
  c.fill("quad", "#ffffff", "#000000");
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
