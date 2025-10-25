let duration = 5;
let radius = 4;
let seed1 = Date.now();
let seed2 = seed1 + 1000;

function setup(c) {
  c.cellModeDouble();
}

function draw(c) {
  c.clear();

  let theta = c.loopAngle(duration);
  c.noiseSeed(seed1);
  let noise = c.noiseLoop(theta, 1);
  c.noiseSeed(seed2);
  let noise2 = c.noiseLoop(theta, 1);
  let x = c.map(noise, 0, 1, radius, c.width - radius);
  let y = c.map(noise2, 0, 1, radius, c.height - radius);
  let color = c.map(noise, 0, 1, 0, 10);
  let colorBg = c.map(noise2, 0, 1, 200, 210);

  c.background("#", colorBg, colorBg);
  c.stroke(" ", color, color);
  c.fill(" ", color, color);
  c.circle(x, y, radius);
}

function onKey(c, e) {
  if (e.key == "space") {
    seed1 = Date.now();
    seed2 = seed1 + 1000;
  }
  if (e.key == "c") {
    c.saveCanvasToMP4("flash.mp4", duration);
  }
}
