function setup(c) {}

function draw(c) {
  c.clear();
  let y1 =
    (Math.sin((c.framecount * 0.2 + 1000) * 0.2) / 2 + 0.5) * c.height * 0.8;
  let x1 =
    (Math.cos((c.framecount * 0.2 + 1000) * 0.2) / 2 + 0.5) * c.width * 0.8;
  let y2 = (Math.sin(c.framecount * 0.1) / 2 + 0.5) * c.height * 0.8;
  let x2 = (Math.cos(c.framecount * 0.1) / 2 + 0.5) * c.width * 0.8;
  c.line(x1, y1, x2, y2);
}
