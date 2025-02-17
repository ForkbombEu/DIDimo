const { parser } = require("./src/generateParser");
const fs = require("fs");
const src = "./src/";
const lib = "./lib/";

// Compile parser
const parserSourceCode = parser.generate();
fs.writeFileSync(lib + "parser.js", parserSourceCode);

// Copy other files
for (const f of fs.readdirSync(src)) {
  if (f === "generateParser.js") continue;

  const s = fs.statSync(src + f);
  if (s.isDirectory()) continue;

  fs.copyFileSync(src + f, lib + f);
}
