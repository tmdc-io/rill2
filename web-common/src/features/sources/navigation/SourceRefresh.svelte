<script lang="ts">
  import RefreshIcon from "@rilldata/web-common/components/icons/RefreshIcon.svelte";
  import { MenuItem } from "@rilldata/web-common/components/menu";
  import { getFilePathFromNameAndType } from "@rilldata/web-common/features/entity-management/entity-mappers";
  import { getFileHasErrors } from "@rilldata/web-common/features/entity-management/resources-store";
  import {
    useIsLocalFileConnector,
    useSource,
    useSourceFromYaml,
  } from "@rilldata/web-common/features/sources/selectors";
  import { overlay } from "@rilldata/web-common/layout/overlay-store";
  import type { V1SourceV2 } from "@rilldata/web-common/runtime-client";
  import { V1ReconcileStatus } from "@rilldata/web-common/runtime-client";
  import { useQueryClient } from "@tanstack/svelte-query";
  import { createEventDispatcher } from "svelte";
  import { runtime } from "../../../runtime-client/runtime-store";
  import { EntityType } from "../../entity-management/types";
  import { refreshSource } from "../refreshSource";
  import Button from "@rilldata/web-common/components/button/Button.svelte";

  export let sourceName: string;
  // manually toggle menu to workaround: https://stackoverflow.com/questions/70662482/react-query-mutate-onsuccess-function-not-responding
  $: filePath = getFilePathFromNameAndType(sourceName, EntityType.Table);

  const queryClient = useQueryClient();

  $: runtimeInstanceId = $runtime.instanceId;

  const dispatch = createEventDispatcher();

  $: sourceQuery = useSource(runtimeInstanceId, sourceName);
  let source: V1SourceV2 | undefined;
  $: source = $sourceQuery.data?.source;

  $: sourceHasError = getFileHasErrors(
    queryClient,
    runtimeInstanceId,
    filePath,
  );
  $: sourceIsIdle =
    $sourceQuery.data?.meta?.reconcileStatus ===
    V1ReconcileStatus.RECONCILE_STATUS_IDLE;

  $: sourceFromYaml = useSourceFromYaml($runtime.instanceId, filePath);

  const onRefreshSource = async (tableName: string) => {
    const connector: string | undefined =
      source?.state?.connector ?? $sourceFromYaml.data?.type;
    if (!connector) {
      // if parse failed or there is no catalog entry, we cannot refresh source
      // TODO: show the import source modal with fixed tableName
      return;
    }
    try {
      await refreshSource(connector, tableName, runtimeInstanceId);
    } catch (err) {
      // no-op
    }
    overlay.set(null);
  };
</script>

<Button
  type="text"
  on:click={(e) => {
    e.stopPropagation();
    onRefreshSource(sourceName);
  }}
>
  <RefreshIcon />
</Button>
