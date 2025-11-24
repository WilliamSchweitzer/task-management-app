/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Enable if you need to serve API routes alongside Next.js
  // or configure rewrites for your backend
  async rewrites() {
    return [
      // Uncomment and configure these when you have your Kong Gateway URL
      // {
      //   source: '/api/auth/:path*',
      //   destination: 'http://your-kong-gateway:8000/auth/:path*',
      // },
      // {
      //   source: '/api/tasks/:path*',
      //   destination: 'http://your-kong-gateway:8000/tasks/:path*',
      // },
    ];
  },
};

export default nextConfig;
