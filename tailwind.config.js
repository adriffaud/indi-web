/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./components/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [require("@catppuccin/tailwindcss")({ prefix: "ctp", defaultFlavour: "macchiato" })],
};
