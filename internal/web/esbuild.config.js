import { existsSync, mkdirSync, copyFile } from 'fs';
import { build } from 'esbuild';
import sveltePlugin from 'esbuild-svelte';
import sveltePreprocess from 'svelte-preprocess';
import tailwindcss from 'tailwindcss';

const WATCH = process.argv.includes('-w');

//make sure the directoy exists before stuff gets put into it
if (!existsSync('./build/')) {
  mkdirSync('./build/');
}

copyFile('./src/index.html', './build/index.html', (err) => {
  if (err) throw err;
});

//build the application
build({
  entryPoints: ['./src/index.js'],
  outdir: './build',
  format: 'esm',
  minify: true,
  bundle: true,
  treeShaking: true,
  splitting: true,
//  watch: WATCH,
 // incremental: WATCH,
  plugins: [
    sveltePlugin({
      compilerOptions: {
        dev: WATCH
      },
      preprocess: sveltePreprocess({
        aliases: ['ts', 'typescript'],
        postcss: {
          plugins: [tailwindcss()]
        }
      })
    })
  ],
  tsconfig: 'tsconfig.json'
}).catch((err) => {
  console.error(err);
  process.exit(1);
});
