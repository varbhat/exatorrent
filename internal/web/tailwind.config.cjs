module.exports = {
  mode: 'jit',
  content: ['./src/*.{html,js,svelte,ts}', './src/**/*.{html,js,svelte,ts}'],
  plugins: [
    require('@tailwindcss/forms')({
      strategy: 'class'
    })
  ]
};
