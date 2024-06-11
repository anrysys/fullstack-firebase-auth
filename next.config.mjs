/** @type {import('next').NextConfig} */
const nextConfig = {
    serverRuntimeConfig: {
    //port: process.env.NEXT_PUBLIC_APP_CLIENT_PORT || 80,
  },
  // webpack: (config, { isServer }) => {
  //   // Only run for server
  //   if (isServer) {
  //     config.externals = {
  //       //...config.externals,
  //       'barcode-website': 'commonjs barcode-website',
  //     };
  //   }

  //   return config;
  // },  
}

export default nextConfig
