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
