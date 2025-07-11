let offset = 1000;

function setup(c) {}

function draw(c) {
  c.clear();
  offset = c.map(c.mouseY, 0, c.height, 0, 1000);
  c.fill("o", "#ffffff", "#000000");
  let angle = c.loopAngle(10);
  c.translate(c.width / 2, c.height / 2);
  c.rotate(angle);

  // NOT GOOD
  // c.triangle(
  //   c.width / 2,
  //   -offset,
  //   c.width + offset,
  //   c.height / 2,
  //   -offset,
  //   c.height / 2,
  // );

  c.circle(0, 0, offset);
}
