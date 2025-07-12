let img;

function setup(c) {
  img = c.loadImage("wish.png");
}

function draw(c) {
  c.clear();
  c.translate(c.width / 2, c.height / 2);
  c.rotate(c.framecount * 0.08);
  c.image(img, 0, 0, 40, 40);
}
