import preprocess from "svelte-preprocess";

const config = {
  preprocess: preprocess(),
  ...(process.env.APP_BASE_PATH ? {
    kit: {
      paths: {
        base: `/${process.env.APP_BASE_PATH}`,
        relative: true,
      },
    }
  } : {})
};

export default config;
