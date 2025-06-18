function setup(c) {
  //c.noLoop();
  c.cellPadding(" ");
}

function draw(c) {
  c.clear();
  c.stroke(".", "#ffffff", "#ffffff");
  c.fill(".", "#ffffff", "#000000");
  c.translate(c.width / 2, c.height / 2);
  c.scale(c.map(Math.sin(c.framecount * 0.1), -1, 1, 1, 4));
  c.rotate(c.framecount * 0.008);
  c.circle(0, 0, 5);
  c.circle(10, 10, 5);
}
