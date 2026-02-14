- The highlighting constants for SKK markers have been renamed to describe their visual shapes (▽ and ▼) rather than their colors. (#2)
  - `WhiteMarkerHighlight` → `TriangleOutlineHighlight` (▽)
  - `BlackMarkerHighlight` → `TriangleFilledHighlight` (▼)
- Integrated `github.com/hymkor/sxencode-go` into the `internal` package. (#3)
- Improved Makefile (#4)
  - Cross-platform support: Enhanced compatibility across UNIX-like systems and Windows.
  - The build process now prioritizes go1.20.14 while falling back to the default go command if unavailable.

v0.6.1
======
Nov 13, 2025

- Maintenance: update dependencies and address staticcheck warnings (#1)
    - Mark `Coloring` as deprecated for compatibility.
    - Fix staticcheck issue (S1001) in `lisp.go`.
    - Update dependency `go-readline-ny` to v1.12.3.

v0.6.0
======
Sep 3, 2025

- Enabled conversion and word registration for words containing slashes in the conversion result
- Added support for evaluating certain Emacs Lisp forms in conversion results, such as `(concat)`, `(pwd)`, `(substring)`, and `(skk-current-date)` (but not `(lambda)` yet)

v0.5.0
======
Jan 29 2025

- Support the new syntax highlighting of go-readline-ny v1.7.4 (See also example2.go)

v0.4.2
======
Nov 28 2024

- Implement `z ` to `\u3000`

v0.4.0
======
Oct 06 2024

- Implement the Hankaku-Kana mode (Ctrl-Q)

v0.3.1
======
Oct 19 2023

- Fix the problem that `UTta` and `UTTa` were converted `打っtあ` and `▽う*t*t` instead of `打った`

v0.3.0
======
Oct 08 2023

- Fix: manually input inverted triangles were recognized as conversion markers

v0.2.0
======
Oct 08 2023

- Add the following the romaji-kana conversions:
    - `z,`→`‥`, `z-`→`～`, `z.`→`…`, `z/`→`・`, `z[`→`『`, `z]`→`』`,
        `z1`→`○`, `z2`→`▽`, `z3`→`△`, `z4`→`□`, `z5`→`◇`,
        `z6`→`☆`, `z7`→`◎`, `z8`→`〔`, `z9`→`〕`, `z0`→`∞`,
        `z^`→`※`, `z\\`→`￥`, `z@`→`〃`, `z;`→`゛`, `z:`→`゜` ,
        `z!`→`●`, `z"`→`▼`, `z#`→`▲`, `z$`→`■ `, `z%`→`◆`,
        `z&`→`★`, `z'`→`♪`, `z(`→`【`, `z)`→`】`, `z=`→`≒`,
        `z~`→`≠`, `z|`→`〒`, ``z` ``→`“`, `z+`→`±`, `z*`→`×`,
        `z<`→`≦`, `z>`→`≧`, `z?`→`÷`, `z_`→`―`,
    - `bya`→`びゃ` or `ビャ` ... `byo`→`びょ` or `ビョ`
    - `pya`→`ぴゃ` or `ピャ` ... `pyo`→`ぴょ` or `ピョ`
    - `tha`→`てぁ` or `テァ` ... `tho`→`てょ` or `テョ`
- Implement `q` that convert mutually between Hiragana and Katakana during conversion.

v0.1.0
======
Oct 06 2023

- The first version for nyagos 4.4.14\_0
