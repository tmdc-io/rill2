import adapter from "@sveltejs/adapter-static";
import preprocess from "svelte-preprocess";

/** @type {import('@sveltejs/kit').Config} */
const config = {
  // Consult https://github.com/sveltejs/svelte-preprocess
  // for more information about preprocessors
  preprocess: preprocess(),

  kit: {
    adapter: adapter({
      fallback: "index.html",
    }),
    files: {
      assets: "../web-common/static",
    },
    ...(process.env.BASE_PATH ? {
      paths: {
        base: `/${process.env.BASE_PATH}`,
        relative: true,
      }
    }: {})
  },
};

export default config;
