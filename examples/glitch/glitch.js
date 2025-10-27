let img;
let glitchPoints = 80;
let glitchSize = 2;

let glitchLines = 5;

function setup(c) {
  img = c.loadImage("mona-lisa.jpg");
  c.fps(5);
}

function draw(c) {
  c.clear();
  c.image(img, c.random(-2, 2), c.random(-2, 2), c.width, c.height);

  for (let i = 0; i < glitchPoints; i++) {
    let part = c.get(
      c.random(0, c.width),
      c.random(0, c.height),
      glitchSize,
      c.random(1, glitchSize + 1),
    );
    c.set(c.random(0, c.width), c.random(0, c.height), part);
  }

  for (let i = 0; i < glitchLines; i++) {
    let part = c.get(0, c.random(0, c.height), c.width, 1);
    c.set(0, c.random(0, c.height), part);
  }
}
