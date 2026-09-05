/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#2563EB',
          hover: '#1D4ED8',
          light: '#DBEAFE',
        },
        sidebar: {
          bg: '#0F172A',
          hover: '#1E293B',
          active: '#2563EB',
        }
      }
    },
  },
  plugins: [],
}
