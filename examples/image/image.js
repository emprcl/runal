let img;

function setup(c) {
  img = c.loadImage("the-great-wave-kanagawa.jpg");
  c.noLoop();
}

function draw(c) {
  c.clear();
  c.image(img, 0, 0, c.width, c.height);
  let fullCanvas = c.get(0, 0, c.width, c.height);
  c.set(c.width / 2, 0, fullCanvas);
}
