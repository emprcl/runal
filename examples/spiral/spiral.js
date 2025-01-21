let cols, rows, size;
let t = 0;
let background = "#000000";
let colors = ["#fcf6bd", "#d0f4de", "#a9def9", "#e4c1f9", "#ff99c8"];

function setup() {
  c.background(" ", background, background);
  c.cellPadding(" ");
}

function draw() {
  c.clear();
  size = 1;
  cols = c.width / size;
  rows = c.height / size;

  for (let i = 0; i < cols; i++) {
    for (let j = 0; j < cols; j++) {
      let x = i * size;
      let y = j * size;
      let d = c.distance(x, y, c.width / 2, c.height / 2);
      let k = 0.6;
      let dx = x - c.width / 2;
      let dy = y - c.height / 2;
      let angle = Math.atan2(dy, dx);
      let spiralPath = Math.sin(d / k + angle + t);
      let df = 2;
      let af = 2;
      threshold = Math.sin(d / df + af * angle);

      c.stroke("⬤", colorGradient(c.width, d), background);

      if (spiralPath > threshold) {
        c.point(x, y);
      }
    }
  }

  t += 0.5;
}

function colorGradient(width, d) {
  let step = width / colors.length;
  for (let i = 0; i < colors.length; i++) {
    if (d <= (i + 1) * step) {
      return colors[i];
    }
  }
  return colors[colors.length - 1];
}
