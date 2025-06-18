function setup(c) {
  //c.noLoop();
  c.cellPadding(" ");
}

function draw(c) {
  c.clear();
  c.stroke(".", "#ffffff", "#ffffff");
  c.fill(".", "#ffffff", "#000000");
  c.translate(c.width / 2, c.height / 2);
  c.rotate(c.framecount * 0.5);
  c.circle(0, 0, 6);
}
