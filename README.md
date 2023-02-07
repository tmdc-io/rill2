# Rill Developer **_(tech preview)_**
Rill Developer makes it effortless to transform your datasets with SQL and create powerful, opinionated dashboards. Rill's principles:

- _**feels good to use**_ – powered by Sveltekit & DuckDB = conversation-fast, not wait-ten-seconds-for-result-set fast
- _**works with your local and remote datasets**_ – imports and exports Parquet and CSV (s3, gcs, https, local)
- _**no more data analysis "side-quests"**_ – helps you build intuition about your dataset through automatic profiling
- _**no "run query" button required**_ – responds to each keystroke by re-profiling the resulting dataset
- _**radically simple dashboards**_ – thoughtful, opinionated defaults to help you quickly derive insights from your data
- _**dashboards as code**_ – each step from data to dashboard has versioning, git sharing, and easy project rehydration 

## Pick an install option
You can get started in less than 2 minutes with our install script (Mac, Linux).

```
curl -s https://cdn.rilldata.com/install.sh | bash
```

After installation, launch the application by running `rill start` or try an example project:
```
rill init --example
```

![home-demo](https://user-images.githubusercontent.com/5587788/207410129-bd4fb84b-dc3d-494c-9cf1-2322fcf0d503.gif "770784519")

## Enable the SQL Assistant powered by OpenAI

Rill Developer includes an optional feature called the SQL Assistant. It leverages the latest deep learning technology from OpenAI to help you write your SQL code. If you'd like to use the feature, you'll need to sign up for an OpenAI account and generate an API key. Follow these steps:

1. [Sign up for an OpenAI account](https://openai.com/api/).
2. [Generate an API key](https://platform.openai.com/account/api-keys).
3. Set the `OPENAI_API_KEY` environment variable to your API key by running the following command in your terminal:

    ```bash
    export OPENAI_API_KEY="sk-..."
    ```
    Note: You may want to add the line above to your terminal dotfile (e.g. ~/.bashrc or ~/.zshrc) to persist the API key.

4. Start Rill with `rill start`.

The SQL Assistant will be enabled and available to use from the sidebar of SQL editor.

## We want to hear from you

You can [file an issue](https://github.com/rilldata/rill-developer/issues/new/choose) or reach us in our [Rill discord](https://bit.ly/3unvA05) channel. Please abide by the [rill community policy](https://github.com/rilldata/rill-developer/blob/main/COMMUNITY-POLICY.md).


## Legal
By downloading and using our application you are agreeing to the [Privacy Policy](https://www.rilldata.com/legal/privacy) and [Rill Terms of Service](https://www.rilldata.com/legal/tos).




