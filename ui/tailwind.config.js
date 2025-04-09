/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './html/**/*.{html,go}',  // HTML templates in `ui/html`
    './static/js/**/*.{js}',  // JS files in `ui/static/js`
  ],
  theme: {
    extend: {
      colors: {
        bbPrimary: '#6b041f', // Burgundy
        bbAccent: '#e49629',   // Mustard
      },
      fontFamily: {
        poppins: ['Poppins', 'sans-serif'],
      },
    },
  },
  plugins: [],
}


