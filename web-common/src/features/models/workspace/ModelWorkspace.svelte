<script lang="ts">
  import { EntityType } from "@rilldata/web-common/features/entity-management/types";
  import { appStore } from "@rilldata/web-local/lib/application-state-stores/app-store";
  import { WorkspaceContainer } from "@rilldata/web-local/lib/components/workspace";
  import Tab from "../../../components/tab/Tab.svelte";
  import TabGroup from "../../../components/tab/TabGroup.svelte";
  import ModelInspectorGpt from "../gpt/ModelInspectorGPT.svelte";
  import ModelInspector from "./inspector/ModelInspector.svelte";
  import ModelBody from "./ModelBody.svelte";
  import ModelWorkspaceHeader from "./ModelWorkspaceHeader.svelte";

  export let modelName: string;
  export let focusEditorOnMount = false;

  const switchToModel = async (modelName: string) => {
    if (!modelName) return;

    appStore.setActiveEntity(modelName, EntityType.Model);
  };

  $: switchToModel(modelName);

  let selectedInspectorTab: "profile" | "gpt" = "profile";
</script>

{#key modelName}
  <WorkspaceContainer assetID={modelName}>
    <div slot="header">
      <ModelWorkspaceHeader {modelName} />
    </div>
    <div slot="body">
      <ModelBody {modelName} {focusEditorOnMount} />
    </div>
    <div slot="inspector">
      <div class="mx-2">
        <TabGroup
          variant="simple"
          on:select={(event) => {
            selectedInspectorTab = event.detail;
          }}
        >
          <Tab
            compact
            selected={selectedInspectorTab === "profile"}
            value={"profile"}>Profiler</Tab
          >
          <Tab compact selected={selectedInspectorTab === "gpt"} value={"gpt"}
            >GPT3</Tab
          >
        </TabGroup>
      </div>
      <hr />
      {#if selectedInspectorTab === "profile"}
        <ModelInspector {modelName} />
      {:else if selectedInspectorTab === "gpt"}
        <ModelInspectorGpt {modelName} />
      {/if}
    </div>
  </WorkspaceContainer>
{/key}
