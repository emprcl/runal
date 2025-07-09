let duration = 5;
let margin = 3;
let pointsNb = 26;
let points = [];

let a = 5;
let b = 2;

function setup(c) {
  c.size(40, 20);
  c.backgroundBg("197");
  c.stroke(".", "255", "197");
}

function draw(c) {
  c.clear();

  let theta = c.loopAngle(duration);
  let x = c.map(Math.sin(a * theta), -1, 1, margin, c.width - margin);
  let y = c.map(Math.sin(b * theta), -1, 1, margin, c.height - margin);
  let x2 = c.map(Math.sin(b * theta), -1, 1, margin, c.width - margin);
  let y2 = c.map(Math.sin(a * theta), -1, 1, margin, c.height - margin);
  if (points.length <= pointsNb * 2) {
    points.push(new Point(x, y, "0"));
    points.push(new Point(x2, y2, "#"));
  }

  for (let i = 0; i < points.length; i++) {
    points[i].update(c);
  }
}

class Point {
  constructor(x, y, char) {
    this.x = x;
    this.y = y;
    this.char = char;
    this.duration = pointsNb;
  }

  update(c) {
    c.strokeText(this.char);
    c.point(this.x, this.y);
    this.duration--;
    if (this.duration < 0) {
      points.shift();
    }
  }
}

function onKey(c, e) {
  if (e.key == "space") {
    c.noiseSeed(Date.now());
  }
}
