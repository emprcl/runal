function setup(c) {
  c.cellPadding(".");
}

function draw(c) {
  c.clear();

  let radius1 = ((Math.sin(c.framecount * 0.1) * 0.5 + 0.5) * c.width) / 2;
  let radius2 = ((Math.sin(c.framecount * 0.2) * 0.5 + 0.5) * c.width) / 3;

  c.stroke("COUCOU", "#ffffff", "#000000");
  c.fill("i", "#ffffff", "#000000");
  c.circle(c.width / 2, c.height / 2, radius1);

  c.stroke("C", "#ffffff", "#000000");
  c.fill("vvvv", "#ffffff", "#000000");
  c.circle(c.width / 2, c.height / 2, radius2);
}
