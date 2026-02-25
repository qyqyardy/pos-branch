export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Manrope', 'ui-sans-serif', 'system-ui', 'sans-serif'],
        brand: ['Fraunces', 'ui-serif', 'Georgia', 'serif'],
      },
      keyframes: {
        'float-in': {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        'pulse-soft': {
          '0%, 100%': { opacity: '1' },
          '50%': { opacity: '.55' },
        },
      },
      animation: {
        'float-in': 'float-in 500ms ease-out both',
        'pulse-soft': 'pulse-soft 1.2s ease-in-out infinite',
      },
    },
  },
  plugins: [],
}
