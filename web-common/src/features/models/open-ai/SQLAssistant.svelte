<script lang="ts">
  import { runtimeStore } from "@rilldata/web-local/lib/application-state-stores/application-store";
  import { invalidateAfterReconcile } from "@rilldata/web-local/lib/svelte-query/invalidation";
  import { useQueryClient } from "@sveltestack/svelte-query";
  import Button from "../../../components/button/Button.svelte";
  import {
    useRuntimeServicePutFileAndReconcile,
    V1ReconcileResponse,
  } from "../../../runtime-client";
  import { fileArtifactsStore } from "../../entity-management/file-artifacts-store";
  import EditSql from "./EditSQL.svelte";
  import GenerateSql from "./GenerateSQL.svelte";
  import { useGetSourcePreview } from "./getSourcePreview";
  import SelectDependency from "./SelectDependency.svelte";

  export let modelName: string;

  let moduleSelection = "generate";

  let sourcePreview: string;
  const sourcePreviewQuery = useGetSourcePreview(); // TODO: this should be a query not a mutation

  function handleSelectDependency(event: CustomEvent) {
    if (!event.detail.dependency) return;
    const dependencyName = event.detail.dependency;
    $sourcePreviewQuery.mutate(
      {
        data: {
          instanceId: $runtimeStore.instanceId,
          sourceName: dependencyName,
        },
      },
      {
        onSuccess: (resp) => {
          sourcePreview = resp.sourcePreview;
        },
        onError: (err) => {
          console.error(err);
        },
      }
    );
  }

  const queryClient = useQueryClient();
  const updateModel = useRuntimeServicePutFileAndReconcile();

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
</script>

<div class="flex flex-col gap-y-4 flex-grow m-4">
  <div>
    Leverage OpenAI to generate and edit your SQL code. Pick a source so that
    its schema and example data can be fed into the prompt.
  </div>
  <SelectDependency {modelName} on:select={handleSelectDependency} />
  <div class="flex flex-row gap-x-2">
    <Button
      type={moduleSelection === "generate" ? "secondary" : "text"}
      on:click={() => (moduleSelection = "generate")}>Generate</Button
    >
    <Button
      type={moduleSelection === "edit" ? "secondary" : "text"}
      on:click={() => (moduleSelection = "edit")}>Edit</Button
    >
  </div>
  {#if moduleSelection === "generate"}
    <GenerateSql {sourcePreview} on:sql={(e) => useSql(e.detail.sql)} />
  {:else if moduleSelection === "edit"}
    <EditSql {modelName} {sourcePreview} on:sql={(e) => useSql(e.detail.sql)} />
  {/if}
</div>
