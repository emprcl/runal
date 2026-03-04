// dissolve — noise-driven glyph landscape
//
// Maps 2D noise values to character pools of increasing visual density.
// The glyph choice itself is the visual — no shapes, no lines, just
// character weight shifting across the screen as the noise field
// drifts through time.

// Character pools ordered by visual density.
// Curating these is where the aesthetic lives —
// different pools make completely different landscapes
// from the same noise field.
var pools = [
  " ",
  "·.:,",
  "░╌┆╍",
  "▒▌▐╎",
  "▓█▀▄"
];

function setup(c) {
  c.fps(10);
}

function draw(c) {
  var t = c.framecount * 0.008;

  for (var y = 0; y < c.height; y++) {
    for (var x = 0; x < c.width; x++) {
      var n = c.noise2D(x * 0.04 + t, y * 0.06 + t * 0.3);

      // Select pool by noise band
      var band = Math.floor(c.map(n, 0, 1, 0, pools.length));
      if (band >= pools.length) band = pools.length - 1;
      var pool = pools[band];

      // Random character from the pool
      var ch = pool[Math.floor(c.random(0, pool.length))];

      // Brightness follows density
      var fg = Math.floor(c.map(n, 0, 1, 236, 255));
      c.stroke(ch, fg, "0");
      c.point(x, y);
    }
  }
}
