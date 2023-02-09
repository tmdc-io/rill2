<script lang="ts">
  import { runtimeStore } from "@rilldata/web-local/lib/application-state-stores/application-store";
  import { invalidateAfterReconcile } from "@rilldata/web-local/lib/svelte-query/invalidation";
  import { useQueryClient } from "@sveltestack/svelte-query";
  import {
    useRuntimeServicePutFileAndReconcile,
    V1ReconcileResponse,
  } from "../../../runtime-client";
  import { fileArtifactsStore } from "../../entity-management/file-artifacts-store";
  import EditSql from "./EditSQL.svelte";
  import GenerateSql from "./GenerateSQL.svelte";
  import { useGetSourcePreview } from "./getSourcePreview";

  export let modelName: string;

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
    Leverage OpenAI to generate and edit your SQL code. Your source schemata
    will be fed into the prompt.
  </div>
  <hr />
  <GenerateSql {sourcePreview} on:sql={(e) => useSql(e.detail.sql)} />
  <hr />
  <EditSql {modelName} {sourcePreview} on:sql={(e) => useSql(e.detail.sql)} />
</div>
