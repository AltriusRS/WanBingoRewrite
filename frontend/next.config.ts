import type {NextConfig} from "next";

const nextConfig: NextConfig = {
    /* config options here */
    images: {
        remotePatterns: [
            new URL('https://i.ytimg.com/vi/**/maxresdefault_live.jpg'),
            new URL('https://cataas.com/cat?width=720&height=480'),
            new URL('https://pbs.floatplane.com/stream_thumbnails/*')
        ]
    }
};

export default nextConfig;
