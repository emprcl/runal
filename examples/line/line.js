function setup() {}

function draw() {
  c.flush();
  let y2 = (Math.sin(c.framecount * 0.1) / 2 + 0.5) * c.height * 0.8;
  c.line("0", 2, 2, c.width - 16, y2);
}
