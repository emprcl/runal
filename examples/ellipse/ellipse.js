function setup(c) {
  c.cellPadding(" ");
}

function draw(c) {
  c.clear();
  //c.stroke("TEST", "#ffffff", "#555555");
  c.fill("blop", "#ffffff", "#000000");
  c.rotate(c.framecount);
  //c.scale(2);
  c.ellipse(10, 5, 15, 5);
  c.circle(10, 35, 5);
}
