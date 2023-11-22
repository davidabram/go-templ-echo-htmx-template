/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./internals/templates/**/*.templ"],
  mode: 'jit',
  theme: {
    extend: {
      fontFamily: {
        sans: ['var(--font-geist-sans)'],
        mono: ['var(--font-geist-mono)'],
      },
    },
  },
  variants: {
    extend: {
      display: ['group-focus'],
      opacity: ['group-focus'],
      inset: ['group-focus'],
    },
  },
  plugins: [],
}
