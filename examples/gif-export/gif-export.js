let duration = 5;
let margin = 3;

function setup(c) {
  c.size(40, 21);
  c.backgroundBg("197");
  c.saveCanvasToGIF("canvas.gif", duration);
}

function draw(c) {
  c.clear();
  c.stroke("RUNAL", "255", "197");
  let theta = c.loopAngle(duration);

  for (let y = margin; y < c.height - margin; y++) {
    let x = c.map(Math.sin(theta + y), -1, 1, margin, c.width - margin);

    c.point(x, y);
  }
}

function onKey(c, e) {
  if (e.key == "space") {
    c.noiseSeed(Date.now());
  }
}
