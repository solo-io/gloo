const plugin = require('tailwindcss/plugin');

module.exports = {
  content: ['./src/**/*.{js,jsx,ts,tsx}'],
  theme: {
    extend: {
      colors: {
        gray: {
          100: '#F9F9F9',
          200: '#F2F2F2',
          300: '#DFDFE0',
          400: '#D4D8DE',
          500: '#C2C8D0',
          600: '#ADB3BC',
          700: '#6E7477',
          800: '#35393B',
          900: '#1a202c',
        },
        blue: {
          100: '#F7FCFF',
          '150gloo': '#EAF8FF',
          200: '#DDF5FF',
          300: '#6AC7F0',
          400: '#54B7E3',
          500: '#2396C9',
          600: '#0F7FB1',
          700: '#253E58',
          '150gloo': '#EAF8FF',
          '200gloo': '#DDF5FF',
          '250gloo': '#B1E7FF',
          '300gloo': '#6AC7F0',
          '400gloo': '#54B7E3',
          '500gloo': '#2396C9',
          '600gloo': '#0F7FB1',
          '700gloo': '#253E58',
          brand: '#2997ca',
        },
      },
    },
  },
  plugins: [
    // @see https://tailwindcss.com/docs/plugins
    plugin(function({ addUtilities, addComponents, e, prefix, config }) {
      // Add your custom styles here
    }),
  ],
};
