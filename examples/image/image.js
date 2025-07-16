let img;

function setup(c) {
  img = c.loadImage("wish.png");
  c.noLoop();
  c.cellPaddingDouble();
}

function draw(c) {
  c.clear();
  c.image(img, 0, 0, c.width, c.height);
  let fullCanvas = c.get(0, 0, c.width, c.height);
  c.image(fullCanvas, c.width / 2, 0, c.width / 2, c.height);
}
