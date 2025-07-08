function setup(c) {}

function draw(c) {
  c.clear();
  c.circle(c.mouseX, c.mouseY, 5);
}

function onKey(c, e) {
  if (e.key == " ") {
    c.backgroundBg(c.random(0, 255));
  }
}

function onMouse(c, e) {
  if (e.action == 0) {
    c.backgroundBg(c.random(0, 255));
  }
}
