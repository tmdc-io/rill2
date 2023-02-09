<script lang="ts">
  import { runtimeStore } from "@rilldata/web-local/lib/application-state-stores/application-store";
  import { invalidateAfterReconcile } from "@rilldata/web-local/lib/svelte-query/invalidation";
  import { useQueryClient } from "@sveltestack/svelte-query";
  import { createForm } from "svelte-forms-lib";
  import { Button } from "../../../components/button";
  import Input from "../../../components/forms/Input.svelte";
  import {
    useRuntimeServicePutFileAndReconcile,
    V1ReconcileResponse,
  } from "../../../runtime-client";
  import { fileArtifactsStore } from "../../entity-management/file-artifacts-store";
  import Spinner from "../../entity-management/Spinner.svelte";
  import { EntityStatus } from "../../entity-management/types";
  import FullPrompt from "./FullPrompt.svelte";
  import { useGetSourcePreview } from "./getSourcePreview";

  export let modelName: string;

  /**
   * Hack: currently, we call the OpenAI API client-side, so we need to pass the API key from the backend to the frontend.
   * TODO: make the OpenAI API calls from the backend.
   */
  $: OPENAI_API_KEY = $runtimeStore.openAIAPIKey;

  let sourcePreview: string;
  const sourcePreviewQuery = useGetSourcePreview(); // TODO: this should be a query not a mutation
  $sourcePreviewQuery.mutate(
    {
      data: {
        instanceId: $runtimeStore.instanceId,
        sourceName: "UFO_Reports",
      },
    },
    {
      onSuccess: (resp) => {
        sourcePreview = resp.sourcePreview;
      },
    }
  );

  let sql: string;
  $: prompt = `#\tDuckDB SQL\n#\n${sourcePreview}\n#\n#\tA query to answer the question: "${
    $form["query"] ?? "[Your question here]"
  }"\n\nSELECT `;
  let isLoading: boolean;
  let error: string;

  const { form, errors, handleSubmit } = createForm({
    initialValues: {
      prompt: "",
    },
    onSubmit: async () => {
      isLoading = true;
      const response = await fetch("https://api.openai.com/v1/completions", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${OPENAI_API_KEY}`,
        },
        body: JSON.stringify({
          model: "code-davinci-002",
          prompt: prompt,
          temperature: 0,
          max_tokens: 150,
          top_p: 1,
          frequency_penalty: 0,
          presence_penalty: 0,
          stop: ["#", ";"],
        }),
      });

      // Postprocess the response
      const data = await response.json();
      // if there's an error, show it
      if (data.error) {
        error = data.error.message;
        isLoading = false;
        return;
      }
      // // add the primer back to the beginning of the response
      sql = "SELECT " + data.choices[0].text;
      // insert a newline before each FROM, WHERE, GROUP BY, ORDER BY, LIMIT
      // sql = sql.replace(/(FROM|WHERE|GROUP BY|ORDER BY|LIMIT)/g, "\n$1"); // TODO: this is only needed sometimes
      // prefix the sql with an informative comment
      sql = `-- Query: ${$form["query"]}\n\n${sql}`;
      useSql(sql);
      isLoading = false;
    },
  });

  const updateModel = useRuntimeServicePutFileAndReconcile();
  const queryClient = useQueryClient();

  function useSql(sql: string) {
    $updateModel.mutateAsync(
      {
        data: {
          instanceId: $runtimeStore.instanceId,
          path: `/models/${modelName}.sql`,
          blob: sql,
        },
      },
      {
        onSuccess: (resp: V1ReconcileResponse) => {
          fileArtifactsStore.setErrors(resp.affectedPaths, resp.errors);
          invalidateAfterReconcile(queryClient, $runtimeStore.instanceId, resp);
        },
        onError: (err) => {
          console.error(err);
        },
      }
    );
  }

  let showFullPromptModal = false;
  function openFullPromptModal() {
    showFullPromptModal = true;
  }
</script>

<div class="flex flex-col gap-y-4 flex-grow m-4">
  <div>
    Leverage OpenAI to generate and edit your SQL code. Your source schemata
    will be fed into the prompt.
  </div>
  <hr />
  <div class="font-bold">Generate SQL</div>
  <div class="flex flex-col gap-y-2">
    Ask a question about your data! For example, you might ask of a COVID
    dataset:
    <ul>
      <li>What country had the most COVID cases in 2022?</li>
      <li>What is the percent change of COVID cases each month?</li>
    </ul>
  </div>
  <form
    id="openai-form"
    autocomplete="off"
    on:submit|preventDefault={handleSubmit}
  >
    <Input
      bind:value={$form["query"]}
      claimFocusOnMount
      error={$errors["query"]}
      id="query"
      label="Question"
      disabled={isLoading}
    />
    <div class="flex flex-row gap-x-2 my-4">
      <Button type="secondary" on:click={openFullPromptModal}
        >See full prompt</Button
      >
      <Button
        type="primary"
        submitForm
        form="openai-form"
        disabled={!$form["query"] || isLoading}>Generate SQL</Button
      >
      {#if isLoading}
        <div class="flex flex-row gap-x-2 items-center">
          <Spinner size="1.5em" status={EntityStatus.Running} />
        </div>
      {/if}
    </div>
    {#if error}
      <div class="text-red-500">{error}</div>
    {/if}
  </form>
</div>

{#if showFullPromptModal}
  <FullPrompt
    {prompt}
    on:close={() => {
      showFullPromptModal = false;
    }}
  />
{/if}
