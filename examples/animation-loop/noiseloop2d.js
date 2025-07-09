let duration = 3;
let scale = 0.3;

function setup(c) {}

function draw(c) {
  c.clear();
  let theta = c.loopAngle(duration);

  for (let x = 0; x < c.width; x++) {
    for (let y = 0; y < c.height; y++) {
      let noise = c.noiseLoop2D(theta, 1, x * scale, y * scale);
      let color = c.map(noise, 0, 1, 232, 255);
      c.stroke("§", color, 0);
      c.point(x, y);
    }
  }
}

function onKey(c, e) {
  if (e.Key == " ") {
    c.noiseSeed(Date.now());
  }
}

function onMouse(c, e) {
  // Handle mouse events here
}
