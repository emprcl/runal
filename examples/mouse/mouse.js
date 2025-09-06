function setup(c) {}

function draw(c) {
  c.clear();
  c.circle(c.mouseX, c.mouseY, 5);
}

function onKey(c, event) {
  if (event.key == "space") {
    c.backgroundBg(c.random(0, 255));
  }
}

function onMouse(c, event) {
  if (event.type == "click") {
    c.backgroundBg(c.random(0, 255));
  }
}
