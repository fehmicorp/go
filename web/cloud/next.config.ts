import type { NextConfig } from 'next';
import path from 'path';
// const outputMode = 'export';
const outputMode = process.env.OUTPUT === 'local' ? 'export' : 'standalone';
const nextConfig: NextConfig = {
  output: outputMode as 'export' | 'standalone',
  images: {
    unoptimized: true,
  },
  turbopack: {
    root: path.resolve(__dirname, '..'),
    resolveAlias: {
      '@config': path.resolve(__dirname, '../config'),
    },
  },
};

export default nextConfig;