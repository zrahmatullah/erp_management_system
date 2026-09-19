/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontSize: {
        '2xs': ['0.75rem', { lineHeight: '1rem' }],       // ~12.5px
        'xs': ['0.845rem', { lineHeight: '1.25rem' }],    // ~14px (upscaled from 12px)
        'sm': ['0.95rem', { lineHeight: '1.4rem' }],      // ~15.7px (upscaled from 14px)
        'base': ['1.0625rem', { lineHeight: '1.65rem' }], // ~17.5px (upscaled from 16px)
        'lg': ['1.2rem', { lineHeight: '1.75rem' }],      // ~19.8px (upscaled from 18px)
        'xl': ['1.375rem', { lineHeight: '1.875rem' }],   // ~22.7px (upscaled from 20px)
        '2xl': ['1.65rem', { lineHeight: '2.1rem' }],     // ~27.2px (upscaled from 24px)
        '3xl': ['2.1rem', { lineHeight: '2.5rem' }],
      },
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
