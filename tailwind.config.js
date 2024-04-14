/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    'templates/*.templ',
  ],
  theme: {
    extend: {
      fontFamily: {
          header: ['"Press Start 2P"', 'system-ui'], // Define custom font-family for headers
          main: ['Roboto', 'sans-serif'], // Define custom font-family for body text
      }
    }
  },
  plugins: []
}
