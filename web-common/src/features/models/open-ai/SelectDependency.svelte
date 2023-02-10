<script lang="ts">
  import { runtimeStore } from "@rilldata/web-local/lib/application-state-stores/application-store";
  import { createEventDispatcher, onMount } from "svelte";
  import { SelectMenu } from "../../../components/menu";
  import { runtimeServiceQuery } from "../../../runtime-client";

  let selection;
  let options;

  const dispatch = createEventDispatcher();

  onMount(async () => {
    const sourceNames = await getAllSourceNames($runtimeStore.instanceId);
    options = sourceNames.map((name) => {
      return {
        key: name,
        main: name,
      };
    });
    selection = options[0];
    dispatch("select", { dependency: selection?.key });
  });

  async function getAllSourceNames(instanceId: string): Promise<string[]> {
    // Note: currently excluding models (VIEWs in DuckDB) because they don't include the "rowid" meta-column, which we use to get the data tail
    const resp = await runtimeServiceQuery(instanceId, {
      sql: "select table_name from information_schema.tables where table_schema = 'main' and table_type = 'BASE TABLE';",
      priority: 1,
    });

    return resp.data.map((row) => row["table_name"]);
  }
</script>

<div class="flex flex-col gap-x-2">
  <div style="font-size: 11px;">Source</div>
  <td class="border rounded border-gray-400 max-w-max mt-1" style:height="32px">
    <SelectMenu
      {options}
      block
      tailwindClasses="overflow-hidden"
      alignment="end"
      bind:selection
      on:select={() => dispatch("select", { dependency: selection?.key })}
    >
      {selection?.main}
    </SelectMenu>
  </td>
</div>
