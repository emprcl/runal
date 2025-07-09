function setup(c) {}

function draw(c) {
  c.clear();
  c.circle(c.mouseX, c.mouseY, 5);
}

function onMouse(c, event) {
  c.backgroundBg(c.random(0, 255));
}
