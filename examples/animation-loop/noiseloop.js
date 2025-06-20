let duration = 1;
let radius = 4;

function setup(c) {
  c.cellPadding(".");
}

function draw(c) {
  c.clear();
  c.strokeText("0");
  let theta = c.loopAngle(duration);

  let x = c.map(c.noiseLoop(theta, 1), 0, 1, radius, c.width - radius);
  let y = c.map(Math.sin(theta), -1, 1, radius, c.height - radius);

  c.circle(x, y, radius);
}

function onKey(c, key) {
  if (key == " ") {
    c.noiseSeed(Date.now());
  }
}
