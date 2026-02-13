import tailwind from 'bun-plugin-tailwind'

const result = await Bun.build({
  outdir: './internal/ui/assets',
  minify: true,
  sourcemap: 'linked',
  naming: {
    chunk: '[name].[ext]',
    asset: '[name].[ext]',
  },
  entrypoints: ['./ui/index.html'],
  publicPath: '/assets/',
  plugins: [tailwind],
})

console.log(result)
