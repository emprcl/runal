let img;

function setup(c) {
  img = c.loadImage("../examples/image/wish.png");
  c.noLoop();
  c.cellPaddingDouble();
}

function draw(c) {
  c.clear();
  //c.translate(c.width / 2, c.height / 2);
  //c.rotate(c.framecount * 0.08);
  c.image(img, 0, 0, c.width - 2, c.height - 2);
}
