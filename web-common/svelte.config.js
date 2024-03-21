import preprocess from "svelte-preprocess";

const config = {
  preprocess: preprocess(),
  ...(process.env.BASE_PATH ? {
    kit: {
      paths: {
        base: `/${process.env.BASE_PATH}`,
        relative: true,
      },
    }
  } : {})
};

export default config;
