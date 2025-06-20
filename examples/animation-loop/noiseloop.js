let duration = 3;
let radius = 4;

function setup(c) {
  c.cellPadding(".");
}

function draw(c) {
  c.clear();

  let theta = c.loopAngle(duration);
  let noise = c.noiseLoop(theta, 1);

  let x = c.map(noise, 0, 1, radius, c.width - radius);
  let y = c.map(Math.sin(theta), -1, 1, radius, c.height - radius);
  let color = c.map(noise, 0, 1, 0, 255);

  c.background("#", 3, 3);
  c.stroke(" ", color, color);
  c.fill(" ", color, color);
  c.circle(x, y, radius);
}

function onKey(c, key) {
  if (key == " ") {
    c.noiseSeed(Date.now());
  }
}
