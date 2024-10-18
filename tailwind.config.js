/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./ui/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [require("@catppuccin/tailwindcss")({ prefix: "ctp", defaultFlavour: "macchiato" })],
};
