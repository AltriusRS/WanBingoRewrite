import type {NextConfig} from "next";

const nextConfig: NextConfig = {
    images: {
        remotePatterns: [
            {
                protocol: "https",
                hostname: "i.ytimg.com",
                pathname: "/vi/**",
            },
            {
                protocol: "https",
                hostname: "cataas.com",
                pathname: "/cat",
            },
            {
                protocol: "https",
                hostname: "pbs.floatplane.com",
                pathname: "/stream_thumbnails/**",
            },
        ],
    },
};

export default nextConfig;
