let x1 = 0,
  y1 = 0,
  x2 = 0,
  y2 = 0;

function setup(c) {
  c.stroke(
    "make yourselves sheep and the wolves will eat you ",
    "#000000",
    "#000000",
  );
  c.noLoop();
}

function draw(c) {
  if (x1 == 0 && y1 == 0) {
    return;
  }
  c.line(x1, y1, x2, y2);
}

function onMouse(c, event) {
  if (event.type != "click") {
    return;
  }
  // set stroke color to one of the ansi colors, but not black (1)
  c.strokeFg("" + Math.ceil(c.random(1, 255)));
  x1 = x2;
  y1 = y2;
  if (event.button == "left") {
    x2 = c.mouseX;
    y2 = c.mouseY;
  }
  c.redraw();
}
