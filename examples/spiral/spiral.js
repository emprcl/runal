let cols, rows, size;
let t = 0;

function setup() {}

function draw() {
  runal.flush();
  size = 1;
  cols = runal.width() / size;
  rows = runal.height() / size;

  for (let i = 0; i < cols; i++) {
    for (let j = 0; j < cols; j++) {
      let x = i * size;
      let y = j * size;
      let d = runal.dist(x, y, runal.width() / 2, runal.height() / 2);
      let k = 0.6;
      let dx = x - runal.width() / 2;
      let dy = y - runal.height() / 2;
      let angle = Math.atan2(dy, dx);
      let spiralPath = Math.sin(d / k + angle + t);
      let df = 2;
      let af = 2;
      threshold = Math.sin(d / df + af * angle);

      if (spiralPath > threshold) {
        runal.text("⬤", x, y);
      }
    }
  }

  t += 0.5;
}
