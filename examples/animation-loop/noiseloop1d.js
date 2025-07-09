let duration = 1;
let scale = 1;

function setup(c) {}

function draw(c) {
  c.clear();
  c.strokeText("THIS IS A 1D NOISE LOOP EXAMPLE");
  let theta = c.loopAngle(duration);

  for (let x = 0; x < c.width; x++) {
    let noise = c.noiseLoop1D(theta, 0.1, x * scale);
    let y = c.map(noise, 0, 1, 0, c.height);
    c.point(x, y);
  }
}

function onKey(c, e) {
  if (e.key == "space") {
    c.noiseSeed(Date.now());
  }
}

function onMouse(c, e) {
  // Handle mouse events here
}
