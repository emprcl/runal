let xList = [1, 21, 31, 41];

function setup(c) {}

function draw(c) {
  c.clear();

  // rects

  c.stroke("1234567890", "#fffb00", "#0004ff");
  c.fill("1234567890", "#0004ff", "#fffb00");
  c.rect(1, 1, 11, 5);
  c.stroke("1234567890", "255", "#ee0000");
  c.text("RECT", 5, 2);
  c.text("Stroke: YES", 5, 3);
  c.text("Fill:   YES", 5, 4);

  c.stroke("1234567890", "#fffb00", "#0004ff");
  c.noFill();
  c.rect(21, 1, 11, 5);
  c.stroke("1234567890", "255", "#ee0000");
  c.text("RECT", 25, 2);
  c.text("Stroke: YES", 25, 3);
  c.text("Fill:   NO ", 25, 4);

  c.fill("1234567890", "#0004ff", "#fffb00");
  c.noStroke();
  c.rect(41, 1, 11, 5);
  c.stroke("1234567890", "255", "#ee0000");
  c.text("RECT", 45, 2);
  c.text("Stroke: NO ", 45, 3);
  c.text("Fill:   YES", 45, 4);

  // circles

  c.stroke("1234567890", "#fffb00", "#0004ff");
  c.fill("1234567890", "#0004ff", "#fffb00");
  c.circle(7, 13, 5);
  c.stroke("1234567890", "255", "#ee0000");
  c.text("CIRCLE", 5, 12);
  c.text("Stroke: YES", 5, 13);
  c.text("Fill:   YES", 5, 14);

  c.stroke("1234567890", "#fffb00", "#0004ff");
  c.noFill();
  c.circle(27, 13, 5);
  c.stroke("1234567890", "255", "#ee0000");
  c.text("CIRCLE", 25, 12);
  c.text("Stroke: YES", 25, 13);
  c.text("Fill:   NO ", 25, 14);

  c.fill("1234567890", "#0004ff", "#fffb00");
  c.noStroke();
  c.circle(47, 13, 5);
  c.stroke("1234567890", "255", "#ee0000");
  c.text("CIRCLE", 45, 12);
  c.text("Stroke: NO ", 45, 13);
  c.text("Fill:   YES", 45, 14);

  // ellipses

  c.stroke("1234567890", "#fffb00", "#0004ff");
  c.fill("1234567890", "#0004ff", "#fffb00");
  c.ellipse(7, 25, 5, 5);
  c.stroke("1234567890", "255", "#ee0000");
  c.text("ELLIPSE", 5, 24);
  c.text("Stroke: YES", 5, 25);
  c.text("Fill:   YES", 5, 26);

  c.stroke("1234567890", "#fffb00", "#0004ff");
  c.noFill();
  c.ellipse(27, 25, 5, 5);
  c.stroke("1234567890", "255", "#ee0000");
  c.text("ELLIPSE", 25, 24);
  c.text("Stroke: YES", 25, 25);
  c.text("Fill:   NO ", 25, 26);

  c.fill("1234567890", "#0004ff", "#fffb00");
  c.noStroke();
  c.ellipse(47, 25, 5, 5);
  c.stroke("1234567890", "255", "#ee0000");
  c.text("ELLIPSE", 45, 24);
  c.text("Stroke: NO ", 45, 25);
  c.text("Fill:   YES", 45, 26);

  // quads

  c.stroke("1234567890", "#fffb00", "#0004ff");
  c.fill("1234567890", "#0004ff", "#fffb00");
  c.quad(64, 1, 74, 1, 74, 6, 64, 6);
  c.stroke("1234567890", "255", "#ee0000");
  c.text("QUAD", 70, 2);
  c.text("Stroke: YES", 70, 3);
  c.text("Fill:   YES", 70, 4);

  c.stroke("1234567890", "#fffb00", "#0004ff");
  c.noFill();
  c.quad(84, 1, 94, 1, 94, 6, 84, 6);
  c.stroke("1234567890", "255", "#ee0000");
  c.text("QUAD", 88, 2);
  c.text("Stroke: YES", 88, 3);
  c.text("Fill:   NO ", 88, 4);

  c.fill("1234567890", "#0004ff", "#fffb00");
  c.noStroke();
  c.quad(104, 1, 114, 1, 114, 6, 104, 6);
  c.stroke("1234567890", "255", "#ee0000");
  c.text("QUAD", 108, 2);
  c.text("Stroke: NO ", 108, 3);
  c.text("Fill:   YES", 108, 4);

  // triangles

  c.stroke("1234567890", "#fffb00", "#0004ff");
  c.fill("1234567890", "#0004ff", "#fffb00");
  c.triangle(64, 9, 74, 9, 74, 18);

  c.stroke("1234567890", "#fffb00", "#0004ff");
  c.noFill();
  c.triangle(84, 9, 94, 9, 94, 18);

  c.fill("1234567890", "#0004ff", "#fffb00");
  c.noStroke();
  c.triangle(104, 9, 114, 9, 114, 18);

  // line

  c.stroke("1234567890", "#fffb00", "#0004ff");
  c.fill("1234567890", "#0004ff", "#fffb00");
  c.line(64, 25, 74, 25);

  c.stroke("1234567890", "#fffb00", "#0004ff");
  c.noFill();
  c.line(84, 25, 94, 25);

  c.fill("1234567890", "#0004ff", "#fffb00");
  c.noStroke();
  c.line(104, 25, 114, 25);
}
