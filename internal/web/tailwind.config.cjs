const colors = require('tailwindcss/colors');
const { colors: defaultColors } = require('tailwindcss/defaultTheme');
module.exports = {
  mode: 'jit',
  purge: {
    enabled: true,
    content: ['./src/*.{html,js,svelte,ts}', './src/**/*.{html,js,svelte,ts}']
  },
  darkMode: 'media',
  theme: {
    colors: {
      ...defaultColors,
      gray: colors.trueGray,
      secGray: colors.gray
    },
    extend: {}
  },
  variants: {
    extend: {}
  },
  plugins: [
    require('@tailwindcss/forms')({
      strategy: 'class'
    })
  ]
};
