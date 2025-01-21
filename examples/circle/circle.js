function setup() {
  c.widthPadding(".");
}

function draw() {
  c.clear();
  let radius1 = ((Math.sin(c.framecount * 0.1) * 0.5 + 0.5) * c.width) / 2;
  let radius2 = ((Math.sin(c.framecount * 0.2) * 0.5 + 0.5) * c.width) / 3;
  c.circle("0", "i", c.width / 2, c.height / 2, radius1);
  c.circle("C", "vuivu", c.width / 2, c.height / 2, radius2);
}
